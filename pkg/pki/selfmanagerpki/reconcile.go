package selfmanagerpki

import (
	"context"
	"emperror.dev/errors"
	"fmt"
	"github.com/Orange-OpenSource/nifikop/api/v1alpha1"
	pkicommon "github.com/Orange-OpenSource/nifikop/pkg/util/pki"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"time"
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
	fmt.Printf("Checking validity of %s\n", secret.Name)
	if checkCertValidity(obj) != true {
		// Update this cert...
		fmt.Printf("Cert %s is expiring in less than 1 hour. Starting renewal of this cert.\n", secret.Name)

		// Delete the secret to be recreated
		if err = client.Delete(ctx, secret); err != nil {
			return err
		}

		// If the secret is the CA Cert, return "renewal" error for a complete recreation
		if secret.Name == fmt.Sprint(pkicommon.NodeCACertTemplate, cluster.Name) {
			return errors.New("renewal")
		}
	}
	return nil
}

func checkCertValidity(obj *corev1.Secret) bool {
	// Parse date from validity data
	validity, err := time.Parse(time.UnixDate, string(obj.Data[v1alpha1.CertValidity]))
	if err != nil {
		return false
	}

	// Check if the cert will be outdated in 1 hour
	if time.Now().Add(time.Hour * -1).After(validity) {
		return false
	}

	return true
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
