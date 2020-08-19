package dataflow

import (
	"github.com/Orange-OpenSource/nifikop/pkg/apis/nifi/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func DataflowExist(client client.Client, flow *v1alpha1.NifiDataflow, cluster *v1alpha1.NifiCluster) (bool, error){

	if flow.Status.ProcessGroupID == "" {
		return false, nil
	}

	return true, nil
}

func CreateDataflow(client client.Client, flow *v1alpha1.NifiDataflow, cluster *v1alpha1.NifiCluster) (string, error) {
	return "", nil
}

func RunDataflow(client client.Client, flow *v1alpha1.NifiDataflow, cluster *v1alpha1.NifiCluster) error {
	return nil
}

func DeleteDataflow(client client.Client, flow *v1alpha1.NifiDataflow, cluster *v1alpha1.NifiCluster) error {
	return nil
}

