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

// Package logger contains the configuration structures needed to configure
// the access to CloudWatch server to sending the log messages there.
package logger

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/RedHatInsights/cloudwatch"
	"github.com/Shopify/sarama"
	zlogsentry "github.com/archdx/zerolog-sentry"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
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
	loggingConf LoggingConfiguration, cloudWatchConf CloudWatchConfiguration, sentryConf SentryLoggingConfiguration, additionalWriters ...io.Writer,
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
		cloudWatchWriter, err := setupCloudwatchLogging(cloudWatchConf)
		if err != nil {
			return err
		}

		writers = append(writers, &WorkaroundForRHIOPS729{Writer: cloudWatchWriter})
	}

	if loggingConf.LoggingToSentryEnabled {
		sentryWriter, err := setupSentryLogging(sentryConf)
		if err != nil {
			return err
		}
		writers = append(writers, sentryWriter)
		needClose = append(needClose, sentryWriter)
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
	logLevel := strings.ToLower(strings.TrimSpace(configuration.LogLevel))

	switch logLevel {
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warn", "warning":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "fatal":
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	}
}

func setupCloudwatchLogging(conf CloudWatchConfiguration) (io.Writer, error) {
	conf.StreamName = strings.ReplaceAll(conf.StreamName, "$HOSTNAME", os.Getenv("HOSTNAME"))
	awsLogLevel := aws.LogOff
	if conf.Debug {
		awsLogLevel = aws.LogDebugWithSigning |
			aws.LogDebugWithSigning |
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

func setupSentryLogging(conf SentryLoggingConfiguration) (io.WriteCloser, error) {
	sentryWriter, err := zlogsentry.New(conf.SentryDSN)
	if err != nil {
		return nil, err
	}

	return sentryWriter, nil
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
