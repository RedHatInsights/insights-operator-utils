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

package generators

import (
	"fmt"

	types "github.com/RedHatInsights/insights-results-types"
)

// GenerateCompositeRuleID generates a rule ID in the "rule.module|ERROR_KEY" format
// TODO: validate the rule FQDN and error key properly
func GenerateCompositeRuleID(ruleFQDN types.RuleFQDN, errorKey types.ErrorKey) (
	ruleID types.RuleID,
	err error,
) {
	// check if ruleFQDN is not empty
	if len(ruleFQDN) == 0 {
		err = fmt.Errorf("empty rule FQDN")
		return
	}

	// check if error key is not empty
	if len(errorKey) == 0 {
		err = fmt.Errorf("empty error key")
		return
	}

	ruleID = types.RuleID(fmt.Sprintf("%v|%v", ruleFQDN, errorKey))

	return
}
