package clowder_test

import (
	"github.com/RedHatInsights/insights-operator-utils/clowder"
	"github.com/RedHatInsights/insights-operator-utils/kafka"
	"github.com/RedHatInsights/insights-operator-utils/postgres"
	api "github.com/redhatinsights/app-common-go/pkg/api/v1"
	"github.com/stretchr/testify/assert"
	"github.com/tisnik/go-capture"
	"testing"
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

func TestUseBrokerConfigNoAuthNoPort(t *testing.T) {
	brokerCfg := kafka.BrokerConfiguration{}
	loadedConfig := api.AppConfig{
		Kafka: &api.KafkaConfig{
			Brokers: []api.BrokerConfig{
				{
					Hostname: "test_broker",
					Port:     nil,
				},
			},
		},
	}

	clowder.UseBrokerConfig(&brokerCfg, &loadedConfig)
	assert.Equal(t, "test_broker", brokerCfg.Address)
}

func TestUseBrokerConfigNoAuth(t *testing.T) {
	brokerCfg := kafka.BrokerConfiguration{}
	port := 12345
	loadedConfig := api.AppConfig{
		Kafka: &api.KafkaConfig{
			Brokers: []api.BrokerConfig{
				{
					Hostname: "test_broker",
					Port:     &port,
				},
			},
		},
	}

	clowder.UseBrokerConfig(&brokerCfg, &loadedConfig)
	assert.Equal(t, "test_broker:12345", brokerCfg.Address)
}

func TestUseBrokerConfigAuthEnabledNoSasl(t *testing.T) {
	brokerCfg := kafka.BrokerConfiguration{}
	port := 12345
	authType := api.BrokerConfigAuthtypeSasl
	loadedConfig := api.AppConfig{
		Kafka: &api.KafkaConfig{
			Brokers: []api.BrokerConfig{
				{
					Hostname: "test_broker",
					Port:     &port,
					Authtype: &authType,
				},
			},
		},
	}

	output, _ := capture.StandardOutput(func() {
		clowder.UseBrokerConfig(&brokerCfg, &loadedConfig)
	})

	assert.Equal(t, "test_broker:12345", brokerCfg.Address)
	assert.Contains(t, output, clowder.NoSaslCfg)
}

func TestUseBrokerConfigAuthEnabledWithSaslConfig(t *testing.T) {
	brokerCfg := kafka.BrokerConfiguration{}
	port := 12345
	saslCfg := "user_pwd"
	protocol := "tls"

	authType := api.BrokerConfigAuthtypeSasl
	loadedConfig := api.AppConfig{
		Kafka: &api.KafkaConfig{
			Brokers: []api.BrokerConfig{
				{
					Hostname: "test_broker",
					Port:     &port,
					Authtype: &authType,
					Sasl: &api.KafkaSASLConfig{
						Password:      &saslCfg,
						Username:      &saslCfg,
						SaslMechanism: &saslCfg,
					},
					SecurityProtocol: &protocol,
				},
			},
		},
	}

	output, _ := capture.StandardOutput(func() {
		clowder.UseBrokerConfig(&brokerCfg, &loadedConfig)
	})

	assert.Equal(t, "test_broker:12345", brokerCfg.Address)
	assert.Contains(t, output, "kafka is configured to use authentication")
	assert.Equal(t, saslCfg, brokerCfg.SaslUsername)
	assert.Equal(t, saslCfg, brokerCfg.SaslPassword)
	assert.Equal(t, saslCfg, brokerCfg.SaslMechanism)
	assert.Equal(t, protocol, brokerCfg.SecurityProtocol)
}
