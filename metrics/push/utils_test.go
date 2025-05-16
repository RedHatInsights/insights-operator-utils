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

// Documentation in literate-programming-style is available at:
// https://redhatinsights.github.io/insights-operator-utils/packages/metrics/push/utils_test.html

import (
	"testing"

	"github.com/RedHatInsights/insights-operator-utils/metrics/push"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
)

var (
	validPrefix = "a_valid_"
)

func TestNewCounter(t *testing.T) {
	t.Run("valid counter", func(t *testing.T) {
		_, err := push.NewCounterWithError(prometheus.CounterOpts{Name: validPrefix + "counter"})
		assert.NoError(t, err)
	})
}

func TestNewCounterVec(t *testing.T) {
	t.Run("valid counter", func(t *testing.T) {
		_, err := push.NewCounterVecWithError(prometheus.CounterOpts{Name: validPrefix + "counter_vec"}, []string{})
		assert.NoError(t, err)
	})
}

func TestNewGauge(t *testing.T) {
	t.Run("valid gauge", func(t *testing.T) {
		_, err := push.NewGaugeWithError(prometheus.GaugeOpts{Name: validPrefix + "gauge"})
		assert.NoError(t, err)
	})
}

func TestNewHistogram(t *testing.T) {
	t.Run("valid histogram", func(t *testing.T) {
		_, err := push.NewHistogramWithError(prometheus.HistogramOpts{Name: validPrefix + "histogram"})
		assert.NoError(t, err)
	})
}
