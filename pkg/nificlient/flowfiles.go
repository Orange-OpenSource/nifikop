package nificlient

import (
	nigoapi "github.com/erdrix/nigoapi/pkg/nifi"
)

func (n *nifiClient) GetDropRequest(connectionId, id string)(*nigoapi.DropRequestEntity, error) {
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

func (n *nifiClient) CreateDropRequest(registryClient nigoapi.RegistryClientEntity)(*nigoapi.RegistryClientEntity, error) {
	// Get nigoapi client, favoring the one associated to the coordinator node.
	client := n.privilegeCoordinatorClient()
	if client == nil {
		log.Error(ErrNoNodeClientsAvailable, "Error during creating node client")
		return nil, ErrNoNodeClientsAvailable
	}

	// Request on Nifi Rest API to create the registry client
	regCliEntity, rsp, err := client.ControllerApi.CreateRegistryClient(nil, registryClient)
	if err := errorCreateOperation(rsp, err); err != nil {
		return nil, err
	}

	return &regCliEntity, nil
}