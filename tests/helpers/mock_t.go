package helpers

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

func (t *MockT) ExpectFailOnError(err error) {
	t.Expects.EXPECT().Errorf(
		gomock.Any(),
		gomock.Any(),
	)

	t.Expects.EXPECT().Fatal(err)
}

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

func (t *MockT) Error(args ...interface{}) {
	t.Expects.Error(args...)
}

func (t *MockT) Errorf(format string, args ...interface{}) {
	t.Expects.Errorf(format, args...)
}

func (t *MockT) Fail() {
	t.Expects.Fail()
}

func (t *MockT) FailNow() {
	t.Expects.FailNow()
}

func (t *MockT) Failed() bool {
	return t.Expects.Failed()
}

func (t *MockT) Fatal(args ...interface{}) {
	t.Expects.Fatal(args...)
}

func (t *MockT) Fatalf(format string, args ...interface{}) {
	t.Expects.Fatalf(format, args...)
}

func (t *MockT) Log(args ...interface{}) {
	t.Expects.Log(args...)
}

func (t *MockT) Logf(format string, args ...interface{}) {
	t.Expects.Logf(format, args...)
}

func (t *MockT) Skip(args ...interface{}) {
	t.Expects.Skip(args...)
}

func (t *MockT) SkipNow() {
	t.Expects.SkipNow()
}

func (t *MockT) Skipf(format string, args ...interface{}) {
	t.Expects.Skipf(format, args...)
}

func (t *MockT) Skipped() bool {
	return t.Expects.Skipped()
}
