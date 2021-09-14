package config

import (
	"github.com/Orange-OpenSource/nifikop/api/v1alpha1"
	"github.com/Orange-OpenSource/nifikop/pkg/nificlient/config/nificluster"
	"github.com/Orange-OpenSource/nifikop/pkg/nificlient/config/tls"
	"github.com/Orange-OpenSource/nifikop/pkg/util/clientconfig"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var MockClientConfig = v1alpha1.ClientConfigType("mock")

func GetClientConfigManager(client client.Client, clusterRef v1alpha1.ClusterReference) clientconfig.Manager {
	switch clusterRef.Type {
	case v1alpha1.ClientConfigNiFiCluster:
		return nificluster.New(client, clusterRef)
	case v1alpha1.ClientConfigExternalTLS:
		return tls.New(client, clusterRef)
	//case v1alpha1.ClientConfigExternalBasic:
	//	return
	case MockClientConfig:
		return newMockClientConfig(client, clusterRef)
	default:
		return nificluster.New(client, clusterRef)
	}
}

// Mock types and functions
type mockClientConfig struct {
	clientconfig.Manager
	client     client.Client
	clusterRef v1alpha1.ClusterReference
}

func newMockClientConfig(client client.Client, clusterRef v1alpha1.ClusterReference) clientconfig.Manager {
	return &mockClientConfig{client: client, clusterRef: clusterRef}
}

func (n *mockClientConfig) BuildConfig() (*clientconfig.NifiConfig, error) {
	return nil, nil
}

func (n *mockClientConfig) BuildConnect() (cluster clientconfig.ClusterConnect, err error) {
	return
}

func (n *mockClientConfig) IsExternal() bool {
	return true
}
