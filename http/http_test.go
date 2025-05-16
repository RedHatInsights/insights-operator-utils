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
// https://redhatinsights.github.io/insights-operator-utils/packages/http/http_test.html

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/stretchr/testify/assert"

	httputils "github.com/RedHatInsights/insights-operator-utils/http"
)

var baseEndpoint = "organizations/{org_id}/clusters/{cluster}/users/{user_id}/report"

func TestMakeURLToEndpoint(t *testing.T) {
	assert.Equal(
		t,
		"api/prefix/organizations/2/clusters/cluster_id/users/1/report",
		httputils.MakeURLToEndpoint("api/prefix/", baseEndpoint, 2, "cluster_id", 1),
	)
}

func TestMakeURLToEndpointFromArray(t *testing.T) {
	assert.Equal(
		t,
		"api/prefix/organizations/2/clusters/cluster_id/users/1/report",
		httputils.MakeURLToEndpoint("api/prefix/", baseEndpoint, []interface{}{2, "cluster_id", 1}...),
	)
}

func TestMakeURLToEndpointWithSpaces(t *testing.T) {
	assert.Equal(
		t,
		"api/prefix/organizations/2/clusters/cluster%20id/users/1/report",
		httputils.MakeURLToEndpoint("api/prefix/", baseEndpoint, 2, "cluster id", 1),
	)
}

func TestMakeURLToEndpoint_NonParsedURLURL(t *testing.T) {
	assert.Equal(t, "\x00", httputils.MakeURLToEndpoint("\x00", ""))
}

func TestMakeURLToEndpointMap(t *testing.T) {
	assert.Equal(
		t,
		"api/prefix/organizations/2/clusters/cluster_id/users/1/report",
		httputils.MakeURLToEndpointMap("api/prefix/", baseEndpoint, map[string]interface{}{
			"org_id":  2,
			"cluster": "cluster_id",
			"user_id": 1,
		}),
	)
}

func TestMakeURLToEndpointMapString(t *testing.T) {
	assert.Equal(
		t,
		"api/prefix/organizations/2/clusters/cluster_id/users/1/report",
		httputils.MakeURLToEndpointMapString("api/prefix/", baseEndpoint, map[string]string{
			"org_id":  "2",
			"cluster": "cluster_id",
			"user_id": "1",
		}),
	)
}

func TestSetHTTPPrefix(t *testing.T) {
	assert.Equal(
		t,
		"http://someaddress:8000",
		httputils.SetHTTPPrefix("someaddress:8000"),
	)
}

func TestSetHTTPPrefixWithCompleteURL(t *testing.T) {
	assert.Equal(
		t,
		"http://someaddress:8000",
		httputils.SetHTTPPrefix("http://someaddress:8000"),
	)
}

// TestSendRequest test the behaviour of function SendRequest
func TestSendRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(`{"clusters":["first_cluster"]}`))
		if err != nil {
			log.Fatal().Msg(err.Error())
		}
	}))

	req, err := http.NewRequest(http.MethodGet, server.URL, http.NoBody)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	body, err := httputils.SendRequest(req, 15*time.Second)
	assert.NoError(t, err)
	assert.Equal(t, "{\"clusters\":[\"first_cluster\"]}", string(body))
}

// TestSendRequestErrorResponseHandling test the behaviour of function
// SendRequest when wrong URL or scheme is used
func TestSendRequestErrorResponseHandling(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "", http.NoBody)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	_, err = httputils.SendRequest(req, 15*time.Second)
	assert.Error(t, err)
}

// TestSendRequestEmptyBody test the behaviour of function SendRequest
func TestSendRequestEmptyBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	}))

	req, err := http.NewRequest(http.MethodGet, server.URL, http.NoBody)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	body, err := httputils.SendRequest(req, 15*time.Second)
	assert.NoError(t, err)
	assert.Empty(t, string(body))
}
