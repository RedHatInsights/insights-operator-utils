// Copyright 2020 Red Hat, Inc
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

package metrics_test

import (
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/RedHatInsights/insights-operator-utils/tests/helpers"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	prommodels "github.com/prometheus/client_model/go"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"

	"github.com/RedHatInsights/insights-operator-utils/metrics"
)

const (
	testCaseTimeLimit = 5 * time.Second
	apiPrefix         = "/api/"
	microAddress      = ":8080"
	testEndpoint      = "test"
)

func init() {
	zerolog.SetGlobalLevel(zerolog.WarnLevel)
}

func getCounterValue(counter prometheus.Counter) float64 {
	pb := &prommodels.Metric{}
	err := counter.Write(pb)
	if err != nil {
		panic(fmt.Sprintf("Unable to get counter from counter %v", err))
	}

	return pb.GetCounter().GetValue()
}

func getCounterVecValue(counterVec *prometheus.CounterVec, labels map[string]string) float64 {
	counter, err := counterVec.GetMetricWith(labels)
	if err != nil {
		panic(fmt.Sprintf("Unable to get counter from counterVec %v", err))
	}

	return getCounterValue(counter)
}

func prepareServer(status int) *helpers.MicroHTTPServer {
	server := helpers.NewMicroHTTPServer(microAddress, apiPrefix)
	server.Router.Use(metrics.LogRequest)
	server.AddEndpoint(testEndpoint, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
	})

	return server
}

func TestAPIRequestsMetric(t *testing.T) {
	helpers.RunTestWithTimeout(t, func(t *testing.T) {
		// resetting since go runs tests in 1 process
		metrics.APIRequests.Reset()
		finalEndpoint := apiPrefix + testEndpoint

		assert.Equal(t, 0.0, getCounterVecValue(metrics.APIRequests, map[string]string{
			"endpoint": finalEndpoint,
		}))

		server := prepareServer(http.StatusOK)

		helpers.AssertAPIRequest(t, server, apiPrefix, &helpers.APIRequest{
			Method:   http.MethodGet,
			Endpoint: testEndpoint,
		}, &helpers.APIResponse{
			StatusCode: http.StatusOK,
		})

		assert.Equal(t, 1.0, getCounterVecValue(metrics.APIRequests, map[string]string{
			"endpoint": finalEndpoint,
		}))
	}, testCaseTimeLimit)
}

func TestAPIResponsesTimeMetric(t *testing.T) {
	metrics.APIResponsesTime.Reset()

	err := testutil.CollectAndCompare(metrics.APIResponsesTime, strings.NewReader(""))
	helpers.FailOnError(t, err)

	metrics.APIResponsesTime.With(prometheus.Labels{"endpoint": "test"}).Observe(5.6)

	expected := `
		# HELP api_endpoints_response_time API endpoints response time
		# TYPE api_endpoints_response_time histogram
		api_endpoints_response_time_bucket{endpoint="test",le="0"} 0
		api_endpoints_response_time_bucket{endpoint="test",le="20"} 1
		api_endpoints_response_time_bucket{endpoint="test",le="40"} 1
		api_endpoints_response_time_bucket{endpoint="test",le="60"} 1
		api_endpoints_response_time_bucket{endpoint="test",le="80"} 1
		api_endpoints_response_time_bucket{endpoint="test",le="100"} 1
		api_endpoints_response_time_bucket{endpoint="test",le="120"} 1
		api_endpoints_response_time_bucket{endpoint="test",le="140"} 1
		api_endpoints_response_time_bucket{endpoint="test",le="160"} 1
		api_endpoints_response_time_bucket{endpoint="test",le="180"} 1
		api_endpoints_response_time_bucket{endpoint="test",le="200"} 1
		api_endpoints_response_time_bucket{endpoint="test",le="220"} 1
		api_endpoints_response_time_bucket{endpoint="test",le="240"} 1
		api_endpoints_response_time_bucket{endpoint="test",le="260"} 1
		api_endpoints_response_time_bucket{endpoint="test",le="280"} 1
		api_endpoints_response_time_bucket{endpoint="test",le="300"} 1
		api_endpoints_response_time_bucket{endpoint="test",le="320"} 1
		api_endpoints_response_time_bucket{endpoint="test",le="340"} 1
		api_endpoints_response_time_bucket{endpoint="test",le="360"} 1
		api_endpoints_response_time_bucket{endpoint="test",le="380"} 1
		api_endpoints_response_time_bucket{endpoint="test",le="+Inf"} 1
		api_endpoints_response_time_sum{endpoint="test"} 5.6
		api_endpoints_response_time_count{endpoint="test"} 1
	`
	err = testutil.CollectAndCompare(metrics.APIResponsesTime, strings.NewReader(expected))
	helpers.FailOnError(t, err)
}

func TestApiResponseStatusCodesMetric_StatusOK(t *testing.T) {
	helpers.RunTestWithTimeout(t, func(t *testing.T) {
		metrics.APIResponseStatusCodes.Reset()

		assert.Equal(t, 0.0, getCounterVecValue(metrics.APIResponseStatusCodes, map[string]string{
			"status_code": fmt.Sprint(http.StatusOK),
		}))

		server := prepareServer(http.StatusOK)

		for i := 0; i < 15; i++ {
			helpers.AssertAPIRequest(t, server, apiPrefix, &helpers.APIRequest{
				Method:   http.MethodGet,
				Endpoint: testEndpoint,
			}, &helpers.APIResponse{
				StatusCode: http.StatusOK,
			})
		}

		assert.Equal(t, 15.0, getCounterVecValue(metrics.APIResponseStatusCodes, map[string]string{
			"status_code": fmt.Sprint(http.StatusOK),
		}))
	}, testCaseTimeLimit)
}

func TestApiResponseStatusCodesMetric_StatusBadRequest(t *testing.T) {
	helpers.RunTestWithTimeout(t, func(t *testing.T) {
		metrics.APIResponseStatusCodes.Reset()

		assert.Equal(t, 0.0, getCounterVecValue(metrics.APIResponseStatusCodes, map[string]string{
			"status_code": fmt.Sprint(http.StatusBadRequest),
		}))

		server := prepareServer(http.StatusBadRequest)

		helpers.AssertAPIRequest(t, server, apiPrefix, &helpers.APIRequest{
			Method:   http.MethodGet,
			Endpoint: testEndpoint,
		}, &helpers.APIResponse{
			StatusCode: http.StatusBadRequest,
		})

		assert.Equal(t, 1.0, getCounterVecValue(metrics.APIResponseStatusCodes, map[string]string{
			"status_code": fmt.Sprint(http.StatusBadRequest),
		}))
	}, testCaseTimeLimit)
}
