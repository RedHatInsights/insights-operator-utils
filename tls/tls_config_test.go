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
