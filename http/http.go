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

// HTTP-related utility functions.
package httputils

// Documentation in literate-programming-style is available at:
// https://redhatinsights.github.io/insights-operator-utils/packages/http/http.html

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

// MakeURLToEndpoint creates URL to endpoint, use constants from file endpoints.go
func MakeURLToEndpoint(apiPrefix, endpoint string, args ...interface{}) string {
	endpoint = ReplaceParamsInEndpointAndTrimLeftSlash(endpoint, "%v")

	if apiPrefix != "/" && len(endpoint) > 0 {
		apiPrefix = strings.TrimRight(apiPrefix, "/")
	}

	nonParsedURL := apiPrefix
	endpointWithArgs := fmt.Sprintf(endpoint, args...)
	if len(endpointWithArgs) > 0 {
		nonParsedURL += "/" + endpointWithArgs
	}

	resultingURL, err := url.Parse(nonParsedURL)

	if err != nil {
		return nonParsedURL
	}

	return resultingURL.String()
}

// ReplaceParamsInEndpointAndTrimLeftSlash replaces params in endpoint and trims left slash
func ReplaceParamsInEndpointAndTrimLeftSlash(endpoint, replacer string) string {
	re := regexp.MustCompile(`\{[a-zA-Z_0-9]+\}`)

	endpoint = re.ReplaceAllString(endpoint, replacer)
	endpoint = strings.TrimLeft(endpoint, "/")

	return endpoint
}

// MakeURLToEndpointMapString creates URL to endpoint using arguments in map in string format, use constants from file endpoints.go
func MakeURLToEndpointMapString(apiPrefix, endpoint string, args map[string]string) string {
	newArgs := make(map[string]interface{})

	for key, val := range args {
		newArgs[key] = val
	}

	return MakeURLToEndpointMap(apiPrefix, endpoint, newArgs)
}

// MakeURLToEndpointMap creates URL to endpoint using arguments in map, use constants from file endpoints.go
func MakeURLToEndpointMap(apiPrefix, endpoint string, args map[string]interface{}) string {
	endpoint = strings.TrimLeft(endpoint, "/")
	for key, val := range args {
		endpoint = strings.ReplaceAll(endpoint, fmt.Sprintf("{%v}", key), fmt.Sprint(val))
	}

	apiPrefix = strings.TrimRight(apiPrefix, "/")

	return apiPrefix + "/" + endpoint
}

// SetHTTPPrefix adds HTTP prefix if it is not already present in the given string
func SetHTTPPrefix(url string) string {
	if !strings.HasPrefix(url, "http") {
		// if no protocol is specified in given URL, assume it is not
		// needed to use https
		url = "http://" + url
	}
	return url
}

// SendRequest sends the given request, reads the body and handles related errors
func SendRequest(req *http.Request, timeout time.Duration) ([]byte, error) {
	client := &http.Client{
		Timeout: timeout,
	}

	response, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	// Read body from response
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	err = response.Body.Close()
	if err != nil {
		return nil, err
	}
	return body, nil
}
