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
	brokerCfg := kafka.BrokersConfig{
		{
			Topic: originalTopicName,
		},
	}
	kafkaTopics := map[string]api.TopicConfig{
		originalTopicName: {
			Name: clowderTopicName,
		},
		"topic2": {
			Name: "AnotherTopicName",
		},
	}

	clowder.UseClowderTopics(brokerCfg, kafkaTopics)
	assert.Equal(t, clowderTopicName, brokerCfg[0].Topic, "Clowder topic name was not used")
}

func TestUseClowderTopicsTopicFoundMultiBrokers(t *testing.T) {
	originalTopicName := "topic1"
	clowderTopicName := "NewTopicName"
	brokerCfg := kafka.BrokersConfig{
		{
			Topic: originalTopicName,
		},
	}
	kafkaTopics := map[string]api.TopicConfig{
		originalTopicName: {
			Name: clowderTopicName,
		},
		"topic2": {
			Name: "AnotherTopicName",
		},
	}

	clowder.UseClowderTopics(brokerCfg, kafkaTopics)
	assert.Equal(t, clowderTopicName, brokerCfg[0].Topic, "Clowder topic name was not used")
}

func TestUseClowderTopicsTopicNotFound(t *testing.T) {
	originalTopicName := "topic1"

	brokerCfg := kafka.BrokersConfig{
		{
			Topic: originalTopicName,
		},
	}
	kafkaTopics := map[string]api.TopicConfig{
		"topic2": {
			Name: "AnotherTopicName",
		},
	}

	output, _ := capture.StandardOutput(func() {
		clowder.UseClowderTopics(brokerCfg, kafkaTopics)
	})
	assert.Equal(t, originalTopicName, brokerCfg[0].Topic, "topic name should not change")
	assert.Contains(t, output, "warning: no kafka mapping found for topic topic1")
}

func TestGetBrokersAddressesNoBrokerConfig(t *testing.T) {
	cfg := kafka.BrokersConfig{}
	assert.Equal(t, []string{}, kafka.GetBrokersAddresses(cfg))
}

func TestGetBrokersAddressesSingleBrokerConfig(t *testing.T) {
	const addr = "some_addr"
	cfg := kafka.BrokersConfig{
		{Address: addr},
	}
	assert.Equal(t, []string{addr}, kafka.GetBrokersAddresses(cfg))
}

func TestGetBrokersAddressesMultipleBrokerConfig(t *testing.T) {
	const addr, addr2 = "some_addr", "some_other_addr"
	cfg := kafka.BrokersConfig{
		{Address: addr},
		{Address: addr2},
	}
	assert.Equal(t, []string{addr, addr2}, kafka.GetBrokersAddresses(cfg))
}

func TestUseBrokerConfigNoClowderKafkaConfig(t *testing.T) {
	brokerCfg := kafka.BrokersConfig{{}}
	loadedConfig := api.AppConfig{}

	output, _ := capture.StandardOutput(func() {
		clowder.UseBrokerConfig(brokerCfg, &loadedConfig)
	})
	assert.Contains(t, output, clowder.NoBrokerCfg)
}

func TestUseBrokerConfigNoOriginalKafkaBrokers(t *testing.T) {
	brokerCfg := kafka.BrokersConfig{}
	loadedConfig := api.AppConfig{
		Kafka: &api.KafkaConfig{},
	}

	output, _ := capture.StandardOutput(func() {
		clowder.UseBrokerConfig(brokerCfg, &loadedConfig)
	})
	assert.Contains(t, output, clowder.NoOriginalBroker)
}

func TestUseBrokerConfigNoClowderKafkaBrokers(t *testing.T) {
	brokerCfg := kafka.BrokersConfig{{}}
	loadedConfig := api.AppConfig{
		Kafka: &api.KafkaConfig{},
	}

	output, _ := capture.StandardOutput(func() {
		clowder.UseBrokerConfig(brokerCfg, &loadedConfig)
	})
	assert.Contains(t, output, clowder.NoBrokerCfg)
}

func TestUseBrokerConfigMultipleClowderKafkaBrokers(t *testing.T) {
	addr1 := "test_broker"
	addr2 := "test_broker_backup"
	port := 12345
	brokerCfg := kafka.BrokersConfig{{}}
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

	brokerCfg = clowder.UseBrokerConfig(brokerCfg, &loadedConfig)
	assert.Equal(t, fmt.Sprintf("%s:%d", addr1, port), brokerCfg[0].Address)
	assert.Equal(t, addr2, brokerCfg[1].Address)
}

func TestUseBrokerConfigNoAuthNoPort(t *testing.T) {
	addr := "test_broker"
	brokerCfg := kafka.BrokersConfig{{}}
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

	brokerCfg = clowder.UseBrokerConfig(brokerCfg, &loadedConfig)
	assert.Equal(t, addr, brokerCfg[0].Address)
}

func TestUseBrokerConfigNoAuth(t *testing.T) {
	brokerCfg := kafka.BrokersConfig{{}}
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

	brokerCfg = clowder.UseBrokerConfig(brokerCfg, &loadedConfig)
	assert.Equal(t, fmt.Sprintf("%s:%d", addr, port), brokerCfg[0].Address)
}

func TestUseBrokerConfigAuthEnabledNoSasl(t *testing.T) {
	brokerCfg := kafka.BrokersConfig{{}}
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
		brokerCfg = clowder.UseBrokerConfig(brokerCfg, &loadedConfig)
	})

	assert.Equal(t, fmt.Sprintf("%s:%d", addr, port), brokerCfg[0].Address)
	assert.Contains(t, output, clowder.NoSaslCfg)
}

func TestUseBrokerConfigAuthEnabledWithSaslConfig(t *testing.T) {
	brokerCfg := kafka.BrokersConfig{{}}
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
		brokerCfg = clowder.UseBrokerConfig(brokerCfg, &loadedConfig)
	})

	assert.Contains(t, output, "kafka is configured to use authentication")

	assert.Equal(t, fmt.Sprintf("%s:%d", addr, port), brokerCfg[0].Address)
	assert.Equal(t, saslUsr, brokerCfg[0].SaslUsername)
	assert.Equal(t, saslPwd, brokerCfg[0].SaslPassword)
	assert.Equal(t, saslMechanism, brokerCfg[0].SaslMechanism)
	assert.Equal(t, protocol, brokerCfg[0].SecurityProtocol)

	assert.Equal(t, fmt.Sprintf("%s:%d", addr2, port), brokerCfg[1].Address)
	assert.Equal(t, saslUsr2, brokerCfg[1].SaslUsername)
	assert.Equal(t, saslPwd, brokerCfg[1].SaslPassword)
	assert.Equal(t, saslMechanism, brokerCfg[1].SaslMechanism)
	assert.Equal(t, protocol, brokerCfg[1].SecurityProtocol)
}
