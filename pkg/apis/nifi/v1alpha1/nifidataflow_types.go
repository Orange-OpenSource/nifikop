package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NifiDataflowSpec defines the desired state of NifiDataflow
// +k8s:openapi-gen=true
type NifiDataflowSpec struct {
	// The id of the parent process group where you want to deploy your dataflow, if not set deploy at root level
	ParentProcessGroupID string `json:"parentProcessGroupID,omitempty"`
	// The UUID of the Bucket containing the flow.
	BucketId string `json:"bucketId"`
	// The UUID of the flow to run.
	FlowId string `json:"flowId"`
	// The version of the flow to run, if not present or equals to -1, then the latest version of flow will be used.
	FlowVersion *int32 `json:"flowVersion,omitempty"`
	// Object that will be passed to the NiFi Flow as parameteres.
	ParameterContextRef *ParameterContextReference `json:"parameterContextRef,omitempty"`
	// If the flow will be ran once or continuously checked
	RunOnce *bool `json:"runOnce,omitempty"`
	//
	SkipInvalidControllerService bool `json:"skipInvalidControllerService,omitempty"`
	//
	SkipInvalidComponent bool `json:"skipInvalidComponent,omitempty"`
	//
	ClusterRef ClusterReference `json:"clusterRef,omitempty"`
	//
	RegistryClientRef *RegistryClientReference `json:"registryClientRef,omitempty"`
	//
	// +kubebuilder:validation:Enum={"drop","drain"}
	UpdateStrategy DataflowUpdateStrategy `json:"updateStrategy"`
}

type UpdateRequest struct {
	// Defines the type of versioned flow update request.
	Type DataflowUpdateRequestType `json:"type"`
	// The id of the update request.
	Id string `json:"id"`
	// The uri for this request.
	Uri string `json:"uri"`
	// The last time this request was updated.
	LastUpdated string `json:"lastUpdated"`
	// Whether or not this request has completed.
	Complete bool `json:"complete"`
	// An explication of why the request failed, or null if this request has not failed.
	FailureReason string `json:"failureReason"`
	// The percentage complete of the request, between 0 and 100.
	PercentCompleted int32 `json:"percentCompleted"`
	// The state of the request
	State string `json:"state"`
}

type DropRequest struct {
	ConnectionId string `json:"connectionId"`
	// The id for this drop request.
	Id string `json:"id"`
	// The uri for this request.
	Uri string `json:"uri"`
	// The last time this request was updated.
	LastUpdated string `json:"lastUpdated"`
	// Whether the request has finished.
	Finished bool `json:"finished"`
	// An explication of why the request failed, or null if this request has not failed.
	FailureReason string `json:"failureReason"`
	// The percentage complete of the request, between 0 and 100.
	PercentCompleted int32 `json:"percentCompleted"`
	// The number of flow files currently queued.
	CurrentCount int32 `json:"currentCount"`
	// The size of flow files currently queued in bytes.
	CurrentSize int64 `json:"currentSize"`
	// The count and size of flow files currently queued.
	Current string `json:"current"`
	// The number of flow files to be dropped as a result of this request.
	OriginalCount int32 `json:"originalCount"`
	// The size of flow files to be dropped as a result of this request in bytes.
	OriginalSize int64 `json:"originalSize"`
	// The count and size of flow files to be dropped as a result of this request.
	Original string `json:"original"`
	// The number of flow files that have been dropped thus far.
	DroppedCount int32 `json:"droppedCount"`
	// The size of flow files currently queued in bytes.
	DroppedSize int64 `json:"droppedSize"`
	// The count and size of flow files that have been dropped thus far.
	Dropped string `json:"dropped"`
	// The state of the request
	State string `json:"state"`
}

// NifiDataflowStatus defines the observed state of NifiDataflow
// +k8s:openapi-gen=true
type NifiDataflowStatus struct {
	// Queued flow files
	// Process Group ID
	ProcessGroupID string `json:"processGroupID"`
	//
	State DataflowState `json:"state"`
	//
	LatestUpdateRequest *UpdateRequest `json:"latestUpdateRequest,omitempty"`
	//
	LatestDropRequest *DropRequest `json:"latestDropRequest,omitempty"`
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
