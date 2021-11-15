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

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
)

// metricsPusher let us mock the pushCollectors function
type metricsPusher interface {
	Push() error
	Client(c push.HTTPDoer) *push.Pusher
	Collector(c prometheus.Collector) *push.Pusher
}

// pushCollectors pushes the metrics using a metricsPusher interface
func pushCollectors(p metricsPusher, gatewayAuthToken string, collectors []prometheus.Collector) error {
	client := PushGatewayClient{
		AuthToken:  gatewayAuthToken,
		HTTPClient: http.Client{},
	}
	for _, collector := range collectors {
		p.Collector(collector)
	}
	p.Client(&client)
	return p.Push()
}
