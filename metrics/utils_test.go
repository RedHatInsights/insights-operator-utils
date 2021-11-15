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

package metrics_test

import (
	"testing"

	"github.com/RedHatInsights/insights-operator-utils/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
)

var (
	invalidName = "invalid name"
	validPrefix = "a_valid_"
)

func TestNewCounterWithError(t *testing.T) {
	t.Run("invalid counter", func(t *testing.T) {

		_, err := metrics.NewCounterWithError(prometheus.CounterOpts{Name: invalidName})
		assert.Error(t, err)
	})

	t.Run("valid counter", func(t *testing.T) {
		_, err := metrics.NewCounterWithError(prometheus.CounterOpts{Name: validPrefix + "counter"})
		assert.NoError(t, err)
	})
}

func TestNewGaugeWithError(t *testing.T) {
	t.Run("invalid gauge", func(t *testing.T) {

		_, err := metrics.NewGaugeWithError(prometheus.GaugeOpts{Name: invalidName})
		assert.Error(t, err)
	})

	t.Run("valid gauge", func(t *testing.T) {
		_, err := metrics.NewGaugeWithError(prometheus.GaugeOpts{Name: validPrefix + "gauge"})
		assert.NoError(t, err)
	})
}

func TestNewHistogramWithError(t *testing.T) {
	t.Run("invalid histogram", func(t *testing.T) {

		_, err := metrics.NewHistogramWithError(prometheus.HistogramOpts{Name: invalidName})
		assert.Error(t, err)
	})

	t.Run("valid histogram", func(t *testing.T) {
		_, err := metrics.NewHistogramWithError(prometheus.HistogramOpts{Name: validPrefix + "histogram"})
		assert.NoError(t, err)
	})
}
