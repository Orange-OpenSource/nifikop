package selfmanagerpki

import (
	"context"
	"fmt"
	"github.com/Orange-OpenSource/nifikop/api/v1alpha1"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// reconcile ensures the given kubernetes object
func reconcile(ctx context.Context, log logr.Logger, client client.Client, object runtime.Object, cluster *v1alpha1.NifiCluster) (err error) {
	switch object.(type) {
	case *corev1.Secret:
		secret, _ := object.(*corev1.Secret)
		return reconcileSecret(ctx, log, client, secret, cluster)
	case *v1alpha1.NifiUser:
		user, _ := object.(*v1alpha1.NifiUser)
		return reconcileUser(ctx, log, client, user, cluster)
	default:
		panic(fmt.Sprintf("Invalid object type: %v", reflect.TypeOf(object)))
	}
}

// reconcileSecret ensures a Kubernetes secret
func reconcileSecret(ctx context.Context, log logr.Logger, client client.Client, secret *corev1.Secret, cluster *v1alpha1.NifiCluster) error {
	obj := &corev1.Secret{}
	var err error
	if err = client.Get(ctx, types.NamespacedName{Name: secret.Name, Namespace: secret.Namespace}, obj); err != nil {
		if !apierrors.IsNotFound(err) {
			return err
		}
		return client.Create(ctx, secret)
	}
	return nil
}

// reconcileUser ensures a v1alpha1.NifiUser
func reconcileUser(ctx context.Context, log logr.Logger, client client.Client, user *v1alpha1.NifiUser, cluster *v1alpha1.NifiCluster) error {
	obj := &v1alpha1.NifiUser{}
	var err error
	if err = client.Get(ctx, types.NamespacedName{Name: user.Name, Namespace: user.Namespace}, obj); err != nil {
		if !apierrors.IsNotFound(err) {
			return err
		}
		return client.Create(ctx, user)
	}
	return nil
}
