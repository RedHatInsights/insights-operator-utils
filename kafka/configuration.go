/*
Copyright © 2020, 2023 Red Hat, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package kafka contains data types, interfaces, and methods related to
// Kafka that can be used to configure brokers, as well as consume/produce
// messages.
package kafka

import (
	"crypto/sha512"
	"strings"
	"time"

	tlsutils "github.com/RedHatInsights/insights-operator-utils/tls"
	"github.com/Shopify/sarama"
	"github.com/rs/zerolog/log"
)

// SASLConfiguration represents configuration of SASL authentication for
// a given Kafka broker
type SASLConfiguration struct {
	SecurityProtocol string `mapstructure:"security_protocol" toml:"security_protocol"`
	CertPath         string `mapstructure:"cert_path" toml:"cert_path"`
	SaslMechanism    string `mapstructure:"sasl_mechanism" toml:"sasl_mechanism"`
	SaslUsername     string `mapstructure:"sasl_username" toml:"sasl_username"`
	SaslPassword     string `mapstructure:"sasl_password" toml:"sasl_password"`
}

// SingleBrokerConfiguration represents configuration of a single-instance Kafka broker
type SingleBrokerConfiguration struct {
	Address          string        `mapstructure:"address" toml:"address"`
	SecurityProtocol string        `mapstructure:"security_protocol" toml:"security_protocol"`
	CertPath         string        `mapstructure:"cert_path" toml:"cert_path"`
	SaslMechanism    string        `mapstructure:"sasl_mechanism" toml:"sasl_mechanism"`
	SaslUsername     string        `mapstructure:"sasl_username" toml:"sasl_username"`
	SaslPassword     string        `mapstructure:"sasl_password" toml:"sasl_password"`
	Topic            string        `mapstructure:"topic" toml:"topic"`
	Timeout          time.Duration `mapstructure:"timeout" toml:"timeout"`
	Group            string        `mapstructure:"group" toml:"group"`
	ClientID         string        `mapstructure:"client_id" toml:"client_id"`
	Enabled          bool          `mapstructure:"enabled" toml:"enabled"`
}

// MultiBrokerConfiguration represents configuration of Kafka broker with
// multiple instances running on different hosts
type MultiBrokerConfiguration struct {
	Addresses        []string            `mapstructure:"addresses" toml:"addresses"`
	SecurityProtocol string              `mapstructure:"security_protocol" toml:"security_protocol"`
	SASLConfigs      []SASLConfiguration `mapstructure:"sasl_configs" toml:"sasl_configs"`
	Topic            string              `mapstructure:"topic" toml:"topic"`
	Timeout          time.Duration       `mapstructure:"timeout" toml:"timeout"`
	Group            string              `mapstructure:"group" toml:"group"`
	ClientID         string              `mapstructure:"client_id" toml:"client_id"`
	Enabled          bool                `mapstructure:"enabled" toml:"enabled"`
}

// SaramaConfigFromBrokerConfig returns a Config struct from broker.Configuration parameters
func SaramaConfigFromBrokerConfig(cfg *SingleBrokerConfiguration) (*sarama.Config, error) {
	saramaConfig := sarama.NewConfig()
	saramaConfig.Version = sarama.V0_10_2_0

	if cfg.Timeout > 0 {
		saramaConfig.Net.DialTimeout = cfg.Timeout
		saramaConfig.Net.ReadTimeout = cfg.Timeout
		saramaConfig.Net.WriteTimeout = cfg.Timeout
	}

	if strings.Contains(cfg.SecurityProtocol, "SSL") {
		saramaConfig.Net.TLS.Enable = true
	}

	if strings.EqualFold(cfg.SecurityProtocol, "SSL") && cfg.CertPath != "" {
		tlsConfig, err := tlsutils.NewTLSConfig(cfg.CertPath)
		if err != nil {
			log.Error().Msgf("Unable to load TLS config for %s cert", cfg.CertPath)
			return nil, err
		}
		saramaConfig.Net.TLS.Config = tlsConfig
	} else if strings.HasPrefix(cfg.SecurityProtocol, "SASL_") {
		log.Info().Msg("Configuring SASL authentication")
		saramaConfig.Net.SASL.Enable = true
		saramaConfig.Net.SASL.User = cfg.SaslUsername
		saramaConfig.Net.SASL.Password = cfg.SaslPassword
		saramaConfig.Net.SASL.Mechanism = sarama.SASLMechanism(cfg.SaslMechanism)

		if strings.EqualFold(cfg.SaslMechanism, sarama.SASLTypeSCRAMSHA512) {
			log.Info().Msg("Configuring SCRAM-SHA512")
			saramaConfig.Net.SASL.Handshake = true
			saramaConfig.Net.SASL.SCRAMClientGeneratorFunc = func() sarama.SCRAMClient {
				return &SCRAMClient{HashGeneratorFcn: sha512.New}
			}
		}
	}

	// ClientID is fully optional, but by setting it, we can get rid of some warning messages in logs
	if cfg.ClientID != "" {
		// if not set, the "sarama" will be used instead
		saramaConfig.ClientID = cfg.ClientID
	}

	// now the config structure is filled-in
	return saramaConfig, nil
}
