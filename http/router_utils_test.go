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
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/RedHatInsights/insights-results-aggregator-data/testdata"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"

	httputils "github.com/RedHatInsights/insights-operator-utils/http"
	"github.com/RedHatInsights/insights-operator-utils/tests/helpers"
	"github.com/RedHatInsights/insights-operator-utils/types"
)

const (
	cluster1ID = "715e10eb-e6ac-49b3-bd72-61734c35b6fb"
	cluster2ID = "931f1495-7b16-4637-a41e-963e117bfd02"
)

func TestGetRouterPositiveIntParam_NonIntError(t *testing.T) {
	request := mustGetRequestWithMuxVars(t, http.MethodGet, "", nil, map[string]string{
		"id": "non int",
	})

	_, err := httputils.GetRouterPositiveIntParam(request, "id")
	assert.EqualError(
		t,
		err,
		"Error during parsing param 'id' with value 'non int'. Error: 'unsigned integer expected'",
	)
}

func TestGetRouterPositiveIntParam(t *testing.T) {
	request := mustGetRequestWithMuxVars(t, http.MethodGet, "", nil, map[string]string{
		"id": "99",
	})

	id, err := httputils.GetRouterPositiveIntParam(request, "id")
	helpers.FailOnError(t, err)

	assert.Equal(t, uint64(99), id)
}

func TestGetRouterPositiveIntParam_ZeroError(t *testing.T) {
	request := mustGetRequestWithMuxVars(t, http.MethodGet, "", nil, map[string]string{
		"id": "0",
	})

	_, err := httputils.GetRouterPositiveIntParam(request, "id")
	assert.EqualError(t, err, "Error during parsing param 'id' with value '0'. Error: 'positive value expected'")
}

func TestGetRouterPositiveIntParam_Missing(t *testing.T) {
	request, err := http.NewRequest(http.MethodGet, "organizations//clusters", nil)
	helpers.FailOnError(t, err)

	_, err = httputils.GetRouterPositiveIntParam(request, "test")
	assert.EqualError(t, err, "Missing required param from request: test")
}

func TestReadParam(t *testing.T) {
	for _, testCase := range []struct {
		TestName   string
		ParamName  string
		ParamValue []interface{}
	}{
		{TestName: "cluster", ParamName: "cluster", ParamValue: []interface{}{testdata.ClusterName}},
		{TestName: "rule_id", ParamName: "rule_id", ParamValue: []interface{}{testdata.Rule1ID}},
		{TestName: "error_key", ParamName: "error_key", ParamValue: []interface{}{testdata.ErrorKey1}},
		{TestName: "organization", ParamName: "organization", ParamValue: []interface{}{testdata.OrgID}},
		{
			TestName:   "organization/with_auth",
			ParamName:  "organization",
			ParamValue: []interface{}{testdata.OrgID},
		},
		{
			TestName:   "clusters",
			ParamName:  "clusters",
			ParamValue: []interface{}{testdata.ClusterName, testdata.ClusterName, testdata.ClusterName},
		},
		{
			TestName:   "organizations",
			ParamName:  "organizations",
			ParamValue: []interface{}{testdata.OrgID, testdata.OrgID},
		},
	} {
		expectedParamValue := paramsToString(",", testCase.ParamValue...)

		t.Run(testCase.TestName, func(t *testing.T) {
			request := mustGetRequestWithMuxVars(t, http.MethodGet, "", nil, map[string]string{
				testCase.ParamName: expectedParamValue,
			})
			recorder := httptest.NewRecorder()

			var (
				value      string
				successful bool
				result     interface{}
			)

			switch testCase.TestName {
			case "cluster":
				result, successful = httputils.ReadClusterName(recorder, request)
			case "rule_id":
				result, successful = httputils.ReadRuleID(recorder, request)
			case "error_key":
				result, successful = httputils.ReadErrorKey(recorder, request)
			case "organization":
				result, successful = httputils.ReadOrganizationID(recorder, request, false)
			case "organization/with_auth":
				result, successful = httputils.ReadOrganizationID(recorder, request, true)
			case "clusters":
				var results []types.ClusterName
				results, successful = httputils.ReadClusterNames(recorder, request)
				result = paramsToString(",", results)
			case "organizations":
				var results []types.OrgID
				results, successful = httputils.ReadOrganizationIDs(recorder, request)
				result = paramsToString(",", results)
			}

			assert.True(t, successful)

			value = fmt.Sprint(result)
			assert.Equal(t, expectedParamValue, value)
		})
	}
}

func TestReadClusterName_Error(t *testing.T) {
	for _, testCase := range []struct {
		TestCaseName  string
		Args          map[string]string
		ExpectedError string
	}{
		{TestCaseName: "Missing", Args: nil, ExpectedError: `{"status":"Missing required param from request: cluster"}`},
		{
			TestCaseName: "BadClusterName",
			Args: map[string]string{
				"cluster": string(testdata.BadClusterName),
			},
			ExpectedError: `{"status":"Error during parsing param 'cluster' with value '` +
				string(testdata.BadClusterName) + `'. Error: 'invalid UUID length: 4'"}`,
		},
	} {
		t.Run(testCase.TestCaseName, func(t *testing.T) {
			testReadParamError(t, "cluster", testCase.Args, testCase.ExpectedError)
		})
	}
}

