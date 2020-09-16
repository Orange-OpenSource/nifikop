// Copyright 2020 Orange SA
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.package apis

package nificlient

import (
	"testing"

	nigoapi "github.com/erdrix/nigoapi/pkg/nifi"
	"github.com/jarcoal/httpmock"
	"github.com/Orange-OpenSource/nifikop/pkg/apis/nifi/v1alpha1"
	nifiutil "github.com/Orange-OpenSource/nifikop/pkg/util/nifi"
)

var (
	nodesId = map[int32]string{0: "12334456", 1: "12334456", 2: "12334456"}
)

type mockNiFiClient struct {
	NifiClient
	opts       *NifiConfig
	client     *nigoapi.APIClient
	nodeClient map[int32]*nigoapi.APIClient
	nodes      []nigoapi.NodeDto

	newClient func(*nigoapi.Configuration) *nigoapi.APIClient
	failOpts  bool
}

func newMockOpts() *NifiConfig {
	return &NifiConfig{}
}

func newMockHttpClient(c *nigoapi.Configuration) *nigoapi.APIClient {
	client := nigoapi.NewAPIClient(c)
	httpmock.Activate()
	return client
}

func newMockClient() *nifiClient {
	return &nifiClient{
		opts:      newMockOpts(),
		newClient: newMockHttpClient,
	}
}

func newBuildedMockClient() *nifiClient {
	client := newMockClient()
	client.Build()
	return client
}

func NewMockNiFiClient() *nifiClient {
	return &nifiClient{
		opts:      newMockOpts(),
		newClient: newMockHttpClient,
	}
}

func NewMockNiFiClientFailOps() *mockNiFiClient {
	return &mockNiFiClient{
		opts:      newMockOpts(),
		newClient: newMockHttpClient,
		failOpts:  true,
	}
}

func MockGetClusterResponse(cluster *v1alpha1.NifiCluster) map[string]interface{} {
	return map[string]interface{}{
		"cluster": map[string]interface{}{
			"nodes": []nigoapi.NodeDto{
				{
					NodeId:  nodesId[0],
					Address: nifiutil.ComputeNodeHostnameFromCluster(0, cluster),
					ApiPort: httpContainerPort,
					Status:  string(v1alpha1.ConnectStatus),
				},
				{
					NodeId:  nodesId[1],
					Address: nifiutil.ComputeNodeHostnameFromCluster(1, cluster),
					ApiPort: httpContainerPort,
					Status:  string(v1alpha1.DisconnectStatus),
				},
				{
					NodeId:  nodesId[2],
					Address: nifiutil.ComputeNodeHostnameFromCluster(2, cluster),
					ApiPort: httpContainerPort,
					Status:  string(v1alpha1.OffloadStatus),
				},
			},
		},
	}
}

func MockGetNodeResponse(nodeId int32, cluster *v1alpha1.NifiCluster) interface{} {
	nodes := map[int32]map[string]interface{}{
		0: {
			"node": nigoapi.NodeDto{
				NodeId:  nodesId[0],
				Address: nifiutil.ComputeNodeHostnameFromCluster(0, cluster),
				ApiPort: httpContainerPort,
				Status:  string(v1alpha1.ConnectStatus),
			},
		},
		1: {
			"node": nigoapi.NodeDto{
				NodeId:  nodesId[1],
				Address: nifiutil.ComputeNodeHostnameFromCluster(1, cluster),
				ApiPort: httpContainerPort,
				Status:  string(v1alpha1.ConnectStatus),
			},
		},
		2: {
			"node": nigoapi.NodeDto{
				NodeId:  nodesId[2],
				Address: nifiutil.ComputeNodeHostnameFromCluster(2, cluster),
				ApiPort: httpContainerPort,
				Status:  string(v1alpha1.ConnectStatus),
			},
		},
	}

	return nodes[nodeId]
}

func testClusterMock(t *testing.T) *v1alpha1.NifiCluster {
	t.Helper()
	cluster := &v1alpha1.NifiCluster{}

	cluster.Name = clusterName
	cluster.Namespace = clusterNamespace
	cluster.Spec = v1alpha1.NifiClusterSpec{}

	cluster.Spec.Nodes = []v1alpha1.Node{
		{Id: 0},
		{Id: 1},
		{Id: 2},
	}

	cluster.Spec.ListenersConfig.InternalListeners = []v1alpha1.InternalListenerConfig{
		{Type: "http", ContainerPort: httpContainerPort},
		{Type: "cluster", ContainerPort: 8083},
		{Type: "s2s", ContainerPort: 8085},
	}
	return cluster
}
