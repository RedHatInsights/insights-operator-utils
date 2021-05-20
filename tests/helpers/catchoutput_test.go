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

package helpers_test

// Documentation in literate-programming-style is available at:
// https://redhatinsights.github.io/insights-operator-utils/packages/tests/helpers/catchoutput_test.html

import (
	"fmt"
	"os"
	"testing"

	"github.com/RedHatInsights/insights-operator-utils/tests/helpers"
	"github.com/stretchr/testify/assert"
)

func TestCatchingOutput(t *testing.T) {
	originalStdout := os.Stdout
	originalStderr := os.Stderr
	helpers.CatchingOutputs(t, func() {
		fmt.Println("Test")
	})

	assert.Equal(t, originalStdout, os.Stdout)
	assert.Equal(t, originalStderr, os.Stderr)
}
