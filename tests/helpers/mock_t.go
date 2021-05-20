// Copyright 2020 Red Hat, Inc
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
// https://redhatinsights.github.io/insights-operator-utils/packages/tests/helpers/mock_t.html

import (
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/RedHatInsights/insights-operator-utils/tests/mock_testing"
)

// MockT wraps testing.T to be able to test functions accepting testing.TB.
// Don't forget to call Finish at the end of the test `defer mockT.Finish()`
type MockT struct {
	*testing.T
	Expects        *mock_testing.MockTB
	mockController *gomock.Controller
}

// NewMockT constructs a new instance of MockT
func NewMockT(t *testing.T) *MockT {
	mockController := gomock.NewController(t)

	mockTB := mock_testing.NewMockTB(mockController)

	return &MockT{
		T:              t,
		Expects:        mockTB,
		mockController: mockController,
	}
}

// Finish cleans up after the MockT
func (t *MockT) Finish() {
	defer t.mockController.Finish()
}

// ExpectFailOnError adds expects corresponding to those called by helpers.FailOnError function
func (t *MockT) ExpectFailOnError(err error) {
	t.Expects.EXPECT().Errorf(
		gomock.Any(),
		gomock.Any(),
	)

	t.Expects.EXPECT().Fatal(err)
}

// ExpectFailOnErrorAnyArgument adds expects corresponding to those called by helpers.FailOnError function
// with any argument
func (t *MockT) ExpectFailOnErrorAnyArgument() {
	t.Expects.EXPECT().Errorf(
		gomock.Any(),
		gomock.Any(),
	)

	t.Expects.EXPECT().Fatal(gomock.Any())
}

// Cleanup mocks Cleanup method of testing.T
func (t *MockT) Cleanup(f func()) {
	t.Expects.Cleanup(f)
}

// Error mocks Error method of testing.T
func (t *MockT) Error(args ...interface{}) {
	t.Expects.Error(args...)
}

// Errorf mocks Errorf method of testing.T
func (t *MockT) Errorf(format string, args ...interface{}) {
	t.Expects.Errorf(format, args...)
}

// Fail mocks Fail method of testing.T
func (t *MockT) Fail() {
	t.Expects.Fail()
}

// FailNow mocks Fail method of testing.T
func (t *MockT) FailNow() {
	t.Expects.FailNow()
}

// Failed mocks Failed method of testing.T
func (t *MockT) Failed() bool {
	return t.Expects.Failed()
}

// Fatal mocks Fatal method of testing.T
func (t *MockT) Fatal(args ...interface{}) {
	t.Expects.Fatal(args...)
}

// Fatalf mocks Fatalf method of testing.T
func (t *MockT) Fatalf(format string, args ...interface{}) {
	t.Expects.Fatalf(format, args...)
}

// Log mocks Log method of testing.T
func (t *MockT) Log(args ...interface{}) {
	t.Expects.Log(args...)
}

// Logf mocks Logf method of testing.T
func (t *MockT) Logf(format string, args ...interface{}) {
	t.Expects.Logf(format, args...)
}

// Skip mocks Skip method of testing.T
func (t *MockT) Skip(args ...interface{}) {
	t.Expects.Skip(args...)
}

// SkipNow mocks SkipNow method of testing.T
func (t *MockT) SkipNow() {
	t.Expects.SkipNow()
}

// Skipf mocks Skipf method of testing.T
func (t *MockT) Skipf(format string, args ...interface{}) {
	t.Expects.Skipf(format, args...)
}

// Skipped mocks Skipped method of testing.T
func (t *MockT) Skipped() bool {
	return t.Expects.Skipped()
}
