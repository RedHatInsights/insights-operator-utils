package httputils_test

import (
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
