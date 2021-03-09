/*
Copyright © 2020 Red Hat, Inc.

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

package logger

import "github.com/rs/zerolog"

// LoggingConfiguration represents configuration for logging in general
type LoggingConfiguration struct {
	// Debug enables pretty colored logging
	Debug     bool `mapstructure:"debug" toml:"debug"`
	UseStderr bool `mapstructure:"use_stderr" toml:"use_stderr"`

	// LogLevel sets logging level to show. Possible values are:
	// "debug"
	// "info"
	// "warn", "warning"
	// "error"
	// "fatal"
	//
	// logging level won't be changed if value is not one of listed above
	LogLevel string `mapstructure:"log_level" toml:"log_level"`

	// LoggingToCloudWatchEnabled enables logging to CloudWatch
	// (configuration for CloudWatch is in CloudWatchConfiguration)
	LoggingToCloudWatchEnabled bool `mapstructure:"logging_to_cloud_watch_enabled" toml:"logging_to_cloud_watch_enabled"`

	// LoggingToSentryEnabled enables logging to Sentry
	// (configuration for Sentry is in SentryLoggingConfiguration)
	LoggingToSentryEnabled bool `mapstructure:"logging_to_sentry_enabled" toml:"logging_to_sentry_enabled"`

	// LoggingToKafkaEnabled enables logging to Kafka
	// (configuration for Kafka logging is in KafkaZerologConfiguration)
	LoggingToKafkaEnabled bool `mapstructure:"logging_to_kafka_enabled" toml:"logging_to_kafka_enabled"`
}

// CloudWatchConfiguration represents configuration of CloudWatch logger
type CloudWatchConfiguration struct {
	AWSAccessID             string `mapstructure:"aws_access_id" toml:"aws_access_id"`
	AWSSecretKey            string `mapstructure:"aws_secret_key" toml:"aws_secret_key"`
	AWSSessionToken         string `mapstructure:"aws_session_token" toml:"aws_session_token"`
	AWSRegion               string `mapstructure:"aws_region" toml:"aws_region"`
	LogGroup                string `mapstructure:"log_group" toml:"log_group"`
	StreamName              string `mapstructure:"stream_name" toml:"stream_name"`
	CreateStreamIfNotExists bool   `mapstructure:"create_stream_if_not_exists" toml:"create_stream_if_not_exists"`

	// enable debug logs for debugging aws client itself
	Debug bool `mapstructure:"debug" toml:"debug"`
}

// SentryLoggingConfiguration represents the configuration of Sentry logger
type SentryLoggingConfiguration struct {
	SentryDSN string `mapstructure:"dsn" toml:"dsn"`
}

// KafkaZerologConfiguration represetns the configuration for sending log messages to a Kafka topic
type KafkaZerologConfiguration struct {
	Broker   string        `mapstructure:"broker" toml:"broker"`
	Topic    string        `mapstructure:"topic" toml:"topic"`
	CertPath string        `mapstructure:"cert_path" toml:"cert_path"`
	Level    zerolog.Level `mapstructure:"level" toml:"level"`
}
