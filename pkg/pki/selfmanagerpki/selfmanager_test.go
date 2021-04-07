package selfmanagerpki

import (
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"github.com/Orange-OpenSource/nifikop/api/v1alpha1"
	certv1 "github.com/jetstack/cert-manager/pkg/apis/certmanager/v1alpha2"
	"k8s.io/client-go/kubernetes/scheme"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"testing"
)

type mockClient struct {
	client.Client
}

func newMock(cluster *v1alpha1.NifiCluster) (manager *SelfManager, err error) {
	certv1.AddToScheme(scheme.Scheme)
	v1alpha1.SchemeBuilder.AddToScheme(scheme.Scheme)
	manager, err = New(fake.NewFakeClientWithScheme(scheme.Scheme), cluster)
	return
}

func newMockCluster() *v1alpha1.NifiCluster {
	cluster := &v1alpha1.NifiCluster{}
	cluster.Name = "test"
	cluster.Namespace = "test-namespace"
	cluster.Spec = v1alpha1.NifiClusterSpec{}
	cluster.Spec.ListenersConfig = v1alpha1.ListenersConfig{}
	cluster.Spec.ListenersConfig.InternalListeners = []v1alpha1.InternalListenerConfig{
		{ContainerPort: 9092},
	}
	cluster.Spec.ListenersConfig.SSLSecrets = &v1alpha1.SSLSecrets{
		TLSSecretName: "test-controller",
		PKIBackend:    v1alpha1.PKIBackendSelfManager,
		Create:        true,
	}

	cluster.Spec.Nodes = []v1alpha1.Node{
		{Id: 0},
		{Id: 1},
		{Id: 2},
	}
	return cluster
}

func TestNew(t *testing.T) {
	pkiManager, err := New(&mockClient{}, newMockCluster())
	if err != nil {
		t.Error("Expected no error from New, got:", err)
	}
	if reflect.TypeOf(pkiManager) != reflect.TypeOf(&SelfManager{}) {
		t.Error("Expected new selfmanager from New, got:", reflect.TypeOf(pkiManager))
	}
}

func TestGenerateUserCert(t *testing.T) {
	manager, err := New(&mockClient{}, newMockCluster())
	if err != nil {
		t.Error("Expected no error from New, got:", err)
	}

	certPEM, certKeyPEM, err := manager.generateUserCert(newMockUser())
	if err != nil {
		t.Error("Expected no error from generateCert, got:", err)
	}
	if reflect.TypeOf(certPEM) != reflect.TypeOf(&bytes.Buffer{}) {
		t.Error("Expected caCert to be bytes.Buffer from setupCA, got:", reflect.TypeOf(certPEM))
	}
	if reflect.TypeOf(certKeyPEM) != reflect.TypeOf(&bytes.Buffer{}) {
		t.Error("Expected cakey to be bytes.Buffer from setupCA, got:", reflect.TypeOf(certKeyPEM))
	}
}

func TestSetupCA(t *testing.T) {
	manager := SelfManager{
		client:  &mockClient{},
		cluster: newMockCluster(),
	}

	if err := manager.setupCA(); err != nil {
		t.Error("Expected no error from setupCA, got:", err)
	}
	if reflect.TypeOf(manager.caCert) != reflect.TypeOf(&x509.Certificate{}) {
		t.Error("Expected caCert to be x509.Certificate from setupCA, got:", reflect.TypeOf(manager.caCert))
	}
	if reflect.TypeOf(manager.caKey) != reflect.TypeOf(&rsa.PrivateKey{}) {
		t.Error("Expected cakey to be rsa.PrivateKey from setupCA, got:", reflect.TypeOf(manager.caKey))
	}
	if reflect.TypeOf(manager.caCertPEM) != reflect.TypeOf(&[]byte{}) {
		t.Error("Expected caCertPEM to be []byte from setupCA, got:", reflect.TypeOf(manager.caCertPEM))
	}
	if reflect.TypeOf(manager.caKeyPEM) != reflect.TypeOf(&[]byte{}) {
		t.Error("Expected cakeyPEM to be []byte from setupCA, got:", reflect.TypeOf(manager.caKeyPEM))
	}
}
