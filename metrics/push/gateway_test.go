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

package push_test

import (
	"net/http"
	"testing"

	"github.com/RedHatInsights/insights-operator-utils/metrics/push"
	"github.com/stretchr/testify/assert"
)

func TestPushGatewayClientDo(t *testing.T) {
	pgc := push.PushGatewayClient{
		AuthToken:  "",
		HTTPClient: http.Client{},
	}
	t.Run("without auth token", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "https://redhat.com", nil)
		_, err := pgc.Do(req)
		assert.NoError(t, err)
	})
	pgc.AuthToken = "random token"
	t.Run("with auth token", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "https://redhat.com", nil)
		_, err := pgc.Do(req)
		assert.NoError(t, err)
	})
}
