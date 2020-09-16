package nificlient

import (
	"strconv"

	"github.com/antihax/optional"
	nigoapi "github.com/erdrix/nigoapi/pkg/nifi"
)

func (n *nifiClient) GetProcessGroup(id string) (*nigoapi.ProcessGroupEntity, error) {
	// Get nigoapi client, favoring the one associated to the coordinator node.
	client := n.privilegeCoordinatorClient()
	if client == nil {
		log.Error(ErrNoNodeClientsAvailable, "Error during creating node client")
		return nil, ErrNoNodeClientsAvailable
	}

	// Request on Nifi Rest API to get the process group informations
	pGEntity, rsp, err := client.ProcessGroupsApi.GetProcessGroup(nil, id)
	if err := errorGetOperation(rsp, err); err != nil {
		return nil, err
	}

	return &pGEntity, nil
}

func (n *nifiClient) CreateProcessGroup(
	entity nigoapi.ProcessGroupEntity,
	pgParentId string) (*nigoapi.ProcessGroupEntity, error) {
	// Get nigoapi client, favoring the one associated to the coordinator node.
	client := n.privilegeCoordinatorClient()
	if client == nil {
		log.Error(ErrNoNodeClientsAvailable, "Error during creating node client")
		return nil, ErrNoNodeClientsAvailable
	}

	// Request on Nifi Rest API to create the versioned process group
	pgEntity, rsp, err := client.ProcessGroupsApi.CreateProcessGroup(nil, pgParentId, entity)
	if err := errorCreateOperation(rsp, err); err != nil {
		return nil, err
	}

	return &pgEntity, nil
}

func (n *nifiClient) UpdateProcessGroup(entity nigoapi.ProcessGroupEntity) (*nigoapi.ProcessGroupEntity, error) {
	// Get nigoapi client, favoring the one associated to the coordinator node.
	client := n.privilegeCoordinatorClient()
	if client == nil {
		log.Error(ErrNoNodeClientsAvailable, "Error during creating node client")
		return nil, ErrNoNodeClientsAvailable
	}

	// Request on Nifi Rest API to update the versioned process group
	pgEntity, rsp, err := client.ProcessGroupsApi.UpdateProcessGroup(nil, entity.Id, entity)
	if err := errorUpdateOperation(rsp, err); err != nil {
		return nil, err
	}

	return &pgEntity, nil
}

func (n *nifiClient) RemoveProcessGroup(entity nigoapi.ProcessGroupEntity) error {
	// Get nigoapi client, favoring the one associated to the coordinator node.
	client := n.privilegeCoordinatorClient()
	if client == nil {
		log.Error(ErrNoNodeClientsAvailable, "Error during creating node client")
		return ErrNoNodeClientsAvailable
	}

	// Request on Nifi Rest API to remove the versioned process group
	_, rsp, err := client.ProcessGroupsApi.RemoveProcessGroup(
		nil,
		entity.Id,
		&nigoapi.ProcessGroupsApiRemoveProcessGroupOpts{
			Version: optional.NewString(strconv.FormatInt(*entity.Revision.Version, 10)),
		})

	return errorDeleteOperation(rsp, err)
}
