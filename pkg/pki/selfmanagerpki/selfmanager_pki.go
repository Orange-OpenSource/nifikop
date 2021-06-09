package selfmanagerpki

import (
	"context"
	"fmt"
	"github.com/Orange-OpenSource/nifikop/api/v1alpha1"
	"github.com/Orange-OpenSource/nifikop/pkg/errorfactory"
	certutil "github.com/Orange-OpenSource/nifikop/pkg/util/cert"
	pkicommon "github.com/Orange-OpenSource/nifikop/pkg/util/pki"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"time"
)

func (s *SelfManager) ReconcilePKI(ctx context.Context, logger logr.Logger, scheme *runtime.Scheme, externalHostnames []string) error {
	logger.Info("Reconciling selfmanager PKI")

	resources, err := s.fullPKI(s.cluster, scheme, externalHostnames)
	if err != nil {
		return err
	}

	for _, o := range resources {
		if err := reconcile(ctx, logger, s.client, o, s.cluster); err != nil {
			return err
		}
	}

	return nil
}

func (s SelfManager) FinalizePKI(ctx context.Context, logger logr.Logger) error {
	logger.Info("Removing selfmanager certificates and secrets")

	// Safety check that we are actually doing something
	if s.cluster.Spec.ListenersConfig.SSLSecrets == nil {
		return nil
	}

	// Names of our secrets
	var objNames []types.NamespacedName

	// Node secrets
	for _, node := range s.cluster.Spec.Nodes {
		objNames = append(objNames, types.NamespacedName{Name: fmt.Sprintf(pkicommon.NodeServerCertTemplate, s.cluster.Name, node.Id), Namespace: s.cluster.Namespace})
	}

	// Controller cert
	objNames = append(
		objNames,
		types.NamespacedName{Name: fmt.Sprintf(pkicommon.NodeControllerTemplate, s.cluster.Name), Namespace: s.cluster.Namespace})

	for _, obj := range objNames {
		// Delete all secrets
		secret := &corev1.Secret{}
		if err := s.client.Get(ctx, obj, secret); err != nil {
			if apierrors.IsNotFound(err) {
				continue
			} else {
				return err
			}
		}
		if err := s.client.Delete(ctx, secret); err != nil {
			return err
		}
	}

	return nil
}

// Return the list of all objects needed for PKI
func (s *SelfManager) fullPKI(cluster *v1alpha1.NifiCluster, scheme *runtime.Scheme, externalHostnames []string) ([]runtime.Object, error) {
	var objects []runtime.Object

	// Ca cert
	caSecret, err := s.caCertForCluster(cluster, scheme)
	if err != nil {
		return objects, err
	}
	objects = append(objects, caSecret)

	// Controller cert
	controllerSecret, err := s.clusterSecretForController()
	if err != nil {
		return objects, err
	}
	objects = append(objects, controllerSecret)

	objects = append(objects, pkicommon.ControllerUserForCluster(cluster))
	// Node "users"
	for _, user := range pkicommon.NodeUsersForCluster(cluster, externalHostnames) {
		// User
		objects = append(objects, user)

		// User's cert
		userSecret, err := s.clusterSecretForUser(user, scheme)
		if err != nil {
			return objects, err
		}
		objects = append(objects, userSecret)
	}
	return objects, nil
}

// Return the 'ca-certificate" secret
func (s *SelfManager) caCertForCluster(cluster *v1alpha1.NifiCluster, scheme *runtime.Scheme) (*corev1.Secret, error) {
	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf(pkicommon.NodeCACertTemplate, cluster.Name),
			Namespace: cluster.Namespace,
			Labels:    pkicommon.LabelsForNifiPKI(cluster.Name),
		},
		Data: map[string][]byte{
			v1alpha1.CoreCACertKey:  s.caCertPEM,
			corev1.TLSCertKey:       s.caCertPEM,
			corev1.TLSPrivateKeyKey: s.caKeyPEM,
			v1alpha1.CertValidity:   []byte(s.caCert.NotAfter.Format(time.UnixDate)),
		},
		Type: corev1.SecretTypeTLS,
	}, nil
}

// Return a secret for specified User
func (s *SelfManager) clusterSecretForUser(user *v1alpha1.NifiUser, scheme *runtime.Scheme) (secret *corev1.Secret, err error) {

	cert, certPEM, keyPEM, err := s.generateUserCert(user)
	if err != nil {
		return nil, err
	}

	secret = &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      user.Spec.SecretName,
			Namespace: user.GetNamespace(),
		},
		Data: map[string][]byte{
			v1alpha1.CoreCACertKey:  s.caCertPEM,
			corev1.TLSCertKey:       certPEM,
			corev1.TLSPrivateKeyKey: keyPEM,
			v1alpha1.CertValidity:   []byte(cert.NotAfter.Format(time.UnixDate)),
		},
	}

	if user.Spec.IncludeJKS {
		secret, err = certutil.EnsureSecretPassJKS(secret)
		if err != nil {
			return
		}

		keystore, truststore, err := s.generateJKSstores(secret.Data[v1alpha1.PasswordKey], certPEM, keyPEM)
		if err != nil {
			return secret, err
		}

		secret.Data[v1alpha1.TLSJKSKeyStore] = []byte(keystore)
		secret.Data[v1alpha1.TLSJKSTrustStore] = []byte(truststore)
	}

	controllerutil.SetControllerReference(user, secret, scheme)
	return
}

// Return  a special secret for the 'controller'
func (s *SelfManager) clusterSecretForController() (secret *corev1.Secret, err error) {

	cert, certPEM, keyPEM, err := s.generateControllerCertPEM()
	if err != nil {
		return nil, err
	}

	secret = &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf(pkicommon.NodeControllerTemplate, s.cluster.Name),
			Namespace: s.cluster.Namespace,
		},
		Data: map[string][]byte{
			v1alpha1.CoreCACertKey:  s.caCertPEM,
			corev1.TLSCertKey:       certPEM,
			corev1.TLSPrivateKeyKey: keyPEM,
			v1alpha1.CertValidity:   []byte(cert.NotAfter.Format(time.UnixDate)),
		},
	}

	return
}

func caValuesFromSecretCert(ctx context.Context, client client.Client, cluster *v1alpha1.NifiCluster) (caCert []byte, caKey []byte, err error) {
	secret := &corev1.Secret{}
	var name = fmt.Sprintf(pkicommon.NodeCACertTemplate, cluster.Name)
	err = client.Get(ctx, types.NamespacedName{Namespace: cluster.Namespace, Name: name}, secret)
	if err != nil {
		if apierrors.IsNotFound(err) {
			err = errorfactory.New(errorfactory.ResourceNotReady{}, err, "could not find provided tls secret")
		} else {
			err = errorfactory.New(errorfactory.APIFailure{}, err, "could not lookup provided tls secret")
		}
		return
	}

	caCert = secret.Data[v1alpha1.CoreCACertKey]
	caKey = secret.Data[v1alpha1.TLSKey]
	return
}
