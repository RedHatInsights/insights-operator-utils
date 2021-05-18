package helpers_test

// Documentation in literate-programming-style is available at:
// https://redhatinsights.github.io/insights-operator-utils/packages/tests/helpers/mock_t_test.html

import (
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/RedHatInsights/insights-operator-utils/tests/helpers"
)

func TestMockT_WrappedMethods(t *testing.T) {
	mockT := helpers.NewMockT(t)
	defer mockT.Finish()

	expect := mockT.Expects.EXPECT()

	expect.Cleanup(gomock.Any())
	mockT.Cleanup(func() {})

	expect.Error(gomock.Any())
	mockT.Error()

	expect.Fail()
	mockT.Fail()

	expect.FailNow()
	mockT.FailNow()

	expect.Failed()
	mockT.Failed()

	expect.Fatal(gomock.Any())
	mockT.Fatal("")

	expect.Log(gomock.Any())
	mockT.Log()

	expect.Skip(gomock.Any())
	mockT.Skip("")

	expect.SkipNow()
	mockT.SkipNow()

	expect.Skipped()
	mockT.Skipped()
}

// cuz linters are crazy
func TestMockT_WrappedFMethods(t *testing.T) {
	mockT := helpers.NewMockT(t)
	defer mockT.Finish()

	expect := mockT.Expects.EXPECT()

	expect.Errorf(gomock.Any())
	mockT.Errorf("")

	expect.Fatalf(gomock.Any())
	mockT.Fatalf("")

	expect.Logf(gomock.Any())
	mockT.Logf("")

	expect.Skipf(gomock.Any(), gomock.Any())
	mockT.Skipf("")
}
