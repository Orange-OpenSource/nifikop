package usergroup

import (
	"github.com/Orange-OpenSource/nifikop/pkg/apis/nifi/v1alpha1"
	"github.com/Orange-OpenSource/nifikop/pkg/clientwrappers"
	"github.com/Orange-OpenSource/nifikop/pkg/clientwrappers/accesspolicies"
	"github.com/Orange-OpenSource/nifikop/pkg/controller/common"
	"github.com/Orange-OpenSource/nifikop/pkg/nificlient"
	nigoapi "github.com/erdrix/nigoapi/pkg/nifi"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

var log = logf.Log.WithName("usergroup-method")

func ExistUserGroup(client client.Client, userGroup *v1alpha1.NifiUserGroup,
	cluster *v1alpha1.NifiCluster) (bool, error) {

	if userGroup.Status.Id == "" {
		return false, nil
	}

	nClient, err := common.NewNodeConnection(log, client, cluster)
	if err != nil {
		return false, err
	}

	entity, err := nClient.GetUserGroup(userGroup.Status.Id)
	if err := clientwrappers.ErrorGetOperation(log, err, "Get user-group"); err != nil {
		if err == nificlient.ErrNifiClusterReturned404 {
			return false, nil
		}
		return false, err
	}

	return entity != nil, nil
}

func CreateUserGroup(client client.Client, userGroup *v1alpha1.NifiUserGroup, users []*v1alpha1.NifiUser,
	cluster *v1alpha1.NifiCluster) (*v1alpha1.NifiUserGroupStatus, error) {
	nClient, err := common.NewNodeConnection(log, client, cluster)
	if err != nil {
		return nil, err
	}

	scratchEntity := nigoapi.UserGroupEntity{}
	updateUserGroupEntity(userGroup, users, &scratchEntity)

	entity, err := nClient.CreateUserGroup(scratchEntity)
	if err := clientwrappers.ErrorCreateOperation(log, err, "Create user-group"); err != nil {
		return nil, err
	}

	return &v1alpha1.NifiUserGroupStatus{
		Id:      entity.Id,
		Version: *entity.Revision.Version,
	}, nil
}

func SyncUserGroup(client client.Client, userGroup *v1alpha1.NifiUserGroup, users []*v1alpha1.NifiUser,
	cluster *v1alpha1.NifiCluster) (*v1alpha1.NifiUserGroupStatus, error) {

	nClient, err := common.NewNodeConnection(log, client, cluster)
	if err != nil {
		return nil, err
	}

	entity, err := nClient.GetUserGroup(userGroup.Status.Id)
	if err := clientwrappers.ErrorGetOperation(log, err, "Get user-group"); err != nil {
		return nil, err
	}

	if !userGroupIsSync(userGroup, users, entity) {
		updateUserGroupEntity(userGroup, users, entity)
		entity, err = nClient.UpdateUserGroup(*entity)
		if err := clientwrappers.ErrorUpdateOperation(log, err, "Update user-group"); err != nil {
			return nil, err
		}
	}

	status := userGroup.Status
	status.Version = *entity.Revision.Version
	status.Id = entity.Id

	// Remove from access policy
	for _, entity := range entity.Component.AccessPolicies {
		contains := false
		for _,  accessPolicy := range userGroup.Spec.AccessPolicies {
			if entity.Component.Action == string(accessPolicy.Action) &&
				entity.Component.Resource == accessPolicy.GetResource(cluster) {
				contains = true
				break
			}
		}
		if !contains {
			if err := accesspolicies.UpdateAccessPolicyEntity(client, &entity,
				[]*v1alpha1.NifiUser{}, []*v1alpha1.NifiUser{},
				[]*v1alpha1.NifiUserGroup{}, []*v1alpha1.NifiUserGroup{userGroup}, cluster); err != nil {
				return &status, err
			}
		}
	}

	// add
	for _,  accessPolicy := range userGroup.Spec.AccessPolicies {
		contains := false
		for _, entity := range entity.Component.AccessPolicies {
			if entity.Component.Action == string(accessPolicy.Action) &&
				entity.Component.Resource == accessPolicy.GetResource(cluster) {
				contains = true
				break
			}
		}
		if !contains {
			if err := accesspolicies.UpdateAccessPolicy(client, &accessPolicy,
				[]*v1alpha1.NifiUser{}, []*v1alpha1.NifiUser{},
				[]*v1alpha1.NifiUserGroup{userGroup}, []*v1alpha1.NifiUserGroup{}, cluster); err != nil {
				return &status, err
			}
		}
	}

	return &status, nil
}

func RemoveUserGroup(client client.Client, userGroup *v1alpha1.NifiUserGroup, users []*v1alpha1.NifiUser,
	cluster *v1alpha1.NifiCluster) error {
	nClient, err := common.NewNodeConnection(log, client, cluster)
	if err != nil {
		return err
	}

	entity, err := nClient.GetUserGroup(userGroup.Status.Id)
	if err := clientwrappers.ErrorGetOperation(log, err, "Get user-group"); err != nil {
		if err == nificlient.ErrNifiClusterReturned404 {
			return nil
		}
		return err
	}

	updateUserGroupEntity(userGroup, users, entity)
	err = nClient.RemoveUserGroup(*entity)

	return clientwrappers.ErrorRemoveOperation(log, err, "Remove user-group")
}

func userGroupIsSync(
	userGroup *v1alpha1.NifiUserGroup,
	users []*v1alpha1.NifiUser,
	entity *nigoapi.UserGroupEntity) bool {

	if userGroup.GetIdentity() != entity.Component.Identity {
		return false
	}

	for _, expected := range users {
		notFound := true
		for _, tenant := range entity.Component.Users {
			if expected.Status.Id == tenant.Id {
				notFound = false
				break
			}
		}
		if notFound {
			return false
		}
	}
	return true
}

func updateUserGroupEntity(userGroup *v1alpha1.NifiUserGroup, users []*v1alpha1.NifiUser, entity *nigoapi.UserGroupEntity) {

	var defaultVersion int64 = 0

	if entity == nil {
		entity = &nigoapi.UserGroupEntity{}
	}

	if entity.Component == nil {
		entity.Revision = &nigoapi.RevisionDto{
			Version: &defaultVersion,
		}
	}

	if entity.Component == nil {
		entity.Component = &nigoapi.UserGroupDto{
		}
	}

	entity.Component.Identity = userGroup.GetIdentity()

	for _, user := range users {
		entity.Component.Users = append(entity.Component.Users, nigoapi.TenantEntity{Id: user.Status.Id})
	}
}
