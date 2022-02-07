// Copyright 2020, 2021, 2022 Red Hat, Inc
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

package httputils

// Documentation in literate-programming-style is available at:
// https://redhatinsights.github.io/insights-operator-utils/packages/http/openapi.html

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/rs/zerolog/log"

	"github.com/RedHatInsights/insights-operator-utils/types"
)

// FilterOutDebugMethods returns the same openapi spec, but without endpoints tagged as debug.
func FilterOutDebugMethods(openAPIFileContent string) (string, error) {
	swagger, err := openapi3.NewSwaggerLoader().LoadSwaggerFromData([]byte(openAPIFileContent))
	if err != nil {
		return "", err
	}

	newPaths := make(openapi3.Paths)

	for path, pathItem := range swagger.Paths {
		for method, operation := range pathItem.Operations() {
			debugTagFound := false
			for _, tag := range operation.Tags {
				if strings.EqualFold(strings.TrimSpace(tag), "debug") {
					debugTagFound = true
					break
				}
			}

			if debugTagFound {
				pathItem.SetOperation(method, nil)
			}
		}

		if len(pathItem.Operations()) > 0 {
			newPaths[path] = pathItem
		}
	}

	swagger.Paths = newPaths

	openAPIBytes, err := swagger.MarshalJSON()
	return string(openAPIBytes), err
}

// CreateOpenAPIHandler creates a handler for a server to send OpenAPI file.
// Optionally, you can turn on or off debug to filter out debug endpoints.
// Optionally, you can turn on caching by setting cacheFile to true,
// then you will have to restart a server on each file change
func CreateOpenAPIHandler(filePath string, debug, cacheFile bool) func(writer http.ResponseWriter, request *http.Request) {
	var fileContent []byte

	return func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")

		if !cacheFile || len(fileContent) == 0 {
			var err error
			// it's not supposed that we'll accept the path from a user
			fileContent, err = ioutil.ReadFile(filePath) // #nosec G304  (CWE-22): Potential file inclusion via variable
			if err != nil {
				log.Error().Err(err).Msg("error reading openapi.json file")
				types.HandleServerError(writer, err)
				return
			}

			if !debug {
				filteredFileContent, err := FilterOutDebugMethods(string(fileContent))
				if err != nil {
					log.Error().Err(err).Msg("error filtering openapi.json file, using the original version")
				} else {
					fileContent = []byte(filteredFileContent)
				}
			}
		}

		_, err := writer.Write(fileContent)
		if err != nil {
			log.Error().Err(err).Msg("error writing openapi.json file")
			types.HandleServerError(writer, err)
			return
		}
	}
}
