package helpers_test

import (
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/RedHatInsights/insights-operator-utils/tests/helpers"
)

func TestMockT_WrappedMethods(t *testing.T) {
	for _, method := range []string{
		"Cleanup", "Error", "Errorf", "Fail", "FailNow", "Failed", "Fatal",
		"Fatalf", "Log", "Logf", "Skip", "SkipNow", "Skipf", "Skipped",
	} {
		t.Run(method, func(t *testing.T) {
			mockT := helpers.NewMockT(t)
			defer mockT.Finish()

			expect := mockT.Expects.EXPECT()

			switch method {
			case "Cleanup":
				expect.Cleanup(gomock.Any())
				mockT.Cleanup(func() {})
			case "Error":
				expect.Error(gomock.Any())
				mockT.Error()
			case "Errorf":
				expect.Errorf(gomock.Any())
				mockT.Errorf("")
			case "Fail":
				expect.Fail()
				mockT.Fail()
			case "FailNow":
				expect.FailNow()
				mockT.FailNow()
			case "Failed":
				expect.Failed()
				mockT.Failed()
			case "Fatal":
				expect.Fatal(gomock.Any())
				mockT.Fatal("")
			case "Fatalf":
				expect.Fatalf(gomock.Any())
				mockT.Fatalf("")
			case "Log":
				expect.Log(gomock.Any())
				mockT.Log()
			case "Logf":
				expect.Logf(gomock.Any())
				mockT.Logf("")
			case "Skip":
				expect.Skip(gomock.Any())
				mockT.Skip("")
			case "SkipNow":
				expect.SkipNow()
				mockT.SkipNow()
			case "Skipf":
				expect.Skipf(gomock.Any(), gomock.Any())
				mockT.Skipf("")
			case "Skipped":
				expect.Skipped()
				mockT.Skipped()
			}
		})
	}
}
