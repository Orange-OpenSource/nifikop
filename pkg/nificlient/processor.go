package nificlient

import nigoapi "github.com/erdrix/nigoapi/pkg/nifi"

func (n *nifiClient) UpdateProcessor(entity nigoapi.ProcessorEntity) (*nigoapi.ProcessorEntity, error) {
	// Get nigoapi client, favoring the one associated to the coordinator node.
	client := n.privilegeCoordinatorClient()
	if client == nil {
		log.Error(ErrNoNodeClientsAvailable, "Error during creating node client")
		return nil, ErrNoNodeClientsAvailable
	}

	// Request on Nifi Rest API to update the versioned processor
	processor, rsp, err := client.ProcessorsApi.UpdateProcessor(nil, entity.Id, entity)
	if err := errorUpdateOperation(rsp, err); err != nil {
		return nil, err
	}

	return &processor, nil
}

func (n *nifiClient) UpdateProcessorRunStatus(
	id string,
	entity nigoapi.ProcessorRunStatusEntity) (*nigoapi.ProcessorEntity, error) {

	// Get nigoapi client, favoring the one associated to the coordinator node.
	client := n.privilegeCoordinatorClient()
	if client == nil {
		log.Error(ErrNoNodeClientsAvailable, "Error during creating node client")
		return nil, ErrNoNodeClientsAvailable
	}

	// Request on Nifi Rest API to update the processor run status
	processor, rsp, err := client.ProcessorsApi.UpdateRunStatus(nil, id, entity)
	if err := errorUpdateOperation(rsp, err); err != nil {
		return nil, err
	}

	return &processor, nil
}
