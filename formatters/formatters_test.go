// Copyright 2022 Red Hat, Inc
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

// Package formatters contains various text formatters utility functions.
package formatters_test

// Documentation in literate-programming-style is available at:
// https://redhatinsights.github.io/insights-operator-utils/packages/formatters/formatters.html

import (
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/RedHatInsights/insights-operator-utils/formatters"
)

func TestFormatNullTimeForNullInput(t *testing.T) {
	var input sql.NullTime
	formatted := formatters.FormatNullTime(input)

	assert.Empty(t, formatted)
}

func TestFormatNullTimeForNonNullInput(t *testing.T) {
	now := time.Now()

	input := sql.NullTime{
		Time:  now,
		Valid: true,
	}
	formatted := formatters.FormatNullTime(input)

	assert.NotEmpty(t, formatted)

	expected := now.Format(time.RFC3339)
	assert.Equal(t, expected, formatted)
}
