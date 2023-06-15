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
