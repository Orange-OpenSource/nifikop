// Copyright 2020 Orange SA
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.package apis

package certmanagerpki

import (
	"context"
	"errors"
	"fmt"

	"github.com/Orange-OpenSource/nifikop/api/v1alpha1"
	"github.com/Orange-OpenSource/nifikop/pkg/errorfactory"
	"github.com/Orange-OpenSource/nifikop/pkg/k8sutil"
	certutil "github.com/Orange-OpenSource/nifikop/pkg/util/cert"
	pkicommon "github.com/Orange-OpenSource/nifikop/pkg/util/pki"
	certv1 "github.com/jetstack/cert-manager/pkg/apis/certmanager/v1alpha2"
	certmeta "github.com/jetstack/cert-manager/pkg/apis/meta/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// FinalizeUserCertificate for cert-manager backend auto returns because controller references handle cleanup
func (c *certManager) FinalizeUserCertificate(ctx context.Context, user *v1alpha1.NifiUser) (err error) {
	return
}

// ReconcileUserCertificate ensures a certificate/secret combination using cert-manager
func (c *certManager) ReconcileUserCertificate(ctx context.Context, user *v1alpha1.NifiUser, scheme *runtime.Scheme) (*pkicommon.UserCertificate, error) {
	var err error
	var secret *corev1.Secret
	// See if we have an existing certificate for this user already
	_, err = c.getUserCertificate(ctx, user)

	if err != nil && apierrors.IsNotFound(err) {
		// the certificate does not exist, let's make one
		cert := c.clusterCertificateForUser(user, scheme)
		if err = c.client.Create(ctx, cert); err != nil {
			return nil, errorfactory.New(errorfactory.APIFailure{}, err, "could not create user certificate")
		}

	} else if err != nil {
		// API failure, requeue
		return nil, errorfactory.New(errorfactory.APIFailure{}, err, "failed looking up user certificate")
	}

	// Get the secret created from the certificate
	secret, err = c.getUserSecret(ctx, user)
	if err != nil {
		return nil, err
	}

	if user.Spec.IncludeJKS {
		if secret, err = c.injectJKS(ctx, user, secret); err != nil {
			return nil, err
		}
	}

	// Ensure controller reference on user secret
	if err = c.ensureControllerReference(ctx, user, secret, scheme); err != nil {
		return nil, err
	}

	// Ensure that the secret is populated with the required values.
	for k, v := range secret.Data {
		if len(v) == 0 && k != v1alpha1.CoreCACertKey {
			return nil, errorfactory.New(errorfactory.APIFailure{}, errors.New("not all secret value populated"), "secret is not ready")
		}
	}

	return &pkicommon.UserCertificate{
		CA:          secret.Data[v1alpha1.CoreCACertKey],
		Certificate: secret.Data[corev1.TLSCertKey],
		Key:         secret.Data[corev1.TLSPrivateKeyKey],
	}, nil
}

// injectJKS ensures that a secret contains JKS format when requested
func (c *certManager) injectJKS(ctx context.Context, user *v1alpha1.NifiUser, secret *corev1.Secret) (*corev1.Secret, error) {
	var err error
	if secret, err = certutil.EnsureSecretJKS(secret); err != nil {
		return nil, errorfactory.New(errorfactory.InternalError{}, err, "could not inject secret with jks")
	}
	if err = c.client.Update(ctx, secret); err != nil {
		return nil, errorfactory.New(errorfactory.APIFailure{}, err, "could not update secret with jks")
	}
	// Fetch the updated secret
	return c.getUserSecret(ctx, user)
}

// ensureControllerReference ensures that a NifiUser owns a given Secret
func (c *certManager) ensureControllerReference(ctx context.Context, user *v1alpha1.NifiUser, secret *corev1.Secret, scheme *runtime.Scheme) error {
	err := controllerutil.SetControllerReference(user, secret, scheme)
	if err != nil && !k8sutil.IsAlreadyOwnedError(err) {
		return errorfactory.New(errorfactory.InternalError{}, err, "error checking controller reference on user secret")
	} else if err == nil {
		if err = c.client.Update(ctx, secret); err != nil {
			return errorfactory.New(errorfactory.APIFailure{}, err, "could not update secret with controller reference")
		}
	}
	return nil
}

// getUserCertificate fetches the cert-manager Certificate for a user
func (c *certManager) getUserCertificate(ctx context.Context, user *v1alpha1.NifiUser) (*certv1.Certificate, error) {
	cert := &certv1.Certificate{}
	err := c.client.Get(ctx, types.NamespacedName{Name: user.Name, Namespace: user.Namespace}, cert)
	return cert, err
}

// getUserSecret fetches the secret created from a cert-manager Certificate for a user
func (c *certManager) getUserSecret(ctx context.Context, user *v1alpha1.NifiUser) (secret *corev1.Secret, err error) {
	secret = &corev1.Secret{}
	err = c.client.Get(ctx, types.NamespacedName{Name: user.Spec.SecretName, Namespace: user.Namespace}, secret)
	if err != nil {
		if apierrors.IsNotFound(err) {
			err = errorfactory.New(errorfactory.ResourceNotReady{}, err, "user secret not ready")
		} else {
			err = errorfactory.New(errorfactory.APIFailure{}, err, "failed to get user secret")
		}
	}
	return
}

// clusterCertificateForUser generates a Certificate object for a NifiUser
func (c *certManager) clusterCertificateForUser(user *v1alpha1.NifiUser, scheme *runtime.Scheme) *certv1.Certificate {
	caName, caKind := c.getCA()
	cert := &certv1.Certificate{
		ObjectMeta: metav1.ObjectMeta{
			Name:      user.Name,
			Namespace: user.Namespace,
		},
		Spec: certv1.CertificateSpec{
			SecretName:  user.Spec.SecretName,
			KeyEncoding: certv1.PKCS8,
			CommonName:  user.Name,
			Usages:      []certv1.KeyUsage{certv1.UsageClientAuth, certv1.UsageServerAuth},
			IssuerRef: certmeta.ObjectReference{
				Name: caName,
				Kind: caKind,
			},
		},
	}
	if user.Spec.DNSNames != nil && len(user.Spec.DNSNames) > 0 {
		cert.Spec.DNSNames = user.Spec.DNSNames
	}
	controllerutil.SetControllerReference(user, cert, scheme)
	return cert
}

// getCA returns the CA name/kind for the NifiCluster
func (c *certManager) getCA() (caName, caKind string) {
	caKind = certv1.IssuerKind
	issuerRef := c.cluster.Spec.ListenersConfig.SSLSecrets.IssuerRef
	if issuerRef != nil {
		caName = issuerRef.Name
		caKind = issuerRef.Kind
	} else {
		if c.cluster.Spec.ListenersConfig.SSLSecrets.ClusterScoped {
			caKind = certv1.ClusterIssuerKind
		}
		caName = fmt.Sprintf(pkicommon.NodeIssuerTemplate, c.cluster.Name)
	}
	// TODO: Do we need to ensure this Issuer is exist?
	return
}
