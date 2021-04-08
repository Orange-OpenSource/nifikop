package selfmanagerpki

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/Orange-OpenSource/nifikop/api/v1alpha1"
	"github.com/Orange-OpenSource/nifikop/pkg/errorfactory"
	corev1 "k8s.io/api/core/v1"
	"strings"
)

// GetControllerTLSConfig creates a TLS config from the user secret created for
// cruise control and manager operations
func (s *SelfManager) GetControllerTLSConfig() (config *tls.Config, err error) {
	config = &tls.Config{}

	tlsKeys, err := s.clusterSecretForController()
	if err != nil {
		return
	}

	clientCert := tlsKeys.Data[corev1.TLSCertKey]
	clientKey := tlsKeys.Data[corev1.TLSPrivateKeyKey]
	caCert := tlsKeys.Data[v1alpha1.CoreCACertKey]

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

	return
}
