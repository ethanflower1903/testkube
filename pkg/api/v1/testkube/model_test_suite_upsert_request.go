/*
 * Testkube API
 *
 * Testkube provides a Kubernetes-native framework for test definition, execution and results
 *
 * API version: 1.0.0
 * Contact: testkube@kubeshop.io
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package testkube

// test create request body
type TestSuiteUpsertRequest struct {
	// object kubernetes namespace
	Namespace string `json:"namespace"`
	// object name
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	// Run this step before whole suite
	Before []TestSuiteStep `json:"before,omitempty"`
	// test suite labels
	Labels map[string]string `json:"labels,omitempty"`
	// Steps to run
	Steps []TestSuiteStep `json:"steps"`
	// Run this step after whole suite
	After   []TestSuiteStep `json:"after,omitempty"`
	Repeats int32           `json:"repeats,omitempty"`
}
