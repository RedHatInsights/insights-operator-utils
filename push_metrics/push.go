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

package push_metrics

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
	"github.com/rs/zerolog/log"
)

// collectors stores the prometheus collectors that will be pushed to the push gateway
var collectors []prometheus.Collector

// UnregisterMetrics unregister all prometheus collectors safe in the collectors variable
func UnregisterMetrics() error {
	var errMsg = "cannot unregister metric"
	log.Debug().
		Int("total", len(collectors)).
		Int("progress", 0).
		Msg("Unregistering metrics")
	for i, c := range collectors {
		if ok := prometheus.Unregister(c); !ok {
			log.Warn().Msg(errMsg)
		} else {
			log.Debug().Msg("metric unregistered")
		}

		log.Debug().
			Int("total", len(collectors)).
			Int("progress", i+1).
			Msg("Unregistering metrics")
	}
	return nil
}

// InitMetrics fills the collector variables with some Prometheus metrics and automatically registers them.
func InitMetrics(initFunctions []func() (prometheus.Collector, error)) (err error) {
	// Reset the collector slice
	if len(collectors) > 0 {
		if err = UnregisterMetrics(); err != nil {
			return err
		}
	}

	collectors = []prometheus.Collector{}
	for _, f := range initFunctions {
		if coll, err := f(); err == nil {
			collectors = append(collectors, coll)
		} else {
			return err
		}
	}
	return nil
}

// PushMetrics pushes the metrics to the configured prometheus push gateway
func PushMetrics(job, gatewayURL, gatewayAuthToken string) error {
	// Creates a pusher to the gateway "$PUSHGW_URL/metrics/job/$(job_name)
	log.Debug().
		Str("Job", job).
		Str("url", gatewayURL).
		Msg("Pushing metrics")
	pusher := push.New(gatewayURL, job)

	err := pushCollectors(pusher, gatewayAuthToken, collectors)
	if err != nil {
		log.Err(err).Msg("Couldn't push prometheus metrics")
		return err
	}
	log.Info().Msg("Metrics pushed successfully.")
	return nil
}

// PushMetricsInLoop pushes the metrics in a loop until context is done
func PushMetricsInLoop(ctx context.Context, job, gatewayURL, gatewayAuthToken string, timeBetweenPush time.Duration) {
	if timeBetweenPush < time.Second*1 {
		log.Warn().Msgf("You are trying to push the metrics every %f seconds. This may overload the push gateway, so this operation is blocked.", timeBetweenPush.Seconds())
		return
	}
	ticker := time.NewTicker(timeBetweenPush)
	for {
		select {
		case <-ticker.C:
			log.Debug().Msg("Pushing metrics")
			_ = PushMetrics(job, gatewayURL, gatewayAuthToken)
		case <-ctx.Done():
			return
		}
	}
}
