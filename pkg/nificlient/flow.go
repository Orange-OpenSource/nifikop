package nificlient

import (
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
	if err := errorGetOperation(rsp, err); err != nil {
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
	if err := errorUpdateOperation(rsp, err); err != nil {
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
	if err := errorUpdateOperation(rsp, err); err != nil {
		return nil, err
	}

	return &csEntity, nil
}

func (n *nifiClient) FlowDropRequest(connectionId, id string)(*nigoapi.DropRequestEntity, error) {
	// Get nigoapi client, favoring the one associated to the coordinator node.
	client := n.privilegeCoordinatorClient()
	if client == nil {
		log.Error(ErrNoNodeClientsAvailable, "Error during creating node client")
		return nil, ErrNoNodeClientsAvailable
	}

	// Request on Nifi Rest API to get the drop request information
	dropRequest, rsp, err := client.FlowfileQueuesApi.GetDropRequest(nil, connectionId, id)
	if err := errorGetOperation(rsp, err); err != nil {
		return nil, err
	}

	return &dropRequest, nil
}