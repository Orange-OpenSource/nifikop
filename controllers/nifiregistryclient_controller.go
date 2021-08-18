/*
Copyright 2020.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"fmt"
	"reflect"

	"github.com/Orange-OpenSource/nifikop/api/v1alpha1"
	"github.com/Orange-OpenSource/nifikop/pkg/clientwrappers/registryclient"
	"github.com/Orange-OpenSource/nifikop/pkg/k8sutil"
	"github.com/Orange-OpenSource/nifikop/pkg/util"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var registryClientFinalizer = "nifiregistryclients.nifi.orange.com/finalizer"

// NifiRegistryClientReconciler reconciles a NifiRegistryClient object
type NifiRegistryClientReconciler struct {
	client.Client
	Log             logr.Logger
	Scheme          *runtime.Scheme
	Recorder        record.EventRecorder
	RequeueInterval int
	RequeueOffset   int
}

// +kubebuilder:rbac:groups=nifi.orange.com,resources=nifiregistryclients,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=nifi.orange.com,resources=nifiregistryclients/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=nifi.orange.com,resources=nifiregistryclients/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the NifiRegistryClient object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.0/pkg/reconcile
func (r *NifiRegistryClientReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = r.Log.WithValues("nifiregistryclient", req.NamespacedName)

	var err error

	// Fetch the NifiRegistryClient instance
	var instance = &v1alpha1.NifiRegistryClient{}
	if err = r.Client.Get(ctx, req.NamespacedName, instance); err != nil {
		if apierrors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			return Reconciled()
		}
		// Error reading the object - requeue the request.
		return RequeueWithError(r.Log, err.Error(), err)
	}

	// Get the referenced NifiCluster
	clusterNamespace := GetClusterRefNamespace(instance.Namespace, instance.Spec.ClusterRef)
	var cluster *v1alpha1.NifiCluster
	if cluster, err = k8sutil.LookupNifiCluster(r.Client, instance.Spec.ClusterRef.Name, clusterNamespace); err != nil {
		// This shouldn't trigger anymore, but leaving it here as a safetybelt
		if k8sutil.IsMarkedForDeletion(instance.ObjectMeta) {
			r.Log.Info("Cluster is already gone, there is nothing we can do")
			if err = r.removeFinalizer(ctx, instance); err != nil {
				return RequeueWithError(r.Log, "failed to remove finalizer", err)
			}
			return Reconciled()
		}

		r.Recorder.Event(instance, corev1.EventTypeWarning, "ReferenceClusterError",
			fmt.Sprintf("Failed to lookup reference cluster : %s in %s",
				instance.Spec.ClusterRef.Name, clusterNamespace))
		// the cluster does not exist - should have been caught pre-flight
		return RequeueWithError(r.Log, "failed to lookup referenced cluster", err)
	}

	// Check if marked for deletion and if so run finalizers
	if k8sutil.IsMarkedForDeletion(instance.ObjectMeta) {
		return r.checkFinalizers(ctx, r.Log, instance, cluster)
	}

	r.Recorder.Event(instance, corev1.EventTypeNormal, "Reconciling",
		fmt.Sprintf("Reconciling registry client %s", instance.Name))

	// Check if the NiFi registry client already exist
	exist, err := registryclient.ExistRegistryClient(r.Client, instance, cluster)
	if err != nil {
		return RequeueWithError(r.Log, "failure checking for existing registry client", err)
	}

	if !exist {
		// Create NiFi registry client
		r.Recorder.Event(instance, corev1.EventTypeNormal, "Creating",
			fmt.Sprintf("Creating registry client %s", instance.Name))
		status, err := registryclient.CreateRegistryClient(r.Client, instance, cluster)
		if err != nil {
			return RequeueWithError(r.Log, "failure creating registry client", err)
		}

		instance.Status = *status
		if err := r.Client.Status().Update(ctx, instance); err != nil {
			return RequeueWithError(r.Log, "failed to update NifiRegistryClient status", err)
		}

		r.Recorder.Event(instance, corev1.EventTypeNormal, "Created",
			fmt.Sprintf("Created registry client %s", instance.Name))
	}

	// Sync RegistryClient resource with NiFi side component
	r.Recorder.Event(instance, corev1.EventTypeNormal, "Synchronizing",
		fmt.Sprintf("Synchronizing registry client %s", instance.Name))
	status, err := registryclient.SyncRegistryClient(r.Client, instance, cluster)
	if err != nil {
		r.Recorder.Event(instance, corev1.EventTypeNormal, "SynchronizingFailed",
			fmt.Sprintf("Synchronizing registry client %s failed", instance.Name))
		return RequeueWithError(r.Log, "failed to sync NifiRegistryClient", err)
	}

	instance.Status = *status
	if err := r.Client.Status().Update(ctx, instance); err != nil {
		return RequeueWithError(r.Log, "failed to update NifiRegistryClient status", err)
	}

	r.Recorder.Event(instance, corev1.EventTypeNormal, "Synchronized",
		fmt.Sprintf("Synchronized registry client %s", instance.Name))
	// Ensure NifiCluster label
	if instance, err = r.ensureClusterLabel(ctx, cluster, instance); err != nil {
		return RequeueWithError(r.Log, "failed to ensure NifiCluster label on registry client", err)
	}

	// Ensure finalizer for cleanup on deletion
	if !util.StringSliceContains(instance.GetFinalizers(), registryClientFinalizer) {
		r.Log.Info("Adding Finalizer for NifiRegistryClient")
		instance.SetFinalizers(append(instance.GetFinalizers(), registryClientFinalizer))
	}

	// Push any changes
	if instance, err = r.updateAndFetchLatest(ctx, instance); err != nil {
		return RequeueWithError(r.Log, "failed to update NifiRegistryClient", err)
	}

	r.Recorder.Event(instance, corev1.EventTypeNormal, "Reconciled",
		fmt.Sprintf("Reconciling registry client %s", instance.Name))

	r.Log.Info("Ensured Registry Client")

	interval := util.GetRequeueInterval(r.RequeueInterval, r.RequeueOffset)
	r.Log.Info(fmt.Sprintf("Will requeue registry client task after %v", interval))
	return RequeueAfter(interval)
}

// SetupWithManager sets up the controller with the Manager.
func (r *NifiRegistryClientReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.NifiRegistryClient{}).
		Complete(r)
}

func (r *NifiRegistryClientReconciler) ensureClusterLabel(ctx context.Context, cluster *v1alpha1.NifiCluster,
	registryClient *v1alpha1.NifiRegistryClient) (*v1alpha1.NifiRegistryClient, error) {

	labels := ApplyClusterRefLabel(cluster, registryClient.GetLabels())
	if !reflect.DeepEqual(labels, registryClient.GetLabels()) {
		registryClient.SetLabels(labels)
		return r.updateAndFetchLatest(ctx, registryClient)
	}
	return registryClient, nil
}

func (r *NifiRegistryClientReconciler) updateAndFetchLatest(ctx context.Context,
	registryClient *v1alpha1.NifiRegistryClient) (*v1alpha1.NifiRegistryClient, error) {

	typeMeta := registryClient.TypeMeta
	err := r.Client.Update(ctx, registryClient)
	if err != nil {
		return nil, err
	}
	registryClient.TypeMeta = typeMeta
	return registryClient, nil
}

func (r *NifiRegistryClientReconciler) checkFinalizers(ctx context.Context, reqLogger logr.Logger,
	registryClient *v1alpha1.NifiRegistryClient, cluster *v1alpha1.NifiCluster) (reconcile.Result, error) {

	reqLogger.Info("NiFi registry client is marked for deletion")
	var err error
	if util.StringSliceContains(registryClient.GetFinalizers(), registryClientFinalizer) {
		if err = r.finalizeNifiRegistryClient(reqLogger, registryClient, cluster); err != nil {
			return RequeueWithError(reqLogger, "failed to finalize nifiregistryclient", err)
		}
		if err = r.removeFinalizer(ctx, registryClient); err != nil {
			return RequeueWithError(reqLogger, "failed to remove finalizer from nifiregistryclient", err)
		}
	}
	return Reconciled()
}

func (r *NifiRegistryClientReconciler) removeFinalizer(ctx context.Context, registryClient *v1alpha1.NifiRegistryClient) error {
	registryClient.SetFinalizers(util.StringSliceRemove(registryClient.GetFinalizers(), registryClientFinalizer))
	_, err := r.updateAndFetchLatest(ctx, registryClient)
	return err
}

func (r *NifiRegistryClientReconciler) finalizeNifiRegistryClient(reqLogger logr.Logger, registryClient *v1alpha1.NifiRegistryClient,
	cluster *v1alpha1.NifiCluster) error {

	if err := registryclient.RemoveRegistryClient(r.Client, registryClient, cluster); err != nil {
		return err
	}
	reqLogger.Info("Delete Registry client")

	return nil
}
