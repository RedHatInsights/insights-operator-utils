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
// https://redhatinsights.github.io/insights-operator-utils/packages/parsers/date_test.html

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/RedHatInsights/insights-operator-utils/parsers"
)

// TestParseDates checks the function parsers.ParseDates for valid input.
func TestParseDates(t *testing.T) {
	dateFormat := "02/01/2006"
	startDate := "01/01/2021"
	endDate := "31/12/2021"

	expectedStartDate := time.Date(2021, time.January, 01, 00, 00, 00, 0, time.UTC)
	expectedEndDate := time.Date(2021, time.December, 31, 0, 0, 0, 0, time.UTC)

	t.Run("valid dates", func(t *testing.T) {
		gotStartDate, gotEndDate, err := parsers.ParseDates(dateFormat, startDate, endDate)

		assert.Nil(t, err, "unexpected error")
		assert.Equal(t, expectedStartDate, gotStartDate)
		assert.Equal(t, expectedEndDate, gotEndDate)
	})
}

// TestParseDatesEmptyInput check the function parsers.ParseDates for empty
// input.
func TestParseDatesEmptyInput(t *testing.T) {
	dateFormat := "02/01/2006"
	startDate := "01/01/2021"
	endDate := "31/12/2021"

	t.Run("empty start date", func(t *testing.T) {
		_, _, err := parsers.ParseDates(dateFormat, "", endDate)

		assert.Equal(t, err.Error(), "empty date")
	})

	t.Run("empty end date", func(t *testing.T) {
		_, _, err := parsers.ParseDates(dateFormat, startDate, "")

		assert.Equal(t, err.Error(), "empty date")
	})

	t.Run("empty both dates", func(t *testing.T) {
		_, _, err := parsers.ParseDates(dateFormat, startDate, "")

		assert.Equal(t, err.Error(), "empty date")
	})
}

// TestParseDatesInvalidInput check the function parsers.ParseDates for invalid
// input.
func TestParseDatesInvalidDate(t *testing.T) {
	dateFormat := "02/01/2006"
	startDate := "01/01/2021"
	endDate := "31/12/2021"

	t.Run("invalid start date", func(t *testing.T) {
		_, _, err := parsers.ParseDates(dateFormat, "invalid!!!", endDate)

		assert.Contains(t, err.Error(), "error parsing the start date")
	})

	t.Run("invalid end date", func(t *testing.T) {
		_, _, err := parsers.ParseDates(dateFormat, startDate, "invalid!!!")

		assert.Contains(t, err.Error(), "error parsing the end date")
	})

	t.Run("invalid both dates", func(t *testing.T) {
		_, _, err := parsers.ParseDates(dateFormat, "invalid!!!", "invalid as well!!!")

		assert.Contains(t, err.Error(), "error parsing the start date")
	})
}
