package selfmanagerpki

import (
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
	manager = New()
	manager.SetClientAndCluster(fake.NewFakeClientWithScheme(scheme.Scheme), cluster)
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
	pkiManager := New()
	pkiManager.SetClientAndCluster(&mockClient{}, newMockCluster())

	if reflect.TypeOf(pkiManager) != reflect.TypeOf(&SelfManager{}) {
		t.Error("Expected new selfmanager from New, got:", reflect.TypeOf(pkiManager))
	}
}

func TestGenerateUserCert(t *testing.T) {
	manager := New()
	manager.SetClientAndCluster(&mockClient{}, newMockCluster())

	certPEM, certKeyPEM, err := manager.generateUserCert(newMockUser())
	if err != nil {
		t.Error("Expected no error from generateUserCert, got:", err)
	}
	if certPEM == nil {
		t.Error("Expected caCert not to be nil")
	}
	if certKeyPEM == nil {
		t.Error("Expected cakey not to be nil")
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
	if manager.caCert == nil {
		t.Error("Expected caCert not to be nil")
	}
	if manager.caKey == nil {
		t.Error("Expected cakey not to be nil")
	}
	if manager.caCertPEM == nil {
		t.Error("Expected caCertPEM not to be nil")
	}
	if manager.caKeyPEM == nil {
		t.Error("Expected cakeyPEM not to be nil")
	}
}
