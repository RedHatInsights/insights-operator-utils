// Copyright 2024 Red Hat, Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package clowder_test

import (
	"fmt"
	"testing"

	"github.com/RedHatInsights/insights-operator-utils/clowder"
	"github.com/RedHatInsights/insights-operator-utils/kafka"
	"github.com/RedHatInsights/insights-operator-utils/postgres"
	api "github.com/redhatinsights/app-common-go/pkg/api/v1"
	"github.com/stretchr/testify/assert"
	"github.com/tisnik/go-capture"
)

func TestUseDBConfig(t *testing.T) {
	cfg := postgres.StorageConfiguration{}
	expected := postgres.StorageConfiguration{
		PGUsername: "username",
		PGPassword: "password",
		PGHost:     "hostname",
		PGPort:     1234,
		PGDBName:   "dbname",
	}
	loadedCfg := api.AppConfig{
		Database: &api.DatabaseConfig{
			AdminPassword: "adminpw",
			AdminUsername: "admin",
			Hostname:      "hostname",
			Name:          "dbname",
			Password:      "password",
			Port:          1234,
			RdsCa:         nil,
			Username:      "username",
		},
	}

	clowder.UseDBConfig(&cfg, &loadedCfg)
	assert.Equal(t, expected, cfg, "Clowder database config was not used")
}

func TestUseClowderTopicsTopicFound(t *testing.T) {
	originalTopicName := "topic1"
	clowderTopicName := "NewTopicName"
	brokerCfg := kafka.BrokerConfiguration{
		Topic: originalTopicName,
	}
	kafkaTopics := map[string]api.TopicConfig{
		originalTopicName: {
			Name: clowderTopicName,
		},
		"topic2": {
			Name: "AnotherTopicName",
		},
	}

	clowder.UseClowderTopics(&brokerCfg, kafkaTopics)
	assert.Equal(t, clowderTopicName, brokerCfg.Topic, "Clowder topic name was not used")
}

func TestUseClowderTopicsTopicFoundMultiBrokers(t *testing.T) {
	originalTopicName := "topic1"
	clowderTopicName := "NewTopicName"
	brokerCfg := kafka.BrokerConfiguration{
		Topic: originalTopicName,
	}
	kafkaTopics := map[string]api.TopicConfig{
		originalTopicName: {
			Name: clowderTopicName,
		},
		"topic2": {
			Name: "AnotherTopicName",
		},
	}

	clowder.UseClowderTopics(&brokerCfg, kafkaTopics)
	assert.Equal(t, clowderTopicName, brokerCfg.Topic, "Clowder topic name was not used")
}

func TestUseClowderTopicsTopicNotFound(t *testing.T) {
	originalTopicName := "topic1"

	brokerCfg := kafka.BrokerConfiguration{
		Topic: originalTopicName,
	}
	kafkaTopics := map[string]api.TopicConfig{
		"topic2": {
			Name: "AnotherTopicName",
		},
	}

	output, _ := capture.StandardOutput(func() {
		clowder.UseClowderTopics(&brokerCfg, kafkaTopics)
	})
	assert.Equal(t, originalTopicName, brokerCfg.Topic, "topic name should not change")
	assert.Contains(t, output, "warning: no kafka mapping found for topic topic1")
}

func TestUseBrokerConfigNoKafkaConfig(t *testing.T) {
	brokerCfg := kafka.BrokerConfiguration{}
	loadedConfig := api.AppConfig{}

	output, _ := capture.StandardOutput(func() {
		clowder.UseBrokerConfig(&brokerCfg, &loadedConfig)
	})
	assert.Contains(t, output, clowder.NoBrokerCfg)
}

func TestUseBrokerConfigNoKafkaBrokers(t *testing.T) {
	brokerCfg := kafka.BrokerConfiguration{}
	loadedConfig := api.AppConfig{
		Kafka: &api.KafkaConfig{},
	}

	output, _ := capture.StandardOutput(func() {
		clowder.UseBrokerConfig(&brokerCfg, &loadedConfig)
	})
	assert.Contains(t, output, clowder.NoBrokerCfg)
}

