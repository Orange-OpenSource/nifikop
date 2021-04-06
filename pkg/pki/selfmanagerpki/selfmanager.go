package selfmanagerpki

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"github.com/Orange-OpenSource/nifikop/api/v1alpha1"
	"github.com/Orange-OpenSource/nifikop/pkg/util/pki"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"math/big"
	"net"
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

	// TODO PEMs or objects ?
	caCert     *x509.Certificate
	caCertByte []byte
	caKey      *rsa.PrivateKey
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
	s.caCertByte, err = x509.CreateCertificate(rand.Reader, s.caCert, s.caCert, &s.caKey.PublicKey, s.caKey)
	if err != nil {
		return
	}

	// pem encode
	caPEM := new(bytes.Buffer)
	if err = pem.Encode(caPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: s.caCertByte,
	}); err != nil {
		return
	}

	caPrivKeyPEM := new(bytes.Buffer)
	if err = pem.Encode(caPrivKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(s.caKey),
	}); err != nil {
		return
	}

	return
}

// TODO PEM or Bytes + params ?
// Generate one cert from selfmanager's CA
func (s *SelfManager) generateUserCert(user *v1alpha1.NifiUser) (certPEM *bytes.Buffer, certPrivKeyPEM *bytes.Buffer, err error) {
	// set up our server certificate
	cert := &x509.Certificate{
		SerialNumber: big.NewInt(2019),
		Subject:      subject,
		IPAddresses:  []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(10, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}

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

	certPEM = new(bytes.Buffer)
	if err = pem.Encode(certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	}); err != nil {
		return
	}

	certPrivKeyPEM = new(bytes.Buffer)
	if err = pem.Encode(certPrivKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(certPrivKey),
	}); err != nil {
		return
	}

	return
}

func (s *SelfManager) clusterSecretForUser(user *v1alpha1.NifiUser, scheme *runtime.Scheme) (secret *corev1.Secret, err error) {

	cert, keyPEM, err := s.generateUserCert(user)
	if err != nil {
		return nil, err
	}

	secret = &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      user.GetName(),
			Namespace: user.GetNamespace(),
		},
		Data: map[string][]byte{
			// TODO wich types ?
			v1alpha1.CoreCACertKey: s.caCertByte,
			//v1alpha1.CACertKey:  s.caCert,
			corev1.TLSCertKey:       cert.Bytes(),
			corev1.TLSPrivateKeyKey: keyPEM.Bytes(),
		},
	}

	// TODO what to do for jks ?
	//if user.Spec.IncludeJKS {
	//	secret.Data.Keystores = &certv1.CertificateKeystores{
	//		JKS: &certv1.JKSKeystore{
	//			Create: true,
	//			PasswordSecretRef: certmeta.SecretKeySelector{
	//				LocalObjectReference: certmeta.LocalObjectReference{
	//					Name: user.Spec.SecretName,
	//				},
	//				Key: v1alpha1.PasswordKey,
	//			},
	//		},
	//	}
	//}

	// TODO ??
	//controllerutil.SetControllerReference(user, secret, scheme)
	return
}
