package tlsutil_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	tlsutil "github.com/RedHatInsights/insights-operator-utils/tls"
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
		cfg, err := tlsutil.NewTLSConfig(tc.input)
		if tc.expectedErr {
			assert.Error(t, err)
			assert.Nil(t, cfg)
		} else {
			assert.NoError(t, err)
			assert.NotNil(t, cfg)
		}
	}
}
