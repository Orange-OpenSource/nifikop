package dataflow

import (
	"github.com/Orange-OpenSource/nifikop/pkg/apis/nifi/v1alpha1"
	"github.com/Orange-OpenSource/nifikop/pkg/controller/common"
	"github.com/Orange-OpenSource/nifikop/pkg/errorfactory"
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

func SyncDataflow(client client.Client, flow *v1alpha1.NifiDataflow, cluster *v1alpha1.NifiCluster) (*v1alpha1.NifiDataflowStatus, error) {

	nClient, err := common.NewNodeConnection(log, client, cluster)
	if err != nil {
		return nil, err
	}

	if flow.Spec.UpdateStrategy == v1alpha1.DropStrategy {
		// unschedule processors
		_, err := nClient.UpdateFlowProcessGroup(nigoapi.ScheduleComponentsEntity{
			Id:    flow.Status.ProcessGroupID,
			State: "STOPPED",
		})

		if err == nificlient.ErrNifiClusterNotReturned200 {
			log.Error(err, "Stop flow failed since Nifi node returned non 200")
			return nil, err
		}

		if err != nil {
			log.Error(err, "could not communicate with nifi node")
			return nil, err
		}

		// Get flow
		flowEntity, err := nClient.GetFlow(flow.Status.ProcessGroupID)
		if err == nificlient.ErrNifiClusterNotReturned200 {
			log.Error(err, "Stop flow failed since Nifi node returned non 200")
			return nil, err
		}

		if err != nil {
			log.Error(err, "could not communicate with nifi node")
			return nil, err
		}

		//
		if !flow.Status.LatestDropRequest.Finished {
			dropRequest, err := nClient.GetDropRequest(flow.Status.LatestDropRequest.ConnectionId, flow.Status.LatestDropRequest.Id)
			if err != nificlient.ErrNifiClusterReturned404 {
				if err == nificlient.ErrNifiClusterNotReturned200 {
					log.Error(err, "Stop flow failed since Nifi node returned non 200")
					return nil, err
				}

				if err != nil {
					log.Error(err, "could not communicate with nifi node")
					return nil, err
				}

				if !dropRequest.DropRequest.Finished {
					flow.Status.LatestDropRequest =
						dropRequest2Status(flow.Status.LatestDropRequest.ConnectionId, dropRequest)
					return &flow.Status, errorfactory.NifiConnectionDropping{}
				}
			}

			flow.Status.LatestDropRequest = nil
		}

		// Drop all events in connections
		for _, connection := range flowEntity.ProcessGroupFlow.Flow.Connections {
			if connection.Status.AggregateSnapshot.FlowFilesQueued != 0 {

				break
			}
		}
	}

	return nil, nil
}

func DeleteDataflow(client client.Client, flow *v1alpha1.NifiDataflow, cluster *v1alpha1.NifiCluster) error {
	return nil
}

func dropRequest2Status(connectionId string, dropRequest *nigoapi.DropRequestEntity) *v1alpha1.DropRequest {
	dr := dropRequest.DropRequest
	return &v1alpha1.DropRequest{
		ConnectionId:     connectionId,
		Id:               dr.Id,
		Uri:              dr.Uri,
		LastUpdated:      dr.LastUpdated,
		Finished:         dr.Finished,
		FailureReason:    dr.FailureReason,
		PercentCompleted: dr.PercentCompleted,
		CurrentCount:     dr.CurrentCount,
		CurrentSize:      dr.CurrentSize,
		Current:          dr.Current,
		OriginalCount:    dr.OriginalCount,
		OriginalSize:     dr.OriginalSize,
		Original:         dr.Original,
		DroppedCount:     dr.DroppedCount,
		DroppedSize:      dr.DroppedSize,
		Dropped:          dr.Dropped,
		State:            dr.State,
	}
}




