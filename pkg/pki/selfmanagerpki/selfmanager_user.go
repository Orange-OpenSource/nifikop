package selfmanagerpki

import (
	"context"
	"errors"
	"github.com/Orange-OpenSource/nifikop/api/v1alpha1"
	"github.com/Orange-OpenSource/nifikop/pkg/errorfactory"
	"github.com/Orange-OpenSource/nifikop/pkg/k8sutil"
	certutil "github.com/Orange-OpenSource/nifikop/pkg/util/cert"
	pkicommon "github.com/Orange-OpenSource/nifikop/pkg/util/pki"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (s *SelfManager) ReconcileUserCertificate(ctx context.Context, user *v1alpha1.NifiUser, scheme *runtime.Scheme) (*pkicommon.UserCertificate, error) {
	var err error
	var secret *corev1.Secret

	// Check if there is already a secret
	secret, err = s.getUserSecret(ctx, user)

	if err != nil && apierrors.IsNotFound(err) {
		// No secret found, generate & create one

		if user.Spec.IncludeJKS {
			if err := s.injectJKSPassword(ctx, user); err != nil {
				return nil, err
			}
		}
		secret, err = s.clusterSecretForUser(user, scheme)
		if err != nil {
			return nil, errorfactory.New(errorfactory.APIFailure{}, err, "error while generating user secret")
		}

		if err = s.client.Create(ctx, secret); err != nil {
			return nil, errorfactory.New(errorfactory.APIFailure{}, err, "could not create user secret")
		}

	} else if err != nil {
		// API failure, requeue
		return nil, errorfactory.New(errorfactory.APIFailure{}, err, "failed looking up user secret")
	}

	// Ensure controller reference on user secret
	if err = s.ensureControllerReference(ctx, user, secret, scheme); err != nil {
		return nil, err
	}

	return &pkicommon.UserCertificate{
		CA:          secret.Data[v1alpha1.CoreCACertKey],
		Certificate: secret.Data[corev1.TLSCertKey],
		Key:         secret.Data[corev1.TLSPrivateKeyKey],
	}, nil
}

func (s *SelfManager) FinalizeUserCertificate(ctx context.Context, user *v1alpha1.NifiUser) error {
	// TODO delete user secret ?

	return nil
}

// getUserSecret fetches the secret created for a user
func (s *SelfManager) getUserSecret(ctx context.Context, user *v1alpha1.NifiUser) (secret *corev1.Secret, err error) {
	secret = &corev1.Secret{}
	err = s.client.Get(ctx, types.NamespacedName{Name: user.Spec.SecretName, Namespace: user.Namespace}, secret)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return secret, errorfactory.New(errorfactory.ResourceNotReady{}, err, "user secret not ready")
		}
		return secret, errorfactory.New(errorfactory.APIFailure{}, err, "failed to get user secret")
	}
	if user.Spec.IncludeJKS {
		if len(secret.Data) != 6 {
			return secret, errorfactory.New(errorfactory.ResourceNotReady{}, err, "user secret not populated yet")
		}
	} else {
		if len(secret.Data) != 3 {
			return secret, errorfactory.New(errorfactory.ResourceNotReady{}, err, "user secret not populated yet")
		}
	}

	for _, v := range secret.Data {
		if len(v) == 0 {
			return secret, errorfactory.New(errorfactory.ResourceNotReady{},
				errors.New("not all secret value populated"), "secret is not ready")
		}
	}

	return secret, nil
}

// ensureControllerReference ensures that a NifiUser owns a given Secret
func (s *SelfManager) ensureControllerReference(ctx context.Context, user *v1alpha1.NifiUser, secret *corev1.Secret, scheme *runtime.Scheme) error {
	err := controllerutil.SetControllerReference(user, secret, scheme)
	if err != nil && !k8sutil.IsAlreadyOwnedError(err) {
		return errorfactory.New(errorfactory.InternalError{}, err, "error checking controller reference on user secret")
	} else if err == nil {
		if err = s.client.Update(ctx, secret); err != nil {
			return errorfactory.New(errorfactory.APIFailure{}, err, "could not update secret with controller reference")
		}
	}
	return nil
}

// injectJKSPassword ensures that a secret contains JKS password when requested
func (s *SelfManager) injectJKSPassword(ctx context.Context, user *v1alpha1.NifiUser) error {
	var err error
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      user.Spec.SecretName,
			Namespace: user.Namespace,
		},
		Data: map[string][]byte{},
	}
	secret, err = certutil.EnsureSecretPassJKS(secret)
	if err != nil {
		return errorfactory.New(errorfactory.InternalError{}, err, "could not inject secret with jks password")
	}
	if err = s.client.Create(ctx, secret); err != nil {
		return errorfactory.New(errorfactory.APIFailure{}, err, "could not create secret with jks password")
	}

	return nil
}
