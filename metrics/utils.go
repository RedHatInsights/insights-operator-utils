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

package metrics

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/rs/zerolog/log"
)

const (
	logName = "Name"
	logHelp = "Help"
)

// NewCounterWithError run promauto.NewCounter() catching the panic and returning an error
func NewCounterWithError(opts prometheus.CounterOpts) (counter prometheus.Counter, err error) {
	logStartRegistering(opts.Name, opts.Help)
	defer func() { // catch the possible panic
		if panicErr := recover(); panicErr != nil {
			err = fmt.Errorf("got error while registering the counter, %v", panicErr)
			logErrorRegistering(opts.Name, opts.Help, err)
		}
	}()
	counter = promauto.NewCounter(opts)
	logFinishRegistering(opts.Name, opts.Help)
	return
}

// NewGaugeWithError run promauto.NewGauge() catching the panic and returning an error
func NewGaugeWithError(opts prometheus.GaugeOpts) (gauge prometheus.Gauge, err error) {
	logStartRegistering(opts.Name, opts.Help)
	defer func() { // catch the possible panic
		if panicErr := recover(); panicErr != nil {
			err = fmt.Errorf("got error while registering the gauge, %v", panicErr)
			logErrorRegistering(opts.Name, opts.Help, err)
		}
	}()
	gauge = promauto.NewGauge(opts)
	logFinishRegistering(opts.Name, opts.Help)
	return
}

// NewHistogramWithError run promauto.NewHistogram() catching the panic and returning an error
func NewHistogramWithError(opts prometheus.HistogramOpts) (histogram prometheus.Histogram, err error) {
	logStartRegistering(opts.Name, opts.Help)
	defer func() { // catch the possible panic
		if panicErr := recover(); panicErr != nil {
			err = fmt.Errorf("got error while registering the histogram, %v", panicErr)
			logErrorRegistering(opts.Name, opts.Help, err)
		}
	}()
	histogram = promauto.NewHistogram(opts)
	logFinishRegistering(opts.Name, opts.Help)
	return
}

func logStartRegistering(name, help string) {
	log.Debug().
		Str(logHelp, help).
		Str(logName, name).
		Msg("Registering metric")
}

func logErrorRegistering(name, help string, err error) {
	log.Error().
		Err(err).
		Str(logHelp, help).
		Str(logName, name).
		Msg("Error registering metric")
}

func logFinishRegistering(name, help string) {
	log.Debug().
		Str(logHelp, help).
		Str(logName, name).
		Msg("Metric registered")
}
