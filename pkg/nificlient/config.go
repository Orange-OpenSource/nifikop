package nificlient

import (
	"crypto/tls"
	"fmt"
	"github.com/Orange-OpenSource/nifikop/api/v1alpha1"
	"github.com/Orange-OpenSource/nifikop/pkg/pki"
	"github.com/Orange-OpenSource/nifikop/pkg/util"
	"github.com/Orange-OpenSource/nifikop/pkg/util/nifi"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strings"
)

const (
	nifiDefaultTimeout = int64(5)
)

// NifiConfig are the options to creating a new ClusterAdmin client
type NifiConfig struct {
	nodeURITemplate string
	NodesURI        map[int32]nodeUri
	NifiURI          string
	UseSSL          bool
	TLSConfig        *tls.Config

	OperationTimeout int64
	RootProcessGroupId string
}

type nodeUri struct {
	HostListener string
	RequestHost  string
}

// ClusterConfig creates connection options from a NifiCluster CR
func ClusterConfig(client client.Client, cluster *v1alpha1.NifiCluster) (*NifiConfig, error) {

	conf := &NifiConfig{}
	conf.RootProcessGroupId = cluster.Status.RootProcessGroupId
	conf.nodeURITemplate = generateNodesURITemplate(cluster)
	conf.NodesURI = generateNodesAddress(cluster)
	conf.NifiURI = nifi.GenerateRequestNiFiAllNodeAddressFromCluster(cluster)
	conf.OperationTimeout = nifiDefaultTimeout

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

func generateNodesAddress(cluster *v1alpha1.NifiCluster) map[int32]nodeUri {
	addresses := make(map[int32]nodeUri)

	for nId, state := range cluster.Status.NodesState {
		if !(state.GracefulActionState.State.IsRunningState() || state.GracefulActionState.State.IsRequiredState()) && state.GracefulActionState.ActionStep != v1alpha1.RemoveStatus {
			addresses[util.ConvertStringToInt32(nId)] = nodeUri{
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