func TestUseBrokerConfigMultipleKafkaBrokers(t *testing.T) {
	addr1 := "test_broker"
	addr2 := "test_broker_backup"
	port := 12345
	brokerCfg := kafka.BrokerConfiguration{}
	loadedConfig := api.AppConfig{
		Kafka: &api.KafkaConfig{
			Brokers: []api.BrokerConfig{
				{
					Hostname: addr1,
					Port:     &port,
				},
				{
					Hostname: addr2,
					Port:     nil,
				},
			},
		},
	}

	clowder.UseBrokerConfig(&brokerCfg, &loadedConfig)
	assert.Equal(t, fmt.Sprintf("%s:%d,%s", addr1, port, addr2), brokerCfg.Addresses)
}

func TestUseBrokerConfigNoAuthNoPort(t *testing.T) {
	addr := "test_broker"
	brokerCfg := kafka.BrokerConfiguration{}
	loadedConfig := api.AppConfig{
		Kafka: &api.KafkaConfig{
			Brokers: []api.BrokerConfig{
				{
					Hostname: addr,
					Port:     nil,
				},
			},
		},
	}

	clowder.UseBrokerConfig(&brokerCfg, &loadedConfig)
	assert.Equal(t, addr, brokerCfg.Addresses)
}

func TestUseBrokerConfigNoAuth(t *testing.T) {
	brokerCfg := kafka.BrokerConfiguration{}
	port := 12345
	addr := "test_broker"
	loadedConfig := api.AppConfig{
		Kafka: &api.KafkaConfig{
			Brokers: []api.BrokerConfig{
				{
					Hostname: addr,
					Port:     &port,
				},
			},
		},
	}

	clowder.UseBrokerConfig(&brokerCfg, &loadedConfig)
	assert.Equal(t, fmt.Sprintf("%s:%d", addr, port), brokerCfg.Addresses)
}

func TestUseBrokerConfigAuthEnabledNoSasl(t *testing.T) {
	brokerCfg := kafka.BrokerConfiguration{}
	port := 12345
	addr := "test_broker"
	authType := api.BrokerConfigAuthtypeSasl
	loadedConfig := api.AppConfig{
		Kafka: &api.KafkaConfig{
			Brokers: []api.BrokerConfig{
				{
					Hostname: addr,
					Port:     &port,
					Authtype: &authType,
				},
			},
		},
	}

	output, _ := capture.StandardOutput(func() {
		clowder.UseBrokerConfig(&brokerCfg, &loadedConfig)
	})

	assert.Equal(t, fmt.Sprintf("%s:%d", addr, port), brokerCfg.Addresses)
	assert.Contains(t, output, clowder.NoSaslCfg)
}

func TestUseBrokerConfigAuthEnabledWithSaslConfig(t *testing.T) {
	brokerCfg := kafka.BrokerConfiguration{}
	port := 12345
	addr := "test_broker"
	addr2 := "test_broker_backup"
	saslUsr := "user"
	saslUsr2 := "user2"
	saslPwd := "pwd"
	saslMechanism := "sasl"
	protocol := "tls"

	authType := api.BrokerConfigAuthtypeSasl
	loadedConfig := api.AppConfig{
		Kafka: &api.KafkaConfig{
			Brokers: []api.BrokerConfig{
				{
					Hostname: addr,
					Port:     &port,
					Authtype: &authType,
					Sasl: &api.KafkaSASLConfig{
						Password:      &saslPwd,
						Username:      &saslUsr,
						SaslMechanism: &saslMechanism,
					},
					SecurityProtocol: &protocol,
				},
				{
					Hostname: addr2,
					Port:     &port,
					Authtype: &authType,
					Sasl: &api.KafkaSASLConfig{
						Password:      &saslPwd,
						Username:      &saslUsr2,
						SaslMechanism: &saslMechanism,
					},
					SecurityProtocol: &protocol,
				},
			},
		},
	}

	output, _ := capture.StandardOutput(func() {
		clowder.UseBrokerConfig(&brokerCfg, &loadedConfig)
	})

	assert.Equal(t, fmt.Sprintf("%s:%d,%s:%d", addr, port, addr2, port), brokerCfg.Addresses)
	assert.Contains(t, output, "kafka is configured to use authentication")
	assert.Equal(t, saslUsr, brokerCfg.SaslUsername)
	assert.Equal(t, saslPwd, brokerCfg.SaslPassword)
	assert.Equal(t, saslMechanism, brokerCfg.SaslMechanism)
	assert.Equal(t, protocol, brokerCfg.SecurityProtocol)

}
