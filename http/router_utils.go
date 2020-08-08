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

package httputils

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"

	"github.com/RedHatInsights/insights-operator-utils/types"
)

// getRouterParam retrieves parameter from URL like `/organization/{org_id}`
func getRouterParam(request *http.Request, paramName string) (string, error) {
	value, found := mux.Vars(request)[paramName]
	if !found {
		return "", &types.RouterMissingParamError{ParamName: paramName}
	}

	return value, nil
}

// validateClusterName checks that the cluster name is a valid UUID.
// Converted cluster name is returned if everything is okay, otherwise an error is returned.
func validateClusterName(clusterName string) (types.ClusterName, error) {
	if _, err := uuid.Parse(clusterName); err != nil {
		message := fmt.Sprintf("invalid cluster name: '%s'. Error: %s", clusterName, err.Error())

		log.Error().Err(err).Msg(message)

		return "", &types.RouterParsingError{
			ParamName:  "cluster",
			ParamValue: clusterName,
			ErrString:  err.Error(),
		}
	}

	return types.ClusterName(clusterName), nil
}

func handleClusterNameError(writer http.ResponseWriter, err error) {
	log.Error().Msg(err.Error())

	// query parameter 'cluster' can't be found in request, which might be caused by issue in Gorilla mux
	// (not on client side), but let's assume it won't :)
	types.HandleServerError(writer, err)
}

// ReadClusterName retrieves cluster name from request
// if it's not possible, it writes http error to the writer and returns false
func ReadClusterName(writer http.ResponseWriter, request *http.Request) (types.ClusterName, bool) {
	clusterName, err := getRouterParam(request, "cluster")
	if err != nil {
		handleClusterNameError(writer, err)
		return "", false
	}

	validatedClusterName, err := validateClusterName(clusterName)
	if err != nil {
		handleClusterNameError(writer, err)
		return "", false
	}

	return validatedClusterName, true
}

// ReadErrorKey retrieves rule id from request's url or writes an error to writer.
// The function returns a rule id and a bool indicating if it was successful.
func ReadRuleID(writer http.ResponseWriter, request *http.Request) (types.RuleID, bool) {
	ruleID, err := getRouterParam(request, "rule_id")
	if err != nil {
		const message = "unable to get rule id"
		log.Error().Err(err).Msg(message)
		types.HandleServerError(writer, err)
		return types.RuleID(0), false
	}

	ruleIDValidator := regexp.MustCompile(`^[a-zA-Z_0-9.]+$`)

	isRuleIDValid := ruleIDValidator.Match([]byte(ruleID))

	if !isRuleIDValid {
		err = fmt.Errorf("invalid rule ID, it must contain only from latin characters, number, underscores or dots")
		log.Error().Err(err)
		types.HandleServerError(writer, &types.RouterParsingError{
			ParamName:  "rule_id",
			ParamValue: ruleID,
			ErrString:  err.Error(),
		})
		return types.RuleID(0), false
	}

	return types.RuleID(ruleID), true
}

// ReadErrorKey retrieves error key from request's url or writes an error to writer.
// The function returns an error key and a bool indicating if it was successful.
func ReadErrorKey(writer http.ResponseWriter, request *http.Request) (types.ErrorKey, bool) {
	errorKey, err := getRouterParam(request, "error_key")
	if err != nil {
		const message = "unable to get error_key"
		log.Error().Err(err).Msg(message)
		types.HandleServerError(writer, err)
		return types.ErrorKey(0), false
	}

	return types.ErrorKey(errorKey), true
}
