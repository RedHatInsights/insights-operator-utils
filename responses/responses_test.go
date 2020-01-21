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
	//"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

var mock_payload = map[string]interface{}{
	"numbers_l": []int{1, 2, 3, 5, 8},
	"color_s":   "blue",
	"extra_data_m": map[string]interface{}{
		"param1": 1,
		"param2": false,
		"param3": 0.75,
	},
}

// Test BuildResponse that returns simple map with key "status" and given value
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

// Test BuildResponse that returns simple map with key "status" and value "ok"
func TestBuildOkResponse(t *testing.T) {
	expected_response := map[string]interface{}{
		"status": "ok",
	}
	response := responses.BuildOkResponse()
	if equal := reflect.DeepEqual(expected_response, response); !equal {
		t.Errorf("Expected response %v", expected_response)
	}
}

// Test BuildResponse that returns simple map with key status: ok and given data
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

func TestSendResponse(t *testing.T) {
	test_server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := responses.BuildOkResponseWithData("data", mock_payload)
		responses.SendResponse(w, response)
	}))
	defer test_server.Close()

	res, err := http.Get(test_server.URL)
	if err != nil {
		t.Fatal(err)
	}

	var response map[string]interface{}

	json.NewDecoder(res.Body).Decode(&response)
	expected_response := map[string]interface{}{
		"data":   mock_payload,
		"status": "ok",
	}

	t.Log(response)
	t.Log(expected_response)
	// why aren't they equal???
	if equal := reflect.DeepEqual(expected_response, response); !equal {
		t.Errorf("Expected response %v.", expected_response)
	}
}

func TestSendCreated(t *testing.T) {

}

func TestSendAccepted(t *testing.T) {

}

func TestSendError(t *testing.T) {

}

func TestSendForbidden(t *testing.T) {

}

func TestSendInternalServerError(t *testing.T) {

}

func TestSendUnauthorized(t *testing.T) {

}

/*
   content_type := req.Header.Get("Content-type")
   if content_type != "application/json" {
       http.Error(w, fmt.Sprintf("Unexpected content %s", content_type), http.StatusBadRequest)
   }
*/
