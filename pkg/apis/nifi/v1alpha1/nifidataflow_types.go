package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NifiDataflowSpec defines the desired state of NifiDataflow
// +k8s:openapi-gen=true
type NifiDataflowSpec struct {
	// The id of the parent process group where you want to deploy your dataflow, if not set deploy at root level
	ParentProcessGroupID string                  `json:"parentProcessGroupID,omitempty"`
	// The UUID of the Bucket containing the flow.
	BucketId             string                  `json:"bucketId"`
	// The UUID of the flow to run.
	FlowId               string                  `json:"flowId"`
	// The version of the flow to run, if not present or equals to -1, then the latest version of flow will be used.
	FlowVersion          *int                    `json:"flowVersion,omitempty"`
	// Object that will be passed to the NiFi Flow as parameteres.
	Parameters           []Parameter             `json:"parameters,omitempty"`
	// If the flow will be ran once or continuously checked
	RunOnce              *bool                   `json:"runOnce,omitempty"`
	//
	ClusterRef           ClusterReference        `json:"clusterRef,omitempty"`
	//
	RegistryClientRef    *RegistryClientReference `json:"registryClientRef,omitempty"`
	//
	UpdateStrategy
}

type Parameter struct {
	//
	Name        string `json:"name"`
	//
	Value       string `json:"value,omitempty"`
	//
	Sensitive   bool   `json:"sensitive,omitempty"`
	//
	Description string `json:"description,omitempty"`
	//
	SecretName  string `json:"secretName,omitempty"`
}

// NifiDataflowStatus defines the observed state of NifiDataflow
// +k8s:openapi-gen=true
type NifiDataflowStatus struct {
	// Queued flow files
	// Process Group ID
	ProcessGroupID string `json:"processGroupID"`
	//
	State DataflowState `json:"state"`
}

// Nifi Dataflow is the Schema for the nifi dataflow API
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:openapi-gen=true
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
type NifiDataflow struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NifiDataflowSpec   `json:"spec,omitempty"`
	Status NifiDataflowStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// NifiDataflowList contains a list of NifiDataflow
type NifiDataflowList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NifiDataflow `json:"items"`
}

func init() {
	SchemeBuilder.Register(&NifiDataflow{}, &NifiDataflowList{})
}