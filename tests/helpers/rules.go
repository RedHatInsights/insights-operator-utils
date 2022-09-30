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

package helpers

// Documentation in literate-programming-style is available at:
// https://redhatinsights.github.io/insights-operator-utils/packages/tests/helpers/http.html

import (
	"sort"
	"testing"
	"time"

	types "github.com/RedHatInsights/insights-results-types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// CompareReportResponses compares two RuleOnReport struct field by field,
// except for the CreatedAt field that is compared with createdAt timestamp. If
// createdAt.isZero() the filed will be ignored
func CompareReportResponses(t testing.TB, expected, actual types.RuleOnReport, createdAt time.Time) {
	actualTemplateData := ToJSONString(actual.TemplateData)
	expectedTemplateData := ToJSONString(expected.TemplateData)
	require.JSONEq(t, expectedTemplateData, actualTemplateData)
	assert.Equal(t, actual.Disabled, expected.Disabled)
	assert.Equal(t, actual.DisableFeedback, expected.DisableFeedback)
	assert.Equal(t, actual.DisabledAt, expected.DisabledAt)
	assert.Equal(t, actual.ErrorKey, expected.ErrorKey)
	assert.Equal(t, actual.Module, expected.Module)
	assert.Equal(t, actual.UserVote, expected.UserVote)
	if !createdAt.IsZero() {
		assert.Equal(t, actual.CreatedAt, types.Timestamp(createdAt.UTC().Format(time.RFC3339)))
	}
}

// SortReports sorts a list of RuleOnReport by ErrorKey field
func SortReports(reports []types.RuleOnReport) []types.RuleOnReport {
	errorKeyReport := make(map[string]types.RuleOnReport)
	var errorKeys []string
	var sorted []types.RuleOnReport
	for _, rep := range reports {
		errorKeyStr := string(rep.ErrorKey)
		errorKeyReport[errorKeyStr] = rep
		errorKeys = append(errorKeys, errorKeyStr)
	}
	sort.Strings(errorKeys)
	for _, ek := range errorKeys {
		sorted = append(sorted, errorKeyReport[ek])
	}
	return sorted
}
