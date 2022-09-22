// Copyright 2020, 2021, 2022 Red Hat, Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package httputils_test

// Documentation in literate-programming-style is available at:
// https://redhatinsights.github.io/insights-operator-utils/packages/http/openapi_test.html

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	httputils "github.com/RedHatInsights/insights-operator-utils/http"
	"github.com/RedHatInsights/insights-operator-utils/tests/helpers"
)

func TestFilterOutDebugMethods(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		for input, expectedOutput := range map[string]string{
			openAPIFile: openAPIFileWithoutDebugEndpoints,
		} {
			output, err := httputils.FilterOutDebugMethods(input)
			helpers.FailOnError(t, err)
			assert.JSONEq(t, expectedOutput, output)
		}
	})

	t.Run("Error", func(t *testing.T) {
		expectedErr := `error unmarshaling JSON: Failed to unmarshal extension properties: ` +
			`json: cannot unmarshal string into Go value of type map[string]json.RawMessage` + "\n" +
			`Input: "definitely-not-json"`
		_, err := httputils.FilterOutDebugMethods("definitely-not-json")
		assert.EqualError(t, err, expectedErr)
	})
}

const (
	openAPIFile = `
{
	"openapi": "3.0.0",
	"info": {
		"title": "title",
		"description": "description",
		"version": "1.0.0"
	},
	"components": {},
	"paths": {
		"/openapi.json": {
			"get": {
				"summary": "Returns the OpenAPI specification JSON.",
				"operationId": "getOpenApi",
				"responses": {
					"200": {
						"description": "A JSON containing the OpenAPI specification for this service."
					}
				}
			}
		},
		"/metrics": {
			"get": {
				"summary": "Read all metrics exposed by this service",
				"description": "Currently the following metrics are exposed to be consumed by Prometheus or any other tool compatible with it: 'consumed_messages' the total number of messages consumed from Kafka, 'consuming_errors' the total number of errors during consuming messages from Kafka, 'successful_messages_processing_time' the time to process successfully message, 'failed_messages_processing_time' the time to process message fail, 'last_checked_timestamp_lag_minutes' shows how slow we get messages from clusters, 'produced_messages' the total number of produced messages, 'written_reports' the total number of reports written to the storage, 'feedback_on_rules' the total number of left feedback, 'sql_queries_counter' the total number of SQL queries, 'sql_queries_durations' the SQL queries durations.  Additionally it is possible to consume all metrics provided by Go runtime. There metrics start with 'go_' and 'process_ 'prefixes.",
				"operationId": "getMetrics",
				"responses": {
					"200": {
						"description": "Default response containing all metrics in semi-structured text format"
					}
				}
			}
		},
		"/organizations": {
			"get": {
				"summary": "Returns a list of available organization IDs.",
				"operationId": "getOrganizations",
				"description": "[DEBUG ONLY] List of organizations for which at least one Insights report is available via the API.",
				"responses": {
					"200": {
						"description": "A JSON array of organization IDs."
					}
				},
				"tags": [
					"debug"
				]
			}
		},
		"/test": {
			"get": {
				"summary": "test",
				"operationId": "test",
				"description": "test",
				"responses": {
					"200": {
						"description": "test"
					}
				},
				"tags": [
					"debug"
				]
			},
			"post": {
				"summary": "test",
				"operationId": "test2",
				"description": "test",
				"responses": {
					"200": {
						"description": "test"
					}
				},
			}
		},
		"/organizations/{orgId}/clusters": {
			"get": {
				"summary": "Returns a list of clusters associated with the specified organization ID.",
				"operationId": "getClustersForOrganization",
				"parameters": [
					{
						"name": "orgId",
						"in": "path",
						"required": true,
						"description": "ID of the requested organization.",
						"schema": {
							"type": "integer",
							"format": "int64",
							"minimum": 0
						}
					}
				],
				"responses": {
					"200": {
						"description": "A JSON array of clusters that belong to the specified organization."
					}
				},
				"tags": [
					"prod"
				]
			}
		}
	}
}`

	openAPIFileWithoutDebugEndpoints = `
{
	"openapi": "3.0.0",
	"info": {
		"title": "title",
		"description": "description",
		"version": "1.0.0"
	},
	"components": {},
	"paths": {
		"/openapi.json": {
			"get": {
				"summary": "Returns the OpenAPI specification JSON.",
				"operationId": "getOpenApi",
				"responses": {
					"200": {
						"description": "A JSON containing the OpenAPI specification for this service."
					}
				}
			}
		},
		"/metrics": {
			"get": {
				"summary": "Read all metrics exposed by this service",
				"description": "Currently the following metrics are exposed to be consumed by Prometheus or any other tool compatible with it: 'consumed_messages' the total number of messages consumed from Kafka, 'consuming_errors' the total number of errors during consuming messages from Kafka, 'successful_messages_processing_time' the time to process successfully message, 'failed_messages_processing_time' the time to process message fail, 'last_checked_timestamp_lag_minutes' shows how slow we get messages from clusters, 'produced_messages' the total number of produced messages, 'written_reports' the total number of reports written to the storage, 'feedback_on_rules' the total number of left feedback, 'sql_queries_counter' the total number of SQL queries, 'sql_queries_durations' the SQL queries durations.  Additionally it is possible to consume all metrics provided by Go runtime. There metrics start with 'go_' and 'process_ 'prefixes.",
				"operationId": "getMetrics",
				"responses": {
					"200": {
						"description": "Default response containing all metrics in semi-structured text format"
					}
				}
			}
		},
		"/test": {
			"post": {
				"summary": "test",
				"operationId": "test2",
				"description": "test",
				"responses": {
					"200": {
						"description": "test"
					}
				}
			}
		},
		"/organizations/{orgId}/clusters": {
			"get": {
				"summary": "Returns a list of clusters associated with the specified organization ID.",
				"operationId": "getClustersForOrganization",
				"parameters": [
					{
						"name": "orgId",
						"in": "path",
						"required": true,
						"description": "ID of the requested organization.",
						"schema": {
							"type": "integer",
							"format": "int64",
							"minimum": 0
						}
					}
				],
				"responses": {
					"200": {
						"description": "A JSON array of clusters that belong to the specified organization."
					}
				},
				"tags": [
					"prod"
				]
			}
		}
	}
}`
)

