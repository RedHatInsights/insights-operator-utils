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

package parsers_test

// Documentation in literate-programming-style is available at:
// https://redhatinsights.github.io/insights-operator-utils/packages/parsers/rule_test.html

import (
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/RedHatInsights/insights-operator-utils/parsers"
	"github.com/RedHatInsights/insights-operator-utils/types"
)

// TestParseRuleSelector checks the function parsers.ParseRuleSelector for
// valid input.
func TestParseRuleSelector(t *testing.T) {
	t.Run("rule selector without numbers", func(t *testing.T) {
		component, errorKey, err := parsers.ParseRuleSelector("foo|bar")

		assert.Nil(t, err, "unexpected error")
		assert.Equal(t, types.Component("foo"), component)
		assert.Equal(t, types.ErrorKey("bar"), errorKey)
	})

	t.Run("rule selector with numbers", func(t *testing.T) {
		component, errorKey, err := parsers.ParseRuleSelector("foo1|bar2")

		assert.Nil(t, err, "unexpected error")
		assert.Equal(t, types.Component("foo1"), component)
		assert.Equal(t, types.ErrorKey("bar2"), errorKey)
	})
}

// TestParseRuleSelectorEmptyInput checks the function
// parsers.ParseRuleSelector for empty input.
func TestParseRuleSelectorEmptyInput(t *testing.T) {
	t.Run("empty input", func(t *testing.T) {
		_, _, err := parsers.ParseRuleSelector("")

		assert.Equal(t, err.Error(), "invalid rule ID, it must contain only rule ID and error key separated by |")
	})
}

// TestParseRuleSelectorMoreSeparators checks the function
// parsers.ParseRuleSelector for input with more separators.
func TestParseRuleSelectorMoreSeparators(t *testing.T) {
	t.Run("foo|bar|baz", func(t *testing.T) {
		_, _, err := parsers.ParseRuleSelector("")

		assert.Equal(t, err.Error(), "invalid rule ID, it must contain only rule ID and error key separated by |")
	})
}

// TestParseRuleSelectorImproperInput checks the function
// parsers.ParseRuleSelector for improper input.
func TestParseRuleSelectorImproperInput(t *testing.T) {
	t.Run("rule without error key", func(t *testing.T) {
		_, _, err := parsers.ParseRuleSelector("foo|")

		assert.Equal(t, err.Error(), "invalid rule ID, each part of ID must contain only latin characters, number, underscores or dots")
	})
	t.Run("rule without component", func(t *testing.T) {
		_, _, err := parsers.ParseRuleSelector("|bar")

		assert.Equal(t, err.Error(), "invalid rule ID, each part of ID must contain only latin characters, number, underscores or dots")
	})
	t.Run("rule with improper characters", func(t *testing.T) {
		_, _, err := parsers.ParseRuleSelector("ěšč|řžý")

		assert.Equal(t, err.Error(), "invalid rule ID, each part of ID must contain only latin characters, number, underscores or dots")
	})
}
