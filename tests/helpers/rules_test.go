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

package helpers_test

// Documentation in literate-programming-style is available at:
// https://redhatinsights.github.io/insights-operator-utils/packages/tests/helpers/mock_t_test.html

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/RedHatInsights/insights-operator-utils/tests/helpers"
	types "github.com/RedHatInsights/insights-results-types"
)

func reportWithErrorKey(errorKey string) types.RuleOnReport {
	return types.RuleOnReport{
		Module:          "",
		ErrorKey:        types.ErrorKey(errorKey),
		UserVote:        types.UserVoteNone,
		Disabled:        false,
		DisableFeedback: "",
		DisabledAt:      "",
		TemplateData:    nil,
		CreatedAt:       "",
	}
}

// TestCompareReportResponsesZeroTimestamp function checks the function
// CompareReportResponses when timestamp is set to zero
func TestCompareReportResponsesZeroTimestamp(t *testing.T) {
	// let's use zero time value there
	var timestamp time.Time

	expectedReport := reportWithErrorKey("")

	// exactly the same report
	actualReport := expectedReport

	// reports should be equal
	helpers.CompareReportResponses(t, expectedReport, actualReport, timestamp)
}

// TestCompareReportResponsesRealTimestamp function checks the function
// CompareReportResponses when timestamp is set to real timestamp
func TestCompareReportResponsesRealTimestamp(t *testing.T) {
	timestamp := time.Now()
	formatted := timestamp.UTC().Format(time.RFC3339)

	expectedReport := types.RuleOnReport{
		Module:          "",
		ErrorKey:        "",
		UserVote:        types.UserVoteNone,
		Disabled:        false,
		DisableFeedback: "",
		DisabledAt:      "",
		TemplateData:    nil,
		CreatedAt:       types.Timestamp(formatted),
	}

	// exactly the same report
	actualReport := expectedReport

	// reports should be equal
	helpers.CompareReportResponses(t, expectedReport, actualReport, timestamp)
}

// TestSortReportsEmptySlice function checks the function
// SortReports when the input slice is empty
func TestSortReportsEmptySlice(t *testing.T) {
	var reports []types.RuleOnReport
	sorted := helpers.SortReports(reports)

	assert.Nil(t, sorted)
	assert.Len(t, sorted, 0)
}

// TestSortReportsSliceWithOneItem function checks the function
// SortReports when the input slice contains just one value
func TestSortReportsSliceWithOneItem(t *testing.T) {
	var reports []types.RuleOnReport
	report := types.RuleOnReport{
		Module:          "",
		ErrorKey:        "",
		UserVote:        types.UserVoteNone,
		Disabled:        false,
		DisableFeedback: "",
		DisabledAt:      "",
		TemplateData:    nil,
		CreatedAt:       "",
	}

	// put the only report into the slice
	reports = append(reports, report)

	sorted := helpers.SortReports(reports)

	assert.NotNil(t, sorted)
	assert.Len(t, sorted, 1)

	assert.Equal(t, sorted[0], report)
}

// TestSortReportsSliceWithTwoSortedItems function checks the function
// SortReports when the input slice contains just two already sorted values
func TestSortReportsSliceWithTwoSortedItems(t *testing.T) {
	var reports []types.RuleOnReport
	report1 := reportWithErrorKey("AAA")
	report2 := reportWithErrorKey("BBB")

	// put both reports into the slice
	reports = append(reports, report1, report2)

	sorted := helpers.SortReports(reports)

	assert.NotNil(t, sorted)
	assert.Len(t, sorted, 2)

	assert.Equal(t, sorted[0], report1)
	assert.Equal(t, sorted[1], report2)
}

// TestSortReportsSliceWithTwoUnsortedItems function checks the function
// SortReports when the input slice contains just two unsorted values
func TestSortReportsSliceWithTwoUnsortedItems(t *testing.T) {
	var reports []types.RuleOnReport
	report1 := reportWithErrorKey("BBB")
	report2 := reportWithErrorKey("AAA")

	// put both reports into the slice
	reports = append(reports, report1, report2)

	sorted := helpers.SortReports(reports)

	assert.NotNil(t, sorted)
	assert.Len(t, sorted, 2)

	assert.Equal(t, sorted[0], report2)
	assert.Equal(t, sorted[1], report1)
}