func TestReadRuleID_Error(t *testing.T) {
	for _, testCase := range []struct {
		TestCaseName  string
		Args          map[string]string
		ExpectedError string
	}{
		{TestCaseName: "Missing", Args: nil, ExpectedError: `{"status":"Missing required param from request: rule_id"}`},
		{
			TestCaseName: "BadRuleID",
			Args: map[string]string{
				"rule_id": string(testdata.BadRuleID),
			},
			ExpectedError: `{"status":"Error during parsing param 'rule_id' with value '` +
				string(testdata.BadRuleID) +
				`'. Error: 'invalid rule ID, it must contain only from latin characters, number, underscores or dots'"}`,
		},
	} {
		t.Run(testCase.TestCaseName, func(t *testing.T) {
			testReadParamError(t, "rule_id", testCase.Args, testCase.ExpectedError)
		})
	}
}

func TestReadErrorKey_Error(t *testing.T) {
	testReadParamError(
		t,
		"error_key",
		nil,
		`{"status":"Missing required param from request: error_key"}`,
	)
}

func TestReadOrganization_Error(t *testing.T) {
	testReadParamError(
		t,
		"organization",
		nil,
		`{"status":"Missing required param from request: organization"}`,
	)
	testReadParamError(
		t,
		"organization/with_auth",
		map[string]string{
			"organization": fmt.Sprint(testdata.OrgID),
		},
		`{"status":"you have no permissions to get or change info about this organization"}`,
	)
}

func TestReadClusters_Error(t *testing.T) {
	for _, testCase := range []struct {
		TestCaseName  string
		Args          map[string]string
		ExpectedError string
	}{
		{TestCaseName: "Missing", Args: nil, ExpectedError: `{"status":"Missing required param from request: clusters"}`},
		{
			TestCaseName: "BadClusters",
			Args: map[string]string{
				"clusters": string(testdata.BadClusterName),
			},
			ExpectedError: `{"status":"Error during parsing param 'cluster' with value '` +
				string(testdata.BadClusterName) +
				`'. Error: 'invalid UUID length: 4'"}`,
		},
	} {
		t.Run(testCase.TestCaseName, func(t *testing.T) {
			testReadParamError(t, "clusters", testCase.Args, testCase.ExpectedError)
		})
	}
}

func TestReadOrganizations_Error(t *testing.T) {
	const badOrgID = "non-int"

	for _, testCase := range []struct {
		TestCaseName  string
		Args          map[string]string
		ExpectedError string
	}{
		{
			TestCaseName:  "Missing",
			Args:          nil,
			ExpectedError: `{"status":"Missing required param from request: organizations"}`,
		},
		{
			TestCaseName: "BadOrganizations",
			Args: map[string]string{
				"organizations": badOrgID,
			},
			ExpectedError: `{"status":"Error during parsing param 'organizations' with value '` + badOrgID +
				`'. Error: 'integer array expected'"}`,
		},
	} {
		t.Run(testCase.TestCaseName, func(t *testing.T) {
			testReadParamError(t, "organizations", testCase.Args, testCase.ExpectedError)
		})
	}
}

func mustGetRequestWithMuxVars(
	t *testing.T,
	method string,
	url string,
	body io.Reader,
	vars map[string]string,
) *http.Request {
	request, err := http.NewRequest(method, url, body)
	helpers.FailOnError(t, err)

	request = mux.SetURLVars(request, vars)

	return request
}

func testReadParamError(t *testing.T, paramName string, args map[string]string, expectedError string) {
	request := mustGetRequestWithMuxVars(t, http.MethodGet, "", nil, args)

	recorder := httptest.NewRecorder()

	var successful bool

	switch paramName {
	case "cluster":
		_, successful = httputils.ReadClusterName(recorder, request)
	case "rule_id":
		_, successful = httputils.ReadRuleID(recorder, request)
	case "error_key":
		_, successful = httputils.ReadErrorKey(recorder, request)
	case "organization":
		_, successful = httputils.ReadOrganizationID(recorder, request, false)
	case "organization/with_auth":
		ctx := context.WithValue(request.Context(), types.ContextKeyUser, types.Identity{
			AccountNumber: testdata.UserID,
			Internal: types.Internal{
				OrgID: testdata.Org2ID,
			},
		})
		request = request.WithContext(ctx)
		_, successful = httputils.ReadOrganizationID(recorder, request, true)
	case "organizations":
		_, successful = httputils.ReadOrganizationIDs(recorder, request)
	case "clusters":
		_, successful = httputils.ReadClusterNames(recorder, request)
	default:
		panic("testReadParamError is not implemented for param '" + paramName + "'")
	}

	assert.False(t, successful)

	resp := recorder.Result()
	assert.NotNil(t, resp)

	body, err := ioutil.ReadAll(resp.Body)
	helpers.FailOnError(t, err)

	assert.Equal(t, expectedError, strings.TrimSpace(string(body)))
}

