package selfmanagerpki

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"github.com/Orange-OpenSource/nifikop/api/v1alpha1"
	"github.com/Orange-OpenSource/nifikop/pkg/util/pki"
	pkicommon "github.com/Orange-OpenSource/nifikop/pkg/util/pki"
	"math/big"
	"net/url"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"time"
)

var subject = pkix.Name{
	Organization:  []string{"Orange"},
	Country:       []string{"FR"},
	Province:      []string{""},
	Locality:      []string{"Paris"},
	StreetAddress: []string{"78 Rue Olivier de Serres"},
	PostalCode:    []string{"75015"},
}

type SelfManager struct {
	pki.Manager
	client  client.Client
	cluster *v1alpha1.NifiCluster

	caCert    *x509.Certificate
	caCertPEM []byte
	caKey     *rsa.PrivateKey
	caKeyPEM  []byte
}

// Return a new fully instantiated SelfManager struct
func New(client client.Client, cluster *v1alpha1.NifiCluster) (manager *SelfManager, err error) {
	manager = &SelfManager{
		client:  client,
		cluster: cluster,
	}

	// setting up our ca and server certificate
	if err = manager.setupCA(); err != nil {
		return
	}

	return
}

// Sets up the caCert & caKey variables by setting up a new self signed CA
func (s *SelfManager) setupCA() (err error) {
	// set up our CA certificate
	s.caCert = &x509.Certificate{
		SerialNumber:          big.NewInt(2019),
		Subject:               subject,
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}

	// create our private and public key
	if s.caKey, err = rsa.GenerateKey(rand.Reader, 4096); err != nil {
		return
	}

	// create the CA
	caBytes, err := x509.CreateCertificate(rand.Reader, s.caCert, s.caCert, &s.caKey.PublicKey, s.caKey)
	if err != nil {
		return
	}

	// pem encode
	caPEM := new(bytes.Buffer)
	if err = pem.Encode(caPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caBytes,
	}); err != nil {
		return
	}
	s.caCertPEM = caPEM.Bytes()

	caPrivKeyPEM := new(bytes.Buffer)
	if err = pem.Encode(caPrivKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(s.caKey),
	}); err != nil {
		return
	}
	s.caKeyPEM = caPrivKeyPEM.Bytes()

	return
}

func (s *SelfManager) generateUserCert(user *v1alpha1.NifiUser) (certPEM []byte, certPrivKeyPEM []byte, err error) {
	// set up our server certificate
	cert := &x509.Certificate{
		SerialNumber: big.NewInt(2019),
		Subject:      subject,
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(10, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}

	// URLs
	urls, err := url.Parse(fmt.Sprintf(pkicommon.SpiffeIdTemplate, s.cluster.Name, user.GetNamespace(), user.GetName()))
	if err != nil {
		return
	}
	cert.URIs = []*url.URL{urls}

	// Add DNS Names if provided
	if user.Spec.DNSNames != nil && len(user.Spec.DNSNames) > 0 {
		cert.DNSNames = user.Spec.DNSNames
	}

	certPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, cert, s.caCert, &certPrivKey.PublicKey, s.caKey)
	if err != nil {
		return
	}

	certPEMBuffer := new(bytes.Buffer)
	if err = pem.Encode(certPEMBuffer, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	}); err != nil {
		return
	}
	certPEM = certPEMBuffer.Bytes()

	certPrivKeyPEMBuffer := new(bytes.Buffer)
	if err = pem.Encode(certPrivKeyPEMBuffer, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(certPrivKey),
	}); err != nil {
		return
	}
	certPrivKeyPEM = certPrivKeyPEMBuffer.Bytes()

	return
}

// Generate one cert from selfmanager's CA
func (s *SelfManager) generateCaCertPEM() (certPEM []byte, certPrivKeyPEM []byte, err error) {
	// set up our server certificate
	cert := &x509.Certificate{
		SerialNumber: big.NewInt(2019),
		Subject:      subject,
		//IPAddresses:  []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(10, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}

	certPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, cert, s.caCert, &certPrivKey.PublicKey, s.caKey)
	if err != nil {
		return
	}

	certPEMBuffer := new(bytes.Buffer)
	if err = pem.Encode(certPEMBuffer, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	}); err != nil {
		return
	}
	certPEM = certPEMBuffer.Bytes()

	certPrivKeyPEMBuffer := new(bytes.Buffer)
	if err = pem.Encode(certPrivKeyPEMBuffer, &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(certPrivKey),
	}); err != nil {
		return
	}
	certPrivKeyPEM = certPrivKeyPEMBuffer.Bytes()

	return
}
