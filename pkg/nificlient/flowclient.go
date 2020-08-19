package nificlient

import (
	"emperror.dev/errors"
	nigoapi "github.com/erdrix/nigoapi/pkg/nifi"
)

func (n *nifiClient) GetFlow(id string)(*nigoapi.ProcessGroupFlowEntity, error) {
	// Get nigoapi client, favoring the one associated to the coordinator node.
	client := n.privilegeCoordinatorClient()
	if client == nil {
		log.Error(ErrNoNodeClientsAvailable, "Error during creating node client")
		return nil, ErrNoNodeClientsAvailable
	}

	// Request on Nifi Rest API to get the process group flow informations
	flowPGEntity, rsp, err := client.FlowApi.GetFlow(nil, id)

	if rsp != nil && rsp.StatusCode == 404 {
		return nil, ErrNifiClusterReturned404
	}

	if rsp != nil && rsp.StatusCode != 200 {
		log.Error(errors.New("Non 200 response from nifi node: "+rsp.Status), "Error during talking to nifi node")
		return nil, ErrNifiClusterNotReturned200
	}

	if err != nil || rsp == nil {
		log.Error(err, "Error during talking to nifi node")
		return nil, err
	}

	return &flowPGEntity, nil
}


func (n *nifiClient) UpdateFlowControllerServices(entity nigoapi.ActivateControllerServicesEntity)(*nigoapi.ActivateControllerServicesEntity, error) {

	// Get nigoapi client, favoring the one associated to the coordinator node.
	client := n.privilegeCoordinatorClient()
	if client == nil {
		log.Error(ErrNoNodeClientsAvailable, "Error during creating node client")
		return nil, ErrNoNodeClientsAvailable
	}

	// Request on Nifi Rest API to enable or disable the controller services
	csEntity, rsp, err := client.FlowApi.ActivateControllerServices(nil, entity.Id, entity)

	if rsp != nil && rsp.StatusCode != 200 {
		log.Error(errors.New("Non 200 response from nifi node: "+rsp.Status), "Error during talking to nifi node")
		return nil, ErrNifiClusterNotReturned200
	}

	if err != nil || rsp == nil {
		log.Error(err, "Error during talking to nifi node")
		return nil, err
	}

	return &csEntity, nil
}

func (n *nifiClient) UpdateFlowProcessGroup(entity nigoapi.ScheduleComponentsEntity)(*nigoapi.ScheduleComponentsEntity, error) {

	// Get nigoapi client, favoring the one associated to the coordinator node.
	client := n.privilegeCoordinatorClient()
	if client == nil {
		log.Error(ErrNoNodeClientsAvailable, "Error during creating node client")
		return nil, ErrNoNodeClientsAvailable
	}

	// Request on Nifi Rest API to enable or disable the controller services
	csEntity, rsp, err := client.FlowApi.ScheduleComponents(nil, entity.Id, entity)

	if rsp != nil && rsp.StatusCode != 200 {
		log.Error(errors.New("Non 200 response from nifi node: "+rsp.Status), "Error during talking to nifi node")
		return nil, ErrNifiClusterNotReturned200
	}

	if err != nil || rsp == nil {
		log.Error(err, "Error during talking to nifi node")
		return nil, err
	}

	return &csEntity, nil
}