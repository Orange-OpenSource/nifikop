package nificluster

import (
	"fmt"
	"github.com/Orange-OpenSource/nifikop/api/v1alpha1"
	"github.com/Orange-OpenSource/nifikop/pkg/k8sutil"
	"github.com/Orange-OpenSource/nifikop/pkg/pki"
	"github.com/Orange-OpenSource/nifikop/pkg/util"
	"github.com/Orange-OpenSource/nifikop/pkg/util/clientconfig"
	"github.com/Orange-OpenSource/nifikop/pkg/util/nifi"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strings"
)

func (n *nifiCluster) BuildConfig() (*clientconfig.NifiConfig, error) {
	var cluster *v1alpha1.NifiCluster
	var err error
	if cluster, err = k8sutil.LookupNifiCluster(n.client, n.clusterRef.Name, n.clusterRef.Namespace); err != nil {
		return nil, err
	}

	return clusterConfig(n.client, cluster)
}

func (n *nifiCluster) BuildConnect() (cluster clientconfig.ClusterConnect, err error) {
	cluster, err = k8sutil.LookupNifiCluster(n.client, n.clusterRef.Name, n.clusterRef.Namespace)
	return
}

func (n *nifiCluster) IsExternal() bool {
	return n.IsExternal()
}

// ClusterConfig creates connection options from a NifiCluster CR
func clusterConfig(client client.Client, cluster *v1alpha1.NifiCluster) (*clientconfig.NifiConfig, error) {
	conf := &clientconfig.NifiConfig{}
	conf.RootProcessGroupId = cluster.Status.RootProcessGroupId
	conf.NodeURITemplate = generateNodesURITemplate(cluster)
	conf.NodesURI = generateNodesAddress(cluster)
	conf.NifiURI = nifi.GenerateRequestNiFiAllNodeAddressFromCluster(cluster)
	conf.OperationTimeout = clientconfig.NifiDefaultTimeout

	if cluster.Spec.ListenersConfig.SSLSecrets != nil && UseSSL(cluster) {
		tlsConfig, err := pki.GetPKIManager(client, cluster).GetControllerTLSConfig()
		if err != nil {
			return conf, err
		}
		conf.UseSSL = true
		conf.TLSConfig = tlsConfig
	}
	return conf, nil
}

func UseSSL(cluster *v1alpha1.NifiCluster) bool {
	return cluster.Spec.ListenersConfig.SSLSecrets != nil
}

func generateNodesAddress(cluster *v1alpha1.NifiCluster) map[int32]clientconfig.NodeUri {
	addresses := make(map[int32]clientconfig.NodeUri)

	for nId, state := range cluster.Status.NodesState {
		if !(state.GracefulActionState.State.IsRunningState() || state.GracefulActionState.State.IsRequiredState()) && state.GracefulActionState.ActionStep != v1alpha1.RemoveStatus {
			addresses[util.ConvertStringToInt32(nId)] = clientconfig.NodeUri{
				HostListener: nifi.GenerateHostListenerNodeAddressFromCluster(util.ConvertStringToInt32(nId), cluster),
				RequestHost:  nifi.GenerateRequestNiFiNodeAddressFromCluster(util.ConvertStringToInt32(nId), cluster),
			}
		}
	}
	return addresses
}

func generateNodesURITemplate(cluster *v1alpha1.NifiCluster) string {
	nodeNameTemplate :=
		fmt.Sprintf(nifi.PrefixNodeNameTemplate, cluster.Name) +
			nifi.RootNodeNameTemplate +
			nifi.SuffixNodeNameTemplate

	return nodeNameTemplate + fmt.Sprintf(".%s",
		strings.SplitAfterN(nifi.GenerateRequestNiFiNodeAddressFromCluster(0, cluster), ".", 2)[1],
	)
}
