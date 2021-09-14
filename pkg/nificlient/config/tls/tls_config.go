package tls

import (
	"fmt"
	"github.com/Orange-OpenSource/nifikop/api/v1alpha1"
	"github.com/Orange-OpenSource/nifikop/pkg/common"
	"github.com/Orange-OpenSource/nifikop/pkg/nificlient"
	"github.com/Orange-OpenSource/nifikop/pkg/pki/certmanagerpki"
	"github.com/Orange-OpenSource/nifikop/pkg/util/clientconfig"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var log = ctrl.Log.WithName("tls_config")

func (n *tls) BuildConfig() (*clientconfig.NifiConfig, error) {
	return clusterConfig(n.client, n.clusterRef)
}

func (n *tls) BuildConnect() (cluster clientconfig.ClusterConnect, err error) {
	config, err := n.BuildConfig()
	cluster = &ExternalTLSCluster{
		NodeURITemplate:    n.clusterRef.NodeURITemplate,
		NodeIds:            n.clusterRef.NodeIds,
		NifiURI:             n.clusterRef.NifiURI,
		RootProcessGroupId: n.clusterRef.RootProcessGroupId,

		nifiConfig: config,
	}
	return
}

func (n *tls) IsExternal() bool {
	return false
}

func clusterConfig(client client.Client, ref v1alpha1.ClusterReference) (*clientconfig.NifiConfig, error) {
	nodesURI := generateNodesAddressFromTemplate(ref.NodeIds, ref.NodeURITemplate)

	conf := &clientconfig.NifiConfig{}
	conf.RootProcessGroupId = ref.RootProcessGroupId
	conf.NodeURITemplate = ref.NodeURITemplate
	conf.NodesURI = nodesURI
	conf.NifiURI = ref.NifiURI
	conf.OperationTimeout = clientconfig.NifiDefaultTimeout

	tlsConfig, err := certmanagerpki.GetControllerTLSConfigFromSecret(client, ref.SecretRef)
	if err != nil {
		return conf, err
	}
	conf.UseSSL = true
	conf.TLSConfig = tlsConfig
	return conf, nil
}

func generateNodesAddressFromTemplate(ids []int32, template string) map[int32]clientconfig.NodeUri {
	addresses := make(map[int32]clientconfig.NodeUri)

	for _,nId := range ids {
		addresses[nId] = clientconfig.NodeUri{
			HostListener: fmt.Sprintf(template ,nId),
			RequestHost:  fmt.Sprintf(template ,nId),
		}
	}
	return addresses
}

type ExternalTLSCluster struct {
	NodeURITemplate    string
	NodeIds            []int32
	NifiURI             string
	RootProcessGroupId string

	nifiConfig *clientconfig.NifiConfig
}

func (cluster *ExternalTLSCluster) IsExternal() bool {
	return true
}

func (cluster *ExternalTLSCluster) IsInternal() bool {
	return false
}

func (cluster *ExternalTLSCluster) ClusterLabelString() string {
	return fmt.Sprintf("%s", cluster.NifiURI)
}

func (cluster ExternalTLSCluster) IsReady() bool {
	nClient, err := common.NewClusterConnection(log, cluster.nifiConfig)
	if err != nil {
		return false
	}

	clusterEntity, err := nClient.DescribeCluster()
	if err != nil {
		return false
	}

	for _, node := range clusterEntity.Cluster.Nodes{
		if node.Status != nificlient.CONNECTED_STATUS {
			return false
		}
	}
	return true
}

func (cluster *ExternalTLSCluster) Id() string {
	return cluster.NifiURI
}