// ResponseWriterMock is mock for http.ResponseWriter
type ResponseWriterMock struct {
	headerCalls      int
	writeCalls       int
	writeHeaderCalls int
	writeShouldFail  bool
}

// NewResponseWriterMock is constructor of ResponseWriterMock struct.
// Constructor takes care of all mock sub-structs, call counters etc.
func NewResponseWriterMock(writeShouldFail bool) ResponseWriterMock {
	return ResponseWriterMock{
		headerCalls:      0,
		writeCalls:       0,
		writeHeaderCalls: 0,
		writeShouldFail:  writeShouldFail,
	}
}

// Header is a method that needs to be implemented in order to satisfy
// ResponseWriter interface
func (w *ResponseWriterMock) Header() http.Header {
	w.headerCalls++
	return http.Header{}
}

// Write is a method that needs to be implemented in order to satisfy
// ResponseWriter interface
func (w *ResponseWriterMock) Write([]byte) (int, error) {
	w.writeCalls++
	if w.writeShouldFail {
		return -1, errors.New("Mocked error")
	}
	return 1, nil
}

// WriteHeader is a method that needs to be implemented in order to satisfy
// ResponseWriter interface
func (w *ResponseWriterMock) WriteHeader(statusCode int) {
	w.writeHeaderCalls++
}

// TestCreateAPIHandlerEmptyFilepath test the function CreateOpenAPIHandler
// when empty file name is provided
func TestCreateAPIHandlerEmptyFilepath(t *testing.T) {
	handler := httputils.CreateOpenAPIHandler("", true, false)

	writer := NewResponseWriterMock()
	handler(&writer, nil)

	// writer should be used to response with error
	assert.LessOrEqual(t, 1, writer.headerCalls)
	assert.LessOrEqual(t, 1, writer.writeCalls)
	assert.LessOrEqual(t, 1, writer.writeHeaderCalls)
}

// TestCreateAPIHandlerPathToExistingFile test the function CreateOpenAPIHandler
// when regular file name is provided
func TestCreateAPIHandlerPathToExistingFile(t *testing.T) {
	// that file should exists everywhere
	handler := httputils.CreateOpenAPIHandler("/etc/passwd", true, false)

	writer := NewResponseWriterMock()
	handler(&writer, nil)

	// writer should be used
	assert.Equal(t, 1, writer.headerCalls)
	assert.Equal(t, 1, writer.writeCalls)
	assert.Equal(t, 0, writer.writeHeaderCalls)
}
