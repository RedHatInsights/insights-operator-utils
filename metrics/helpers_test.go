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
	"github.com/RedHatInsights/insights-operator-utils/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

func getMetric1() (prometheus.Collector, error) {
	metric1, err := metrics.NewGaugeWithError(prometheus.GaugeOpts{
		Name: "metric_1",
		Help: "the first metric",
	})
	return metric1, err
}

func getMetric2() (prometheus.Collector, error) {
	metric2, err := metrics.NewCounterWithError(prometheus.CounterOpts{
		Name: "metric_2",
		Help: "the second metric",
	})
	return metric2, err
}

// testInitFunctions let us generate some metrics for tests
var testInitFunctions = []func() (prometheus.Collector, error){
	getMetric1,
	getMetric2,
}
