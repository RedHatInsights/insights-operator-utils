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

package httputils_test

// Documentation in literate-programming-style is available at:
// https://redhatinsights.github.io/insights-operator-utils/packages/http/metrics_middleware_test.html

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"

	httputils "github.com/RedHatInsights/insights-operator-utils/http"
	"github.com/RedHatInsights/insights-operator-utils/metrics"
	"github.com/RedHatInsights/insights-operator-utils/responses"
	"github.com/RedHatInsights/insights-operator-utils/tests/helpers"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

const (
	localhostAddress = "localhost"
	port             = 8080
)

func TestLogRequest(t *testing.T) {
	buf := new(bytes.Buffer)
	log.Logger = zerolog.New(buf).With().Timestamp().Logger()

	server := createTestServer(t, []Endpoint{
		{
			Path: "/",
			Func: func(writer http.ResponseWriter, request *http.Request) {
				err := responses.Send(http.StatusOK, writer, responses.BuildOkResponse())
				helpers.FailOnError(t, err)
			},
			Methods: []string{http.MethodGet},
		},
	})

	resp, err := http.Get(fmt.Sprintf("http://%v:%v/", localhostAddress, port))
	helpers.FailOnError(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	err = server.Shutdown(context.TODO())
	helpers.FailOnError(t, err)

	assert.Contains(t, buf.String(), "Request received - URI: /, Method: GET")

	expected := `
	# HELP api_endpoints_requests The total number of requests per endpoint
	# TYPE api_endpoints_requests counter
	api_endpoints_requests{endpoint="/"} 1
	`
	assertCollectorEquals(t, metrics.APIRequests, expected)

	expected = `
	# HELP api_endpoints_status_codes API endpoints status codes
	# TYPE api_endpoints_status_codes counter
	api_endpoints_status_codes{endpoint="/", status_code="200"} 1
	`
	assertCollectorEquals(t, metrics.APIResponseStatusCodes, expected)
}

type Endpoint struct {
	Path    string
	Func    func(http.ResponseWriter, *http.Request)
	Methods []string
}

func createTestServer(t testing.TB, endpoints []Endpoint) *http.Server {
	router := mux.NewRouter().StrictSlash(true)
	router.Use(httputils.LogRequest)

	for _, endpoint := range endpoints {
		router.HandleFunc(endpoint.Path, endpoint.Func).Methods(endpoint.Methods...)
	}

	server := &http.Server{Addr: fmt.Sprintf("%v:%v", localhostAddress, port), Handler: router}

	listener, err := net.Listen("tcp", server.Addr)
	helpers.FailOnError(t, err)

	go func() {
		if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
			helpers.FailOnError(t, err)
		}
	}()

	return server
}

func assertCollectorEquals(t *testing.T, c prometheus.Collector, expected string) {
	err := testutil.CollectAndCompare(c, strings.NewReader(expected))
	helpers.FailOnError(t, err)
}
