// Copyright 2021 Red Hat, Inc
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

package parsers

import (
	"fmt"
	"regexp"
	"strings"

	types "github.com/RedHatInsights/insights-results-types"
)

// ParseRuleSelector function parses the rule selector which consists of
// component name and error key separaed by "|" character.
func ParseRuleSelector(ruleSelector types.RuleSelector) (types.Component, types.ErrorKey, error) {
	// component and error key is to be separated by "|"
	splitedRuleID := strings.Split(string(ruleSelector), "|")

	// check if both parts have been found
	if len(splitedRuleID) != 2 {
		err := fmt.Errorf("invalid rule ID, it must contain only rule ID and error key separated by |")
		return types.Component(""), types.ErrorKey(""), err
	}

	// check component name and error key content
	IDValidator := regexp.MustCompile(`^[a-zA-Z_0-9.]+$`)

	isRuleIDValid := IDValidator.MatchString(splitedRuleID[0])
	isErrorKeyValid := IDValidator.MatchString(splitedRuleID[1])

	if !isRuleIDValid {
		err := fmt.Errorf("invalid component name: each part of ID must contain only latin characters, number, underscores or dots")
		return types.Component(""), types.ErrorKey(""), err
	}

	if !isErrorKeyValid {
		err := fmt.Errorf("invalid error key: each part of ID must contain only latin characters, number, underscores or dots")
		return types.Component(""), types.ErrorKey(""), err
	}

	return types.Component(splitedRuleID[0]), types.ErrorKey(splitedRuleID[1]), nil
}
