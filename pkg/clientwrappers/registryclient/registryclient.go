package registryclient

import (
	"github.com/Orange-OpenSource/nifikop/pkg/apis/nifi/v1alpha1"
	"github.com/Orange-OpenSource/nifikop/pkg/controller/common"
	"github.com/Orange-OpenSource/nifikop/pkg/nificlient"
	nigoapi "github.com/erdrix/nigoapi/pkg/nifi"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

var log = logf.Log.WithName("registryclient-method")

func CreateRegistryClient(client client.Client, registryClient *v1alpha1.NifiRegistryClient,
		cluster *v1alpha1.NifiCluster) (*v1alpha1.NifiRegistryClientStatus, error) {
	nClient, err := common.NewNodeConnection(log, client, cluster)
	if err != nil {
		return nil, err
	}

	scratchEntity := nigoapi.RegistryClientEntity{}
	updateRegistryClientEntity(registryClient, &scratchEntity)
	entity, err := nClient.CreateRegistryClient(scratchEntity)

	if err != nil && err != nificlient.ErrNifiClusterNotReturned201 {
		log.Error(err, "could not communicate with nifi node")
		return nil, err
	}

	if err == nificlient.ErrNifiClusterNotReturned201 {
		log.Error(err, "Create client registry failed since Nifi node returned non 201")
		return nil, err
	}


	return &v1alpha1.NifiRegistryClientStatus{
		Id:      entity.Id,
		Version: *entity.Revision.Version,
	}, nil
}

func SyncRegistryClient(client client.Client, registryClient *v1alpha1.NifiRegistryClient,
		cluster *v1alpha1.NifiCluster) (*v1alpha1.NifiRegistryClientStatus, error) {

	nClient, err := common.NewNodeConnection(log, client, cluster)
	if err != nil {
		return nil, err
	}

	entity, err := nClient.GetRegistryClient(registryClient.Status.Id)

	if err == nificlient.ErrNifiClusterReturned404 {
		return CreateRegistryClient(client, registryClient, cluster)
	}

	if err == nificlient.ErrNifiClusterNotReturned200 {
		log.Error(err, "Find client registry failed since Nifi node returned non 200")
		return nil, err
	}

	if err != nil {
		log.Error(err, "could not communicate with nifi node")
		return nil, err
	}

	if !registryClientIsSync(registryClient, entity) {
		updateRegistryClientEntity(registryClient, entity)
		entity, err = nClient.UpdateRegistryClient(*entity)
		if err != nil && err != nificlient.ErrNifiClusterNotReturned200 {
			log.Error(err, "could not communicate with nifi node")
			return nil, err
		}

		if err == nificlient.ErrNifiClusterNotReturned200 {
			log.Error(err, "Update client registry failed since Nifi node returned non 200")
			return nil, err
		}
	}

	status := registryClient.Status
	status.Version = *entity.Revision.Version
	status.Id = entity.Id

	return &status, nil
}

func RemoveRegistryClient(client client.Client, registryClient *v1alpha1.NifiRegistryClient,
	cluster *v1alpha1.NifiCluster) error {
	nClient, err := common.NewNodeConnection(log, client, cluster)
	if err != nil {
		return err
	}

	entity, err := nClient.GetRegistryClient(registryClient.Status.Id)
	if err != nil && err != nificlient.ErrNifiClusterNotReturned200 {
		log.Error(err, "could not communicate with nifi node")
		return err
	}

	if err == nificlient.ErrNifiClusterNotReturned200 {
		log.Error(err, "Find client registry failed since Nifi node returned non 200")
		return err
	}

	updateRegistryClientEntity(registryClient, entity)
	err = nClient.RemoveRegistryClient(*entity)
	if err != nil && err != nificlient.ErrNifiClusterNotReturned200 {
		log.Error(err, "could not communicate with nifi node")
		return err
	}

	if err == nificlient.ErrNifiClusterNotReturned200 {
		log.Error(err, "Create client registry failed since Nifi node returned non 200")
		return  err
	}

	return nil
}

func registryClientIsSync(registryClient *v1alpha1.NifiRegistryClient, entity *nigoapi.RegistryClientEntity) bool {
	return registryClient.Name == entity.Component.Name &&
		registryClient.Spec.Description == entity.Component.Description &&
		registryClient.Spec.Uri == entity.Component.Uri
}

func updateRegistryClientEntity(registryClient *v1alpha1.NifiRegistryClient, entity *nigoapi.RegistryClientEntity){

	var defaultVersion int64 = 0

	if entity == nil {
		entity = &nigoapi.RegistryClientEntity{}
	}

	if entity.Component == nil {
		entity.Revision = &nigoapi.RevisionDto{
			Version: &defaultVersion,
		}
	}

	if entity.Component == nil {
		entity.Component = &nigoapi.RegistryDto{
		}
	}

	entity.Component.Name        = registryClient.Name
	entity.Component.Description = registryClient.Spec.Description
	entity.Component.Uri         = registryClient.Spec.Uri
}
