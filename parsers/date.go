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

// Package parsers contains various text parser utility functions.
package parsers

// Documentation in literate-programming-style is available at:
// https://redhatinsights.github.io/insights-operator-utils/packages/parsers/date.html

import (
	"fmt"
	"time"
)

// ParseDates converts two strings to time.Time() format using the provided
// date format. In case of any error or if any input string is empty, error
// object is returned together with "null" dates.
func ParseDates(dateFormat, startDate, endDate string) (startDateParsed, endDateParsed time.Time, err error) {
	// check if both dates are set
	if startDate == "" || endDate == "" {
		return time.Time{}, time.Time{}, fmt.Errorf("empty date")
	}

	// try to parse first date
	startDateParsed, err = time.Parse(dateFormat, startDate)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("error parsing the start date: %v", err)
	}

	// try to parse second date
	endDateParsed, err = time.Parse(dateFormat, endDate)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("error parsing the end date: %v", err)
	}

	// conversion results
	return startDateParsed, endDateParsed, nil
}
