package v1alpha1

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NifiUserGroupSpec defines the desired state of NifiUserGroup
// +k8s:openapi-gen=true
type NifiUserGroupSpec struct {
	// Contains the reference to the NifiCluster with the one the registry client is linked.
	ClusterRef ClusterReference `json:"clusterRef"`
	// Contains the list of reference to NifiUsers that are part to the group.
	UsersRef []UserReference `json:"usersRef,omitempty"`
	AccessPolicies []string `json:"accessPolicies,omitempty"`
}

// NifiUserGroupStatus defines the observed state of NifiUserGroup
// +k8s:openapi-gen=true
type NifiUserGroupStatus struct {
	// The nifi registry client's id
	Id string `json:"id"`
	// The last nifi registry client revision version catched
	Version int64 `json:"version"`
}

// Nifi Registry Client is the Schema for the nifi registry client API
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:openapi-gen=true
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
type NifiUserGroup struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NifiUserGroupSpec   `json:"spec,omitempty"`
	Status NifiUserGroupStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// NifiUserGroupList contains a list of NifiUserGroup
type NifiUserGroupList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NifiUserGroup `json:"items"`
}

func init() {
	SchemeBuilder.Register(&NifiUserGroup{}, &NifiUserGroupList{})
}

func (n NifiUserGroup) GetIdentity() string {
	return fmt.Sprintf("%s-%s", n.Namespace, n.Name)
}
