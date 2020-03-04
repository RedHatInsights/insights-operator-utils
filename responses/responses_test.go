/*
Copyright © 2019, 2020 Red Hat, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package responses_test

import (
	"encoding/json"
	"github.com/RedHatInsights/insights-operator-utils/responses"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

// Define types used in table tests struct
type functionWithData func(http.ResponseWriter, map[string]interface{}) error

type functionWithoutData func(http.ResponseWriter, string) error

const (
	contentType = "Content-Type"
	appJSON     = "application/json; charset=utf-8"
)

// Mock payload to be sent as data
var mockPayload = map[string]interface{}{
	"color_s": "blue",
	"extra_data_m": map[string]interface{}{
		"param1": 1,
		"param2": false,
	},
}

/*
Table tests - two structs are used so that different type of function can be used
It would be possible to only have one struct with both func types and a flag e.g. testPayload bool
determining which one to call.
*/
var headerTestsWithData = []struct {
	testName       string
	fName          functionWithData
	expectedHeader int
}{
	{"responses.SendResponse", responses.SendResponse, http.StatusOK},
	{"responses.SendCreated", responses.SendCreated, http.StatusCreated},
	{"responses.SendAccepted", responses.SendAccepted, http.StatusAccepted},
	{"responses.SendUnauthorized", responses.SendUnauthorized, http.StatusUnauthorized},
}

var headerTestsWithoutData = []struct {
	testName       string
	fName          functionWithoutData
	expectedHeader int
}{
	{"responses.SendError", responses.SendError, http.StatusBadRequest},
	{"responses.SendForbidden", responses.SendForbidden, http.StatusForbidden},
	{"responses.SendInternalServerError", responses.SendInternalServerError, http.StatusInternalServerError},
}

var sendTests = []struct {
	testName       string
	statusCode     int
	expectedStatus int
	data           interface{}
	expectedData   string
}{
	{
		"responses.Send(http.StatusInternalServerError, ...)",
		http.StatusInternalServerError,
		http.StatusInternalServerError,
		"unable to connect to the database",
		`{"status": "unable to connect to the database"}`,
	},
	{
		"responses.Send(http.StatusInternalServerError, ...)",
		http.StatusBadRequest,
		http.StatusBadRequest,
		"wrong ID format",
		`{"status": "wrong ID format"}`,
	},
	{
		"responses.Send(http.StatusInternalServerError, ...)",
		http.StatusOK,
		http.StatusOK,
		map[string]interface{}{"list_of_something": []string{}},
		`{"list_of_something": []}`,
	},
}

// checkResponse checks status code against the expected one and check content-type + payload
func checkResponse(
	url string,
	expectedStatusCode int,
	checkPayload bool,
	expectedBody string,
	t *testing.T,
) {
	res, err := http.Get(url)
	if err != nil {
		t.Fatal(err)
	}

	if res.StatusCode != expectedStatusCode {
		t.Errorf("Expected status code %v, got %v", expectedStatusCode, res.StatusCode)
	}

	if checkPayload {
		contentType := res.Header.Get(contentType)
		if contentType != appJSON {
			t.Errorf("Unexpected content type. Expected %v, got %v", appJSON, contentType)
		}

		body, err := ioutil.ReadAll(res.Body)
		defer res.Body.Close()

		var expected map[string]interface{}
		err = json.NewDecoder(strings.NewReader(expectedBody)).Decode(&expected)
		if err != nil {
			t.Fatal(err)
		}

		var response map[string]interface{}
		err = json.Unmarshal(body, &response)
		if err != nil {
			t.Fatal(err)
		}

		if equal := reflect.DeepEqual(response, expected); !equal {
			t.Errorf(`Expected response "%+v", got "%+v"`, expected, response)
		}
	}
}

// TestBuildResponse tests BuildResponse func that returns simple map with key "status" and given value
func TestBuildResponse(t *testing.T) {
	statusStr := "I'm a teapot"
	expectedResponse := map[string]interface{}{
		"status": statusStr,
	}
	response := responses.BuildResponse(statusStr)
	if equal := reflect.DeepEqual(expectedResponse, response); !equal {
		t.Errorf("Expected response %v", expectedResponse)
	}
}

// TestBuildOkResponse tests BuildResponse that returns simple map with key "status" and value "ok"
func TestBuildOkResponse(t *testing.T) {
	expectedResponse := map[string]interface{}{
		"status": "ok",
	}
	response := responses.BuildOkResponse()
	if equal := reflect.DeepEqual(expectedResponse, response); !equal {
		t.Errorf("Expected response %v", expectedResponse)
	}
}

// TestBuildOkResponseWithData tests that the func returns simple map with key status: ok and given data
func TestBuildOkResponseWithData(t *testing.T) {
	expectedResponse := map[string]interface{}{
		"data":   mockPayload,
		"status": "ok",
	}
	response := responses.BuildOkResponseWithData("data", mockPayload)
	if equal := reflect.DeepEqual(expectedResponse, response); !equal {
		t.Errorf("Expected response %v", expectedResponse)
	}
}

// TestHeaders run table tests to test StatusCodes and payloads.
func TestHeaders(t *testing.T) {
	for _, tt := range headerTestsWithData {
		t.Run(tt.testName, func(t *testing.T) {
			testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				err := tt.fName(w, mockPayload) // call the function
				if err != nil {
					t.Fatal(err)
				}
			}))
			defer testServer.Close()

			const expectedBody = `{"color_s":"blue","extra_data_m":{"param1":1,"param2":false}}`
			checkResponse(
				testServer.URL,
				tt.expectedHeader,
				true,
				expectedBody,
				t,
			)
		})
	}

	for _, tt := range headerTestsWithoutData {
		t.Run(tt.testName, func(t *testing.T) {
			testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				err := tt.fName(w, "Test Status") // call the function
				if err != nil {
					t.Fatal(err)
				}
			}))
			defer testServer.Close()

			checkResponse(
				testServer.URL,
				tt.expectedHeader,
				false,
				"",
				t,
			)
		})
	}

	for _, test := range sendTests {
		t.Run(test.testName, func(t *testing.T) {
			testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				responses.Send(test.statusCode, w, test.data)
			}))
			defer testServer.Close()

			checkResponse(
				testServer.URL,
				test.expectedStatus,
				true,
				test.expectedData,
				t,
			)
		})
	}
}
