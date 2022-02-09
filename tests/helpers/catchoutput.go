// Copyright 2021, 2022 Red Hat, Inc
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
// https://redhatinsights.github.io/insights-operator-utils/packages/tests/helpers/catchoutput.html

import (
	"io"
	"os"
	"testing"
)

// CatchingOutputs execute a function capturing and returning Stdout and Stderr for later checks
func CatchingOutputs(t *testing.T, f func()) (string, string) {
	originalStdout := os.Stdout
	originalStderr := os.Stderr
	defer func() {
		os.Stdout = originalStdout
		os.Stderr = originalStderr
	}()

	stdoutReader, fakeStdout, err := os.Pipe()
	FailOnError(t, err)

	stderrReader, fakeStderr, err := os.Pipe()
	FailOnError(t, err)

	os.Stdout = fakeStdout
	os.Stderr = fakeStderr
	f()

	err = fakeStdout.Close()
	FailOnError(t, err)

	err = fakeStderr.Close()
	FailOnError(t, err)

	stdoutOutput, err := io.ReadAll(stdoutReader)
	FailOnError(t, err)

	stderrOutput, err := io.ReadAll(stderrReader)
	FailOnError(t, err)

	return string(stdoutOutput), string(stderrOutput)
}
