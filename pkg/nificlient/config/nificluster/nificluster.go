package nificluster

import (
	"github.com/Orange-OpenSource/nifikop/api/v1alpha1"
	"github.com/Orange-OpenSource/nifikop/pkg/util/clientconfig"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type NifiCluster interface {
	clientconfig.Manager
}

type nifiCluster struct {
	client     client.Client
	clusterRef v1alpha1.ClusterReference
}

func New(client client.Client, clusterRef v1alpha1.ClusterReference) NifiCluster {
	return &nifiCluster{clusterRef: clusterRef, client: client}
}
