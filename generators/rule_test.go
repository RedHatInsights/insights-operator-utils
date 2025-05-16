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

package generators_test

// Documentation in literate-programming-style is available at:
// https://redhatinsights.github.io/insights-operator-utils/packages/generators/rule_test.html

import (
	"testing"

	"github.com/stretchr/testify/assert"

	types "github.com/RedHatInsights/insights-results-types"

	"github.com/RedHatInsights/insights-operator-utils/generators"
)

const (
	ruleFQDNOK = "rule.module.name"
	errorKeyOK = "ERROR_KEY"
)

// TestGenerateCompositeRuleID checks the function generators.GenerateCompositeRuleID
func TestGenerateCompositeRuleID(t *testing.T) {
	t.Run("everything fine", func(t *testing.T) {
		ruleID, err := generators.GenerateCompositeRuleID(types.RuleFQDN(ruleFQDNOK), types.ErrorKey(errorKeyOK))

		assert.Nil(t, err)
		assert.Equal(t, types.RuleID("rule.module.name|ERROR_KEY"), ruleID)
	})

	t.Run("ruleFQDN empty", func(t *testing.T) {
		_, err := generators.GenerateCompositeRuleID(types.RuleFQDN(""), types.ErrorKey(errorKeyOK))
		assert.Error(t, err)
		assert.Equal(t, err.Error(), "empty rule FQDN")
	})

	t.Run("error key empty", func(t *testing.T) {
		_, err := generators.GenerateCompositeRuleID(types.RuleFQDN(ruleFQDNOK), types.ErrorKey(""))
		assert.Error(t, err)
		assert.Equal(t, err.Error(), "empty error key")
	})

	t.Run("ruleFQDN and error key empty", func(t *testing.T) {
		_, err := generators.GenerateCompositeRuleID(types.RuleFQDN(""), types.ErrorKey(""))
		assert.Error(t, err)
		assert.Equal(t, err.Error(), "empty rule FQDN")
	})
}
