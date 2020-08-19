package dataflow

import (
	"github.com/Orange-OpenSource/nifikop/pkg/apis/nifi/v1alpha1"
	"github.com/Orange-OpenSource/nifikop/pkg/controller/common"
	"github.com/Orange-OpenSource/nifikop/pkg/nificlient"
	nigoapi "github.com/erdrix/nigoapi/pkg/nifi"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

var log = logf.Log.WithName("dataflow-method")

func DataflowExist(client client.Client, flow *v1alpha1.NifiDataflow, cluster *v1alpha1.NifiCluster) (bool, error){

	if flow.Status.ProcessGroupID == "" {
		return false, nil
	}

	return true, nil
}

func CreateDataflow(client client.Client, flow *v1alpha1.NifiDataflow, cluster *v1alpha1.NifiCluster) (string, error) {
	return "", nil
}

func ScheduleDataflow(client client.Client, flow *v1alpha1.NifiDataflow, cluster *v1alpha1.NifiCluster) error {
	return nil
}

func IsOutOfSyncDataflow(client client.Client, flow *v1alpha1.NifiDataflow, cluster *v1alpha1.NifiCluster) (bool, error) {
	return false, nil
}

func SyncDataflow(client client.Client, flow *v1alpha1.NifiDataflow, cluster *v1alpha1.NifiCluster) error {

	nClient, err := common.NewNodeConnection(log, client, cluster)
	if err != nil {
		return err
	}

	if flow.Spec.UpdateStrategy == v1alpha1.DropStrategy {
		// unschedule processors
		_, err := nClient.UpdateFlowProcessGroup(nigoapi.ScheduleComponentsEntity{
			Id:    flow.Status.ProcessGroupID,
			State: "STOPPED",
		})

		if err == nificlient.ErrNifiClusterNotReturned200 {
			log.Error(err, "Stop flow failed since Nifi node returned non 200")
			return err
		}

		if err != nil {
			log.Error(err, "could not communicate with nifi node")
			return err
		}

		// Get flow
		flowEntity, err := nClient.GetFlow(flow.Status.ProcessGroupID)

		if err == nificlient.ErrNifiClusterNotReturned200 {
			log.Error(err, "Stop flow failed since Nifi node returned non 200")
			return err
		}

		if err != nil {
			log.Error(err, "could not communicate with nifi node")
			return err
		}

		// Drop all events in connections
		for _, connection := range flowEntity.ProcessGroupFlow.Flow.Connections {
			if connection.Status.AggregateSnapshot.FlowFilesQueued != 0 {

			}
		}
	}

	return nil
}

func DeleteDataflow(client client.Client, flow *v1alpha1.NifiDataflow, cluster *v1alpha1.NifiCluster) error {
	return nil
}




