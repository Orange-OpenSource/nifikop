package selfmanagerpki

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/Orange-OpenSource/nifikop/pkg/errorfactory"
	pkicommon "github.com/Orange-OpenSource/nifikop/pkg/util/pki"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"strings"
)

// GetControllerTLSConfig creates a TLS config from the user secret created for
// cruise control and manager operations
func (s *SelfManager) GetControllerTLSConfig() (config *tls.Config, err error) {
	config = &tls.Config{}
	tlsKeys := &corev1.Secret{}
	err = s.client.Get(context.TODO(),
		types.NamespacedName{
			Namespace: s.cluster.Namespace,
			Name:      fmt.Sprintf(pkicommon.NodeControllerTemplate, s.cluster.Name),
		},
		tlsKeys,
	)

	if err != nil {
		if apierrors.IsNotFound(err) {
			tlsKeys, err = s.clusterSecretForController()
			if err != nil {
				return
			}
		}
		return
	}
	clientCert := tlsKeys.Data[corev1.TLSCertKey]
	clientKey := tlsKeys.Data[corev1.TLSPrivateKeyKey]
	caCert := s.caCertPEM

	if len(caCert) == 0 {
		certs := strings.SplitAfter(string(clientCert), "-----END CERTIFICATE-----")
		clientCert = []byte(certs[0])
		caCert = []byte(certs[len(certs)-1])
		if len(certs) == 3 {
			caCert = []byte(certs[len(certs)-2])
		}
	}

	x509ClientCert, err := tls.X509KeyPair(clientCert, clientKey)
	if err != nil {
		err = errorfactory.New(errorfactory.InternalError{}, err, "could not decode controller certificate")
		return
	}

	rootCAs := x509.NewCertPool()
	rootCAs.AppendCertsFromPEM(caCert)

	config.Certificates = []tls.Certificate{x509ClientCert}
	config.RootCAs = rootCAs
	//config.InsecureSkipVerify = true // TODO test & find workaround
	return
}
