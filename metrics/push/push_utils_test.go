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

package push

// Documentation in literate-programming-style is available at:
// https://redhatinsights.github.io/insights-operator-utils/packages/metrics/push/push_utils_test.html

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
	"github.com/stretchr/testify/assert"
)

var testGatewayAuthToken = "auth token"

type mockPusher struct {
	callsClient    int
	callsCollector int
}

func (p *mockPusher) Push() error {
	return nil
}
func (p *mockPusher) Client(c push.HTTPDoer) *push.Pusher {
	p.callsClient++
	return push.New("", "")
}
func (p *mockPusher) Collector(c prometheus.Collector) *push.Pusher {
	p.callsCollector++
	return push.New("", "")
}

var testInitFunctions = []func() (prometheus.Collector, error){
	func() (prometheus.Collector, error) {
		metric2, err := NewCounterWithError(prometheus.CounterOpts{
			Name: "a_metric",
			Help: "a metric",
		})
		return metric2, err
	},
}

func TestPushCollectors(t *testing.T) {
	err := InitMetrics(testInitFunctions)
	assert.NoError(t, err)

	pusher := mockPusher{}

	err = pushCollectors(&pusher, testGatewayAuthToken, collectors)
	assert.NoError(t, err)
	assert.Equal(t, 1, pusher.callsClient)
	assert.Equal(t, len(collectors), pusher.callsCollector)

	err = UnregisterMetrics()
	assert.NoError(t, err)
}
