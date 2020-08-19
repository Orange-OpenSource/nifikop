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
	"strconv"

	"emperror.dev/errors"
	"github.com/antihax/optional"
	nigoapi "github.com/erdrix/nigoapi/pkg/nifi"
)

func (n *nifiClient) GetRegistryClient(id string)(*nigoapi.RegistryClientEntity, error) {
	// Get nigoapi client, favoring the one associated to the coordinator node.
	client := n.privilegeCoordinatorClient()
	if client == nil {
		log.Error(ErrNoNodeClientsAvailable, "Error during creating node client")
		return nil, ErrNoNodeClientsAvailable
	}

	// Request on Nifi Rest API to get the registy client informations
	nodeEntity, rsp, err := client.ControllerApi.GetRegistryClient(nil, id)

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

	return &nodeEntity, nil
}

func (n *nifiClient) CreateRegistryClient(registryClient nigoapi.RegistryClientEntity)(*nigoapi.RegistryClientEntity, error) {
	// Get nigoapi client, favoring the one associated to the coordinator node.
	client := n.privilegeCoordinatorClient()
	if client == nil {
		log.Error(ErrNoNodeClientsAvailable, "Error during creating node client")
		return nil, ErrNoNodeClientsAvailable
	}

	// Request on Nifi Rest API to create the registry client
	regCliEntity, rsp, err := client.ControllerApi.CreateRegistryClient(nil, registryClient)

	if rsp != nil && rsp.StatusCode != 201 {
		log.Error(errors.New("Non 201 response from nifi node: "+rsp.Status), "Error during talking to nifi node")
		return nil, ErrNifiClusterNotReturned201
	}

	if err != nil || rsp == nil {
		log.Error(err, "Error during talking to nifi node")
		return nil, err
	}

	return &regCliEntity, nil
}

func (n *nifiClient) UpdateRegistryClient(registryClient nigoapi.RegistryClientEntity)(*nigoapi.RegistryClientEntity, error) {
	// Get nigoapi client, favoring the one associated to the coordinator node.
	client := n.privilegeCoordinatorClient()
	if client == nil {
		log.Error(ErrNoNodeClientsAvailable, "Error during creating node client")
		return nil, ErrNoNodeClientsAvailable
	}

	// Request on Nifi Rest API to update the registry client
	regCliEntity, rsp, err := client.ControllerApi.UpdateRegistryClient(nil, registryClient.Id, registryClient)

	if rsp != nil && rsp.StatusCode != 200 {
		log.Error(errors.New("Non 200 response from nifi node: "+rsp.Status), "Error during talking to nifi node")
		return nil, ErrNifiClusterNotReturned200
	}

	if err != nil || rsp == nil {
		log.Error(err, "Error during talking to nifi node")
		return nil, err
	}

	return &regCliEntity, nil
}

func (n *nifiClient) RemoveRegistryClient(registryClient nigoapi.RegistryClientEntity) error {
	// Get nigoapi client, favoring the one associated to the coordinator node.
	client := n.privilegeCoordinatorClient()
	if client == nil {
		log.Error(ErrNoNodeClientsAvailable, "Error during creating node client")
		return ErrNoNodeClientsAvailable
	}

	// Request on Nifi Rest API to remove the registry client
	_, rsp, err := client.ControllerApi.DeleteRegistryClient(nil, registryClient.Id,
		&nigoapi.ControllerApiDeleteRegistryClientOpts{
			Version: optional.NewString(strconv.FormatInt(*registryClient.Revision.Version, 10)),
	})

	if rsp != nil && rsp.StatusCode == 404 {
		log.Error(errors.New("404 response from nifi node: "+rsp.Status), "No registry client to remove found")
		return nil
	}

	if rsp != nil && rsp.StatusCode != 200 {
		log.Error(errors.New("Non 200 response from nifi node: "+rsp.Status), "Error during talking to nifi node")
		return ErrNifiClusterNotReturned200
	}

	if err != nil || rsp == nil {
		log.Error(err, "Error during talking to nifi node")
		return err
	}

	return nil
}