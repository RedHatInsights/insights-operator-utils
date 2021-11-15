// Copyright 2021 Red Hat, Inc
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

package metrics

import (
	"net/http"

	"github.com/rs/zerolog/log"
)

// PushGatewayClient is a simple wrapper over http.Client so that prometheus
// can do HTTP requests with the given authentication header
type PushGatewayClient struct {
	AuthToken string

	HTTPClient http.Client
}

// Do is a simple wrapper over http.Client.Do method that includes
// the authentication header configured in the PushGatewayClient instance
func (pgc *PushGatewayClient) Do(request *http.Request) (*http.Response, error) {
	if pgc.AuthToken != "" {
		log.Debug().Msg("Adding authorization header to HTTP request")
		request.Header.Set("Authorization", "Basic "+pgc.AuthToken)
	} else {
		log.Debug().Msg("No authorization token provided. Making HTTP request without credentials.")
	}
	log.Debug().
		Str("request", request.URL.String()).
		Str("method", request.Method).
		Msg("Pushing metrics to Prometheus push gateway")
	resp, err := pgc.HTTPClient.Do(request)
	if resp != nil {
		log.Debug().Int("code", resp.StatusCode).Msg("Returned status code")
	}
	return resp, err
}
