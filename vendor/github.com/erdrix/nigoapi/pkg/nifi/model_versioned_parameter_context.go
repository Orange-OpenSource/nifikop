/*
 * NiFi Rest Api
 *
 * The Rest Api provides programmatic access to command and control a NiFi instance in real time. Start and                                              stop processors, monitor queues, query provenance data, and more. Each endpoint below includes a description,                                             definitions of the expected input and output, potential response codes, and the authorizations required                                             to invoke each service.
 *
 * API version: 1.12.0
 * Contact: dev@nifi.apache.org
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package nifi

type VersionedParameterContext struct {
	// The name of the context
	Name string `json:"name,omitempty"`
	// The description of the parameter context
	Description string `json:"description,omitempty"`
	// The parameters in the context
	Parameters []VersionedParameter `json:"parameters,omitempty"`
}
