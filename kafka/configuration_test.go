// Copyright 2022, 2023 Red Hat, Inc
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

package kafka_test

import (
	"testing"
	"time"

	"github.com/RedHatInsights/insights-operator-utils/kafka"
	"github.com/RedHatInsights/insights-operator-utils/tests/helpers"

	"github.com/IBM/sarama"
	"github.com/stretchr/testify/assert"
)

func TestSaramaConfigFromBrokerConfig(t *testing.T) {
	cfg := kafka.BrokerConfiguration{}
	_, err := kafka.SaramaConfigFromBrokerConfig(&cfg)
	helpers.FailOnError(t, err)

	cfg = kafka.BrokerConfiguration{
		Timeout: time.Second,
	}
	saramaConfig, err := kafka.SaramaConfigFromBrokerConfig(&cfg)
	helpers.FailOnError(t, err)
	assert.Equal(t, time.Second, saramaConfig.Net.DialTimeout)
	assert.Equal(t, time.Second, saramaConfig.Net.ReadTimeout)
	assert.Equal(t, time.Second, saramaConfig.Net.WriteTimeout)
	assert.Equal(t, "sarama", saramaConfig.ClientID) // default value

	cfg = kafka.BrokerConfiguration{
		SecurityProtocol: "SSL",
	}

	saramaConfig, err = kafka.SaramaConfigFromBrokerConfig(&cfg)
	helpers.FailOnError(t, err)
	assert.True(t, saramaConfig.Net.TLS.Enable)

	cfg = kafka.BrokerConfiguration{
		SecurityProtocol: "SASL_SSL",
		SaslMechanism:    "PLAIN",
		SaslUsername:     "username",
		SaslPassword:     "password",
		ClientID:         "foobarbaz",
	}
	saramaConfig, err = kafka.SaramaConfigFromBrokerConfig(&cfg)
	helpers.FailOnError(t, err)
	assert.True(t, saramaConfig.Net.TLS.Enable)
	assert.True(t, saramaConfig.Net.SASL.Enable)
	assert.Equal(t, sarama.SASLMechanism("PLAIN"), saramaConfig.Net.SASL.Mechanism)
	assert.Equal(t, "username", saramaConfig.Net.SASL.User)
	assert.Equal(t, "password", saramaConfig.Net.SASL.Password)
	assert.Equal(t, "foobarbaz", saramaConfig.ClientID)

	cfg.SaslMechanism = "SCRAM-SHA-512"
	saramaConfig, err = kafka.SaramaConfigFromBrokerConfig(&cfg)
	helpers.FailOnError(t, err)
	assert.True(t, saramaConfig.Net.TLS.Enable)
	assert.True(t, saramaConfig.Net.SASL.Enable)
	assert.Equal(t, sarama.SASLMechanism(sarama.SASLTypeSCRAMSHA512), saramaConfig.Net.SASL.Mechanism)
	assert.Equal(t, "username", saramaConfig.Net.SASL.User)
	assert.Equal(t, "password", saramaConfig.Net.SASL.Password)
	assert.Equal(t, "foobarbaz", saramaConfig.ClientID)
}

func TestBadConfiguration(t *testing.T) {
	cfg := kafka.BrokerConfiguration{
		SecurityProtocol: "SSL",
		CertPath:         "missing_path",
	}

	saramaCfg, err := kafka.SaramaConfigFromBrokerConfig(&cfg)
	assert.Error(t, err)
	assert.Nil(t, saramaCfg)
}
