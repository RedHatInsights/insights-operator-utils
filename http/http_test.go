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

import (
	"testing"

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
