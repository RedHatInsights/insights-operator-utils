package httputils_test

import (
	"testing"

	httputils "github.com/RedHatInsights/insights-operator-utils/http"
	"github.com/stretchr/testify/assert"
)

// TestNewTLSConfig tests the NewTLSConfig method
func TestNewTLSConfig(t *testing.T) {
	testCases := []struct {
		input       string
		expectedErr bool
	}{
		{"", true},
		{"aaaaa", true},
		{"../testdata/cert.pem", false},
	}

	for _, tc := range testCases {
		cfg, err := httputils.NewTLSConfig(tc.input)
		if tc.expectedErr {
			assert.Error(t, err)
			assert.Nil(t, cfg)
		} else {
			assert.NoError(t, err)
			assert.NotNil(t, cfg)
		}
	}
}
