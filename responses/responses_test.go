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
type functionWithData func(http.ResponseWriter, map[string]interface{})

type functionWithoutData func(http.ResponseWriter, string)

const (
	expectedBody = `{"color_s":"blue","extra_data_m":{"param1":1,"param2":false}}`
	contentType  = "Content-Type"
	appJSON      = "application/json; charset=utf-8"
)

// Mock payload to be sent as data
var mock_payload = map[string]interface{}{
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

// checkResponse checks status code against the expected one and check content-type + payload
func checkResponse(url string, expectedStatusCode int, checkPayload bool, t *testing.T) {
	res, err := http.Get(url)
	if err != nil {
		t.Fatal(err)
	}

	if res.StatusCode != expectedStatusCode {
		t.Errorf("Expected status code %v, got %v", expectedStatusCode, res.StatusCode)
	}

	if checkPayload {
		content_type := res.Header.Get(contentType)
		if content_type != appJSON {
			t.Errorf("Unexpected content type. Expected %v, got %v", appJSON, content_type)
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
			t.Errorf("Expected response %v.", expected)
		}
	}
}

// TestBuildResponse tests BuildResponse func that returns simple map with key "status" and given value
func TestBuildResponse(t *testing.T) {
	status_str := "I'm a teapot"
	expected_response := map[string]interface{}{
		"status": status_str,
	}
	response := responses.BuildResponse(status_str)
	if equal := reflect.DeepEqual(expected_response, response); !equal {
		t.Errorf("Expected response %v", expected_response)
	}
}

// TestBuildOkResponse tests BuildResponse that returns simple map with key "status" and value "ok"
func TestBuildOkResponse(t *testing.T) {
	expected_response := map[string]interface{}{
		"status": "ok",
	}
	response := responses.BuildOkResponse()
	if equal := reflect.DeepEqual(expected_response, response); !equal {
		t.Errorf("Expected response %v", expected_response)
	}
}

// TestBuildOkResponseWithData tests that the func returns simple map with key status: ok and given data
func TestBuildOkResponseWithData(t *testing.T) {
	expected_response := map[string]interface{}{
		"data":   mock_payload,
		"status": "ok",
	}
	response := responses.BuildOkResponseWithData("data", mock_payload)
	if equal := reflect.DeepEqual(expected_response, response); !equal {
		t.Errorf("Expected response %v", expected_response)
	}
}

// TestHeaders run table tests to test StatusCodes and payloads.
func TestHeaders(t *testing.T) {
	for _, tt := range headerTestsWithData {
		t.Run(tt.testName, func(t *testing.T) {
			test_server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				tt.fName(w, mock_payload) // call the function
			}))
			defer test_server.Close()

			checkResponse(test_server.URL, tt.expectedHeader, true, t)
		})
	}

	for _, tt := range headerTestsWithoutData {
		t.Run(tt.testName, func(t *testing.T) {
			test_server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				tt.fName(w, "Test Status") // call the function
			}))
			defer test_server.Close()

			checkResponse(test_server.URL, tt.expectedHeader, false, t)
		})
	}
}
