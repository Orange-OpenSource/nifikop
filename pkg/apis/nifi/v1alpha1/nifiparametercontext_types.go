package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NifiParameterContextSpec defines the desired state of NifiParameterContext
// +k8s:openapi-gen=true
type NifiParameterContextSpec struct {
	// The Description of the Parameter Context.
	Description string `json:"description,omitempty"`
	// The Parameters for the Parameter Context
	Parameters []Parameter `json:"parameters"`
	// contains the reference to the NifiCluster with the one the user is linked
	ClusterRef ClusterReference `json:"clusterRef"`
	// A list of secret containing sensitive parameters (the key will name of the parameter)
	SecretRefs []SecretReference `json:"secretRefs,omitempty"`
}

type Parameter struct {
	// The name of the Parameter.
	Name string `json:"name"`
	// The value of the Parameter.
	Value string `json:"value,omitempty"`
	// The description of the Parameter.
	Description string `json:"description,omitempty"`
}

// NifiParameterContextStatus defines the observed state of NifiParameterContext
// +k8s:openapi-gen=true
type NifiParameterContextStatus struct {
	// Queued flow files
	// The nifi parameter context id
	Id string `json:"id"`
	// The last nifi parameter context revision version catched
	Version int64 `json:"version"`
	// The latest update request
	LatestUpdateRequest *ParameterContextUpdateRequest `json:"latestUpdateRequest,omitempty"`
}

type ParameterContextUpdateRequest struct {
	// The id of the update request.
	Id string `json:"id"`
	// The uri for this request.
	Uri string `json:"uri"`
	// The timestamp of when the request was submitted This property is read only.
	SubmissionTime string `json:"submissionTime"`
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

// NifiParameterContext is the Schema for the nifi parameter context API
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:openapi-gen=true
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
type NifiParameterContext struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NifiParameterContextSpec   `json:"spec,omitempty"`
	Status NifiParameterContextStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// NifiParameterContextList contains a list of NifiParameterContext
type NifiParameterContextList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NifiParameterContext `json:"items"`
}

func init() {
	SchemeBuilder.Register(&NifiParameterContext{}, &NifiParameterContextList{})
}
