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
