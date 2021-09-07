package config

import (
	"github.com/Orange-OpenSource/nifikop/api/v1alpha1"
	"github.com/Orange-OpenSource/nifikop/pkg/nificlient/config/nificluster"
	"github.com/Orange-OpenSource/nifikop/pkg/util/clientconfig"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var MockClientConfig = v1alpha1.ClientConfigType("mock")

func GetClientConfigManager(client client.Client, clusterRef v1alpha1.ClusterReference) clientconfig.Manager {
	switch clusterRef.Type {
	case v1alpha1.ClientConfigNiFiCluster:
		return nificluster.New(client, clusterRef)
	//case v1alpha1.ClientConfigExternalTLS:
	//	return
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

//// external
//func ExternalTLSConfig(ref v1alpha1.ClusterReference) (*NifiConfig, error) {
//	nodesURI := generateNodesAddressFromTemplate(ref.NodeIds, ref.NodeURITemplate)
//
//	conf := &NifiConfig{}
//	conf.RootProcessGroupId = ref.RootProcessGroupId
//	conf.NodeURITemplate = ref.NodeURITemplate
//	conf.NodesURI = nodesURI
//	conf.NifiURI = ref.NifiURI
//	conf.OperationTimeout = nifiDefaultTimeout
//
//	tlsConfig, err := certmanagerpki.GetControllerTLSConfigFromSecret()
//	if err != nil {
//		return conf, err
//	}
//	conf.UseSSL = true
//	conf.TLSConfig = tlsConfig
//	return conf, nil
//}
//
//
//func generateNodesAddressFromTemplate(ids []int32, template string) map[int32]nodeUri {
//	addresses := make(map[int32]nodeUri)
//
//	for _,nId := range ids {
//		addresses[nId] = nodeUri{
//			HostListener: fmt.Sprintf(template ,nId),
//			RequestHost:  fmt.Sprintf(template ,nId),
//		}
//	}
//	return addresses
//}
