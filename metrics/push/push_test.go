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
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/RedHatInsights/insights-operator-utils/metrics/push"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	testJob              = "testjob"
	testGatewayAuthToken = "gateway_auth_token"
)

func TestInitMetrics(t *testing.T) {
	zerolog.SetGlobalLevel(zerolog.DebugLevel) // TODO
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	t.Run("optimal setup", func(t *testing.T) {
		err := push.InitMetrics(testInitFunctions)
		assert.NoError(t, err)
	})
	t.Run("already registered metrics", func(t *testing.T) {
		err := push.InitMetrics(testInitFunctions)
		assert.NoError(t, err)
	})
}

func TestPushMetrics(t *testing.T) {
	err := push.InitMetrics(testInitFunctions)
	require.NoError(t, err)
	t.Run("ok response", func(t *testing.T) {
		// Fake a Pushgateway that responds with 202 to DELETE and with 200 in
		// all other cases.
		pgwOK := httptest.NewServer(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", `text/plain; charset=utf-8`)
				if r.Method == http.MethodDelete {
					w.WriteHeader(http.StatusAccepted)
					return
				}
				w.WriteHeader(http.StatusOK)
			}),
		)
		defer pgwOK.Close()

		err := push.SendMetrics(testJob, pgwOK.URL, testGatewayAuthToken)
		assert.NoError(t, err)
	})

	t.Run("error response", func(t *testing.T) {
		// Fake a Pushgateway that responds with 500.
		pgwErr := httptest.NewServer(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", `text/plain; charset=utf-8`)
				w.WriteHeader(http.StatusInternalServerError)
			}),
		)
		defer pgwErr.Close()

		err := push.SendMetrics(testJob, pgwErr.URL, testGatewayAuthToken)
		assert.Error(t, err)
	})
}

func TestPushMetricsInLoop(t *testing.T) {
	// Fake a Pushgateway that responds with 202 to DELETE and with 200 in
	// all other cases and counts the number of pushes received
	var (
		pushes          int
		expectedPushes  = 3
		timeBetweenPush = 1 * time.Second // s
		totalTime       = 4 * time.Second // give enough time
	)

	pgwOK := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", `text/plain; charset=utf-8`)
			if r.Method == http.MethodDelete {
				w.WriteHeader(http.StatusAccepted)
				return
			}
			w.WriteHeader(http.StatusOK)
			pushes++
		}),
	)
	defer pgwOK.Close()

	t.Run("if TimeBetweenPush != 0", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), totalTime)

		go push.SendMetricsInLoop(ctx, testJob, pgwOK.URL, testGatewayAuthToken, time.Duration(timeBetweenPush))
		time.Sleep(totalTime)
		cancel()
		time.Sleep(1 * time.Second) // give time for the push to complete

		assert.GreaterOrEqual(t, pushes, expectedPushes, fmt.Sprintf("expected more than %d pushes but found %d", expectedPushes, pushes))

		log.Info().Int("pushes", pushes).Msg("debug")
	})

	t.Run("if TimeBetweenPush is 0, don't do anything", func(t *testing.T) {
		lastPushes := pushes
		ctx, cancel := context.WithTimeout(context.Background(), totalTime)

		timeBetweenPush = 0 * time.Second
		go push.SendMetricsInLoop(ctx, testJob, pgwOK.URL, testGatewayAuthToken, time.Duration(timeBetweenPush))
		time.Sleep(totalTime)
		cancel()

		assert.Equal(t, lastPushes, pushes, "expected not to have pushed any more metrics")
		log.Info().Int("pushes", pushes).Msg("debug")
	})
}
