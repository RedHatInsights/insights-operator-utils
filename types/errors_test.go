// Copyright 2023 Red Hat, Inc
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

package types_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/RedHatInsights/insights-operator-utils/types"
)

// TestRouterMissingParamError checks the method Error() for data structure
// RouterMissingParamError
func TestRouterMissingParamError(t *testing.T) {
	// expected error value
	const expected = "Missing required param from request: paramName"

	// construct an instance of error interface
	err := types.RouterMissingParamError{
		ParamName: "paramName",
	}

	// check if error value is correct
	assert.Equal(t, err.Error(), expected)
}

// TestRouterParsingError checks the method Error() for data structure
// RouterParsingError
func TestRouterParsingError(t *testing.T) {
	// expected error value
	const expected = "Error during parsing param 'paramName' with value 'paramValue'. Error: 'errorMessage'"

	// construct an instance of error interface
	err := types.RouterParsingError{
		ParamName:  "paramName",
		ParamValue: "paramValue",
		ErrString:  "errorMessage"}

	// check if error value is correct
	assert.Equal(t, err.Error(), expected)
}

// TestNoContentError checks the method Error() for data structure
// NoContentError
func TestNoContentError(t *testing.T) {
	// expected error value
	const errorString = "error message"
	const expected = errorString

	// construct an instance of error interface
	err := types.NoContentError{
		ErrString: errorString,
	}

	// check if error value is correct
	assert.Equal(t, err.Error(), expected)
}

// TestUnauthorizedError checks the method Error() for data structure
// UnauthorizedError
func TestUnauthorizedError(t *testing.T) {
	// expected error value
	const errorString = "unauthorized error message"
	const expected = errorString

	// construct an instance of error interface
	err := types.UnauthorizedError{
		ErrString: errorString,
	}

	// check if error value is correct
	assert.Equal(t, err.Error(), expected)
}

// TestForbiddenError checks the method Error() for data structure
// ForbiddenError
func TestForbiddenError(t *testing.T) {
	// expected error value
	const errorString = "forbidden error message"
	const expected = errorString

	// construct an instance of error interface
	err := types.ForbiddenError{
		ErrString: errorString,
	}

	// check if error value is correct
	assert.Equal(t, err.Error(), expected)
}

// TestNoBodyError checks the method Error() for data structure
// NoBodyError
func TestNoBodyError(t *testing.T) {
	// expected error value
	const expected = "client didn't provide request body"

	// construct an instance of error interface
	err := types.NoBodyError{}

	// check if error value is correct
	assert.Equal(t, err.Error(), expected)
}

// TestValidationError checks the method Error() for data structure
// ValidationError
func TestValidationError(t *testing.T) {
	// expected error value
	const expected = "Error during validating param 'PARAMETER' with value 'VALUE'. Error: 'ERROR'"

	// construct an instance of error interface
	err := types.ValidationError{
		ParamName:  "PARAMETER",
		ParamValue: "VALUE",
		ErrString:  "ERROR",
	}

	// check if error value is correct
	assert.Equal(t, err.Error(), expected)
}

// TestItemNotFoundError checks the method Error() for data structure
// ItemNotFoundError.
func TestItemNotFoundError(t *testing.T) {
	// expected error value
	const expected = "Item with ID ITEM_ID was not found in the storage"

	// construct an instance of error interface
	err := types.ItemNotFoundError{
		ItemID: "ITEM_ID"}

	// check if error value is correct
	assert.Equal(t, err.Error(), expected)
}

// TestHandleServer error check the function HandleServerError defined in errors.go
func TestHandleServerError(t *testing.T) {
	// check the behaviour with all error types defined in this package
	testResponse(t, &types.RouterMissingParamError{}, http.StatusBadRequest)
	testResponse(t, &types.RouterParsingError{}, http.StatusBadRequest)
	testResponse(t, &types.ItemNotFoundError{}, http.StatusNotFound)
	testResponse(t, &types.UnauthorizedError{}, http.StatusUnauthorized)
	testResponse(t, &types.ForbiddenError{}, http.StatusForbidden)
	testResponse(t, &types.ForbiddenError{}, http.StatusForbidden)

	// we need to retriev json.UnmarshalTypeError
	// so let's try to unmarshal "foo" string into an integer
	var x int
	err := json.Unmarshal([]byte("\"foo\""), &x)

	/// test with json.UnmarshalTypeError
	testResponse(t, err, http.StatusBadRequest)

	// error can be nil
	testResponse(t, nil, http.StatusInternalServerError)
}

// testResponse function uses HTTP server mock to check server response
// handlers
func testResponse(t *testing.T, e error, expectedCode int) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		types.HandleServerError(w, e)
	}))
	defer testServer.Close()

	res, err := http.Get(testServer.URL)
	if err != nil {
		t.Fatal(err)
	}

	if res.StatusCode != expectedCode {
		t.Errorf("Expected status code %v but got %v", expectedCode, res.StatusCode)
	}
}
