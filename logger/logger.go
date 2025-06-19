/*
Copyright © 2020, 2021, 2022, 2023 Red Hat, Inc.

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

// Package logger contains the configuration structures needed to configure
// the access to CloudWatch server to sending the log messages there.
package logger

// Documentation in literate-programming-style is available at:
// https://redhatinsights.github.io/insights-operator-utils/packages/logger/logger.html

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/IBM/sarama"
	"github.com/RedHatInsights/cloudwatch"
	"github.com/RedHatInsights/kafka-zerolog/kafkazerolog"
	zlogsentry "github.com/archdx/zerolog-sentry"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	sentry "github.com/getsentry/sentry-go"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var needClose []io.Closer = []io.Closer{}

// WorkaroundForRHIOPS729 keeps only those fields that are currently getting their way to Kibana
// TODO: delete when RHIOPS-729 is fixed
type WorkaroundForRHIOPS729 struct {
	io.Writer
}

func (writer WorkaroundForRHIOPS729) Write(bytes []byte) (int, error) {
	var obj map[string]interface{}

	err := json.Unmarshal(bytes, &obj)
	if err != nil {
		// it's not JSON object, so we don't do anything
		return writer.Writer.Write(bytes)
	}

	// lowercase the keys
	for key := range obj {
		val := obj[key]
		delete(obj, key)
		obj[strings.ToUpper(key)] = val
	}

	resultBytes, err := json.Marshal(obj)
	if err != nil {
		return 0, err
	}

	written, err := writer.Writer.Write(resultBytes)
	if err != nil {
		return written, err
	}

	if written < len(resultBytes) {
		return written, fmt.Errorf("too few bytes were written")
	}

	return len(bytes), nil
}

// AWSCloudWatchEndpoint allows you to mock cloudwatch client by redirecting requests to a local proxy
var AWSCloudWatchEndpoint string

// InitZerolog initializes zerolog with provided configs to use proper stdout and/or CloudWatch logging
func InitZerolog(
	loggingConf LoggingConfiguration, cloudWatchConf CloudWatchConfiguration, sentryConf SentryLoggingConfiguration,
	kafkazerologConf KafkaZerologConfiguration, additionalWriters ...io.Writer,
) error {
	setGlobalLogLevel(loggingConf)

	var writers []io.Writer
	selectedOutput := os.Stdout

	writers = append(writers, additionalWriters...)

	if loggingConf.UseStderr {
		selectedOutput = os.Stderr
	}

	var consoleWriter io.Writer
	consoleWriter = selectedOutput

	if loggingConf.Debug {
		// nice colored output
		consoleWriter = zerolog.ConsoleWriter{Out: selectedOutput}
	}

	writers = append(writers, consoleWriter)

	if loggingConf.LoggingToCloudWatchEnabled {
		cloudWatchWriter, err := setupCloudwatchLogging(&cloudWatchConf)
		if err != nil {
			err = fmt.Errorf("Error initializing Cloudwatch logging: %s", err.Error())
			return err
		}

		writers = append(writers, &WorkaroundForRHIOPS729{Writer: cloudWatchWriter})
	}

	if loggingConf.LoggingToSentryEnabled {
		sentryWriter, err := setupSentryLogging(sentryConf)
		if err != nil {
			err = fmt.Errorf("Error initializing Sentry logging: %s", err.Error())
			return err
		}
		writers = append(writers, sentryWriter)
		needClose = append(needClose, sentryWriter)
	}

	if loggingConf.LoggingToKafkaEnabled {
		kafkaWriter, err := setupKafkaZerolog(kafkazerologConf)
		if err != nil {
			err = fmt.Errorf("Error initializing Kafka logging: %s", err.Error())
			return err
		}
		writers = append(writers, kafkaWriter)
		needClose = append(needClose, kafkaWriter)
	}

	logsWriter := zerolog.MultiLevelWriter(writers...)
	log.Logger = zerolog.New(logsWriter).With().Timestamp().Logger()

	// zerolog doesn't implement Println required by sarama
	sarama.Logger = &SaramaZerologger{zerologger: log.Logger}
	return nil
}

// CloseZerolog closes properly the zerolog, if needed
func CloseZerolog() {
	for _, toClose := range needClose {
		if err := toClose.Close(); err != nil {
			log.Debug().Err(err).Msg("Error when closing")
		}
	}
}

func setGlobalLogLevel(configuration LoggingConfiguration) {
	logLevel := convertLogLevel(configuration.LogLevel)
	zerolog.SetGlobalLevel(logLevel)
}

func convertLogLevel(level string) zerolog.Level {
	level = strings.ToLower(strings.TrimSpace(level))
	switch level {
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn", "warning":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	case "fatal":
		return zerolog.FatalLevel
	}

	return zerolog.DebugLevel
}

func setupCloudwatchLogging(conf *CloudWatchConfiguration) (io.Writer, error) {
	// os.Hostname is preferred to os.Getenv("HOSTNAME") because the env var might
	// not be populated yet (e.g. GitHub runners)
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	// if no log stream name is explicitly provided, HOSTNAME is used
	if conf.StreamName == "" {
		conf.StreamName = hostname
	} else {
		// take provided log stream name and replace any $HOSTNAME placeholders with real hostname
		conf.StreamName = strings.ReplaceAll(conf.StreamName, "$HOSTNAME", hostname)
	}
	awsLogLevel := aws.LogOff
	if conf.Debug {
		awsLogLevel = aws.LogDebugWithSigning |
			aws.LogDebugWithHTTPBody |
			aws.LogDebugWithEventStreamBody
	}

	awsConf := aws.NewConfig().
		WithCredentials(credentials.NewStaticCredentials(
			conf.AWSAccessID, conf.AWSSecretKey, conf.AWSSessionToken,
		)).
		WithRegion(conf.AWSRegion).
		WithLogLevel(awsLogLevel)

	if len(AWSCloudWatchEndpoint) > 0 {
		awsConf = awsConf.WithEndpoint(AWSCloudWatchEndpoint)
	}

	cloudWatchSession := session.Must(session.NewSession(awsConf))
	CloudWatchClient := cloudwatchlogs.New(cloudWatchSession)

	var cloudWatchWriter io.Writer
	if conf.CreateStreamIfNotExists {
		group := cloudwatch.NewGroup(conf.LogGroup, CloudWatchClient)

		var err error
		cloudWatchWriter, err = group.Create(conf.StreamName)
		if err != nil {
			return nil, err
		}
	} else {
		cloudWatchWriter = cloudwatch.NewWriter(
			conf.LogGroup, conf.StreamName, CloudWatchClient,
		)
	}

	return cloudWatchWriter, nil
}

func sentryBeforeSend(event *sentry.Event, _ *sentry.EventHint) *sentry.Event {
	event.Fingerprint = []string{event.Message}
	return event
}

func setupSentryLogging(conf SentryLoggingConfiguration) (io.WriteCloser, error) {
	sentryWriter, err := zlogsentry.New(
		conf.SentryDSN,
		zlogsentry.WithEnvironment(conf.SentryEnvironment),
		zlogsentry.WithBeforeSend(sentryBeforeSend),
	)
	if err != nil {
		return nil, err
	}

	return sentryWriter, nil
}

func setupKafkaZerolog(conf KafkaZerologConfiguration) (io.WriteCloser, error) {
	return kafkazerolog.NewKafkaLogger(kafkazerolog.KafkaLoggerConf{
		Broker: conf.Broker,
		Topic:  conf.Topic,
		Cert:   conf.CertPath,
		Level:  convertLogLevel(conf.Level),
	})
}

const kafkaErrorPrefix = "kafka: error"

// SaramaZerologger is a wrapper to make sarama log to zerolog
// those logs can be filtered by key "package" with value "sarama"
type SaramaZerologger struct{ zerologger zerolog.Logger }

// Print wraps print method
func (logger *SaramaZerologger) Print(params ...interface{}) {
	var messages []string
	for _, item := range params {
		messages = append(messages, fmt.Sprint(item))
	}

	logger.logMessage("%v", strings.Join(messages, " "))
}

// Printf wraps printf method
func (logger *SaramaZerologger) Printf(format string, params ...interface{}) {
	logger.logMessage(format, params...)
}

// Println wraps println method
func (logger *SaramaZerologger) Println(v ...interface{}) {
	logger.Print(v...)
}

func (logger *SaramaZerologger) logMessage(format string, params ...interface{}) {
	var event *zerolog.Event
	messageStr := fmt.Sprintf(format, params...)

	if strings.HasPrefix(messageStr, kafkaErrorPrefix) {
		event = logger.zerologger.Error()
	} else {
		event = logger.zerologger.Info()
	}

	event = event.Str("package", "sarama")
	event.Msg(messageStr)
}