func paramsToString(separator string, params ...interface{}) string {
	var unpackedParams []interface{}

	for _, param := range params {
		switch reflect.TypeOf(param).Kind() {
		case reflect.Array, reflect.Slice:
			s := reflect.ValueOf(param)

			for i := 0; i < s.Len(); i++ {
				unpackedParams = append(unpackedParams, s.Index(i).Interface())
			}
		default:
			unpackedParams = append(unpackedParams, param)
		}
	}

	params = unpackedParams

	var stringParams []string
	for _, param := range params {
		stringParams = append(stringParams, fmt.Sprint(param))
	}

	res := strings.Join(stringParams, separator)
	return res
}

// TestReadClusterListFromPathMissingClusterList function checks if missing
// cluster list in path is detected and processed correctly by function
// ReadClusterListFromPath.
func TestReadClusterListFromPathMissingClusterList(t *testing.T) {
	request, err := http.NewRequest(http.MethodGet, "", nil)
	helpers.FailOnError(t, err)

	// try to read list of clusters from path
	_, successful := httputils.ReadClusterListFromPath(httptest.NewRecorder(), request)

	// missing list means that the read operation should fail
	assert.False(t, successful)
}

// TestReadClusterListFromPathEmptyClusterList function checks if empty cluster
// list in path is detected and processed correctly by function
// ReadClusterListFromPath.
func TestReadClusterListFromPathEmptyClusterList(t *testing.T) {
	request := mustGetRequestWithMuxVars(t, http.MethodGet, "", nil, map[string]string{
		"cluster_list": "",
	})

	// try to read list of clusters from path
	_, successful := httputils.ReadClusterListFromPath(httptest.NewRecorder(), request)

	// empty list means that the read operation should fail
	assert.False(t, successful)
}

// TestReadClusterListFromPathOneCluster function checks if list with one
// cluster ID is processed correctly by function ReadClusterListFromPath.
func TestReadClusterListFromPathOneCluster(t *testing.T) {
	request := mustGetRequestWithMuxVars(t, http.MethodGet, "", nil, map[string]string{
		"cluster_list": fmt.Sprintf("%v", cluster1ID),
	})

	// try to read list of clusters from path
	list, successful := httputils.ReadClusterListFromPath(httptest.NewRecorder(), request)

	// cluster list exists so the read operation should not fail
	assert.True(t, successful)

	// we expect do get list with one cluster ID
	assert.ElementsMatch(t, list, []string{cluster1ID})
}

// TestReadClusterListFromPathTwoClusters function checks if list with two
// cluster IDs is processed correctly by function ReadClusterListFromPath.
func TestReadClusterListFromPathTwoClusters(t *testing.T) {
	request := mustGetRequestWithMuxVars(t, http.MethodGet, "", nil, map[string]string{
		"cluster_list": fmt.Sprintf("%v,%v", cluster1ID, cluster2ID),
	})

	// try to read list of clusters from path
	list, successful := httputils.ReadClusterListFromPath(httptest.NewRecorder(), request)

	// cluster list exists so the read operation should not fail
	assert.True(t, successful)

	// we expect do get list with two cluster IDs
	assert.ElementsMatch(t, list, []string{cluster1ID, cluster2ID})
}

// TestReadClusterListFromBodyNoJSON function checks if reading list of
// clusters from empty request body is detected properly by function
// ReadClusterListFromBody.
func TestReadClusterListFromBodyNoJSON(t *testing.T) {
	request, err := http.NewRequest(
		http.MethodGet,
		"",
		strings.NewReader(""),
	)
	helpers.FailOnError(t, err)

	// try to read list of clusters from path
	_, successful := httputils.ReadClusterListFromBody(httptest.NewRecorder(), request)

	// the read should fail because of empty request body
	assert.False(t, successful)
}

// TestReadClusterListFromBodyCorrectJSON function checks if reading list of
// clusters from correct request body containing JSON data is done correctly by
// function ReadClusterListFromBody.
func TestReadClusterListFromBodyCorrectJSON(t *testing.T) {
	request, err := http.NewRequest(
		http.MethodGet,
		"",
		strings.NewReader(fmt.Sprintf(`{"clusters": ["%v","%v"]}`, cluster1ID, cluster2ID)),
	)
	helpers.FailOnError(t, err)

	// try to read list of clusters from path
	list, successful := httputils.ReadClusterListFromBody(httptest.NewRecorder(), request)

	// cluster list exists so the call should not fail
	assert.True(t, successful)

	// we expect do get list with two cluster IDs
	assert.ElementsMatch(t, list, []string{cluster1ID, cluster2ID})
}

// TestReadClusterListFromBodyWrongJSON function checks if reading list of
// clusters from request body with improper format is processed correctly by
// function ReadClusterListFromBody.
func TestReadClusterListFromBodyWrongJSON(t *testing.T) {
	request, err := http.NewRequest(
		http.MethodGet,
		"",
		strings.NewReader("this-is-not-json"),
	)
	helpers.FailOnError(t, err)

	// try to read list of clusters from path
	_, successful := httputils.ReadClusterListFromBody(httptest.NewRecorder(), request)

	// the read should fail because of broken JSON
	assert.False(t, successful)
}
