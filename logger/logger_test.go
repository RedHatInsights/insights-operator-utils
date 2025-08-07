// Copyright 2020, 2021, 2022, 2023 Red Hat, Inc
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

package logger_test

// Documentation in literate-programming-style is available at:
// https://redhatinsights.github.io/insights-operator-utils/packages/logger/logger_test.html

import (
	"bytes"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/IBM/sarama"
	"github.com/RedHatInsights/insights-operator-utils/logger"
	"github.com/RedHatInsights/insights-operator-utils/tests/helpers"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

const (
	testTimeout = 10 * time.Second
)

var (
	cloudWatchConf = logger.CloudWatchConfiguration{
		AWSAccessID:     "access ID",
		AWSSecretKey:    "secret",
		AWSSessionToken: "sess token",
		AWSRegion:       "aws region",
		LogGroup:        "log group",
		StreamName:      "stream name",
		Debug:           false,
	}
)

type RemoteLoggingExpect struct {
	ExpectedMethod   string
	ExpectedTarget   string
	ExpectedBody     string
	ResultStatusCode int
	ResultBody       string
}

// getDescribeLogStreamsEvent returns a mock for the cloudwatchwriter2 library's DescribeLogStreams call
// which doesn't include descending and orderBy parameters
func getDescribeLogStreamsEvent(logStreamName string) RemoteLoggingExpect {
	return RemoteLoggingExpect{
		http.MethodPost,
		"Logs_20140328.DescribeLogStreams",
		`{
			"logGroupName": "` + cloudWatchConf.LogGroup + `",
			"logStreamNamePrefix": "` + logStreamName + `"
		}`,
		http.StatusOK,
		`{
			"logStreams": [
				{
					"arn": "arn:aws:logs:` +
			cloudWatchConf.AWSRegion + `:012345678910:log-group:` + cloudWatchConf.LogGroup +
			`:log-stream:` + logStreamName + `",
					"creationTime": 1,
					"firstEventTimestamp": 2,
					"lastEventTimestamp": 3,
					"lastIngestionTime": 4,
					"logStreamName": "` + logStreamName + `",
					"storedBytes": 100,
					"uploadSequenceToken": "1"
				}
			],
			"nextToken": "token1"
		}`,
	}
}

func TestSaramaZerologger(t *testing.T) {
	const expectedStrInfoLevel = "some random message"
	const expectedErrStrErrorLevel = "kafka: error test error"

	buf := new(bytes.Buffer)

	err := logger.InitZerolog(
		logger.LoggingConfiguration{
			Debug:                      false,
			LogLevel:                   "debug",
			LoggingToCloudWatchEnabled: false,
			LoggingToSentryEnabled:     false,
		},
		logger.CloudWatchConfiguration{},
		logger.SentryLoggingConfiguration{},
		zerolog.New(buf),
	)
	helpers.FailOnError(t, err)

	t.Run("InfoLevel", func(t *testing.T) {
		buf.Reset()

		sarama.Logger.Printf(expectedStrInfoLevel)

		assert.Contains(t, buf.String(), `\"level\":\"info\"`)
		assert.Contains(t, buf.String(), expectedStrInfoLevel)
	})

	t.Run("InfoLevel, Println", func(t *testing.T) {
		buf.Reset()

		sarama.Logger.Println(expectedStrInfoLevel)

		assert.Contains(t, buf.String(), `\"level\":\"info\"`)
		assert.Contains(t, buf.String(), expectedStrInfoLevel)
	})

	t.Run("ErrorLevel", func(t *testing.T) {
		buf.Reset()

		sarama.Logger.Print(expectedErrStrErrorLevel)

		assert.Contains(t, buf.String(), `\"level\":\"error\"`)
		assert.Contains(t, buf.String(), expectedErrStrErrorLevel)
	})

	t.Run("ErrorLevel, Println", func(t *testing.T) {
		buf.Reset()

		sarama.Logger.Println(expectedErrStrErrorLevel)

		assert.Contains(t, buf.String(), `\"level\":\"error\"`)
		assert.Contains(t, buf.String(), expectedErrStrErrorLevel)
	})
}

func TestLoggerSetLogLevel(t *testing.T) {
	logLevels := []string{"debug", "info", "warning", "error"}
	for logLevelIndex, logLevel := range logLevels {
		t.Run(logLevel, func(t *testing.T) {
			buf := new(bytes.Buffer)

			err := logger.InitZerolog(
				logger.LoggingConfiguration{
					Debug:                      false,
					LogLevel:                   logLevel,
					LoggingToCloudWatchEnabled: false,
				},
				logger.CloudWatchConfiguration{},
				logger.SentryLoggingConfiguration{},
				zerolog.New(buf),
			)
			helpers.FailOnError(t, err)

			log.Debug().Msg("debug level")
			log.Info().Msg("info level")
			log.Warn().Msg("warning level")
			log.Error().Msg("error level")

			for i := 0; i < len(logLevels); i++ {
				if i < logLevelIndex {
					assert.NotContains(t, buf.String(), logLevels[i]+" level")
				} else {
					assert.Contains(t, buf.String(), logLevels[i]+" level")
				}
			}
		})
	}
}

func TestWorkaroundForRHIOPS729_Write(t *testing.T) {
	for _, testCase := range []struct {
		Name        string
		StrToWrite  string
		ExpectedStr string
		IsJSON      bool
	}{
		{"NotJSON", "some expected string", "some expected string", false},
		{
			"JSON",
			`{"level": "error", "is_something": true}`,
			`{"LEVEL":"error", "IS_SOMETHING": true}`,
			true,
		},
	} {
		t.Run(testCase.Name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			unJSONWriter := logger.WorkaroundForRHIOPS729{Writer: buf}

			writtenBytes, err := unJSONWriter.Write([]byte(testCase.StrToWrite))
			helpers.FailOnError(t, err)

			assert.Equal(t, writtenBytes, len(testCase.StrToWrite))
			if testCase.IsJSON {
				helpers.AssertStringsAreEqualJSON(t, testCase.ExpectedStr, buf.String())
			} else {
				assert.Equal(t, testCase.ExpectedStr, strings.TrimSpace(buf.String()))
			}
		})
	}
}

func gockExpectLogStreamCreation(t testing.TB, baseURL string) {
	logger.AWSCloudWatchEndpoint = baseURL

	// Mock the CreateLogGroup call that cloudwatchwriter2 makes first
	helpers.GockExpectAPIRequest(t, baseURL, &helpers.APIRequest{
		Method:   http.MethodPost,
		Body:     `{"logGroupName":"test-group"}`,
		Endpoint: "",
		ExtraHeaders: http.Header{
			"X-Amz-Target": []string{"Logs_20140328.CreateLogGroup"},
		},
	}, &helpers.APIResponse{
		StatusCode: http.StatusBadRequest,
		Body:       `{"__type": "ResourceAlreadyExistsException", "message": "The specified log group already exists"}`,
		Headers: map[string]string{
			"Content-Type": "application/x-amz-json-1.1",
		},
	})

	// Mock the CreateLogStream call that cloudwatchwriter2 makes
	helpers.GockExpectAPIRequest(t, baseURL, &helpers.APIRequest{
		Method:   http.MethodPost,
		Body:     `{"logGroupName":"test-group","logStreamName":"test-stream"}`,
		Endpoint: "",
		ExtraHeaders: http.Header{
			"X-Amz-Target": []string{"Logs_20140328.CreateLogStream"},
		},
	}, &helpers.APIResponse{
		StatusCode: http.StatusBadRequest,
		Body:       `{"__type": "ResourceAlreadyExistsException", "message": "The specified log stream already exists"}`,
		Headers: map[string]string{
			"Content-Type": "application/x-amz-json-1.1",
		},
	})

	// Mock DescribeLogStreams call that cloudwatchwriter2 makes
	helpers.GockExpectAPIRequest(t, baseURL, &helpers.APIRequest{
		Method:   http.MethodPost,
		Body:     `{"logGroupName":"test-group","logStreamNamePrefix":"test-stream"}`,
		Endpoint: "",
		ExtraHeaders: http.Header{
			"X-Amz-Target": []string{"Logs_20140328.DescribeLogStreams"},
		},
	}, &helpers.APIResponse{
		StatusCode: http.StatusOK,
		Body:       `{"logStreams":[{"logStreamName":"test-stream","uploadSequenceToken":"1"}]}`,
		Headers: map[string]string{
			"Content-Type": "application/x-amz-json-1.1",
		},
	})
}

// TestInitZerolog_DebugEnabled check if/how instance of zerolog is constructed
// when debug output is enabled.
func TestInitZerolog_DebugEnabled(t *testing.T) {
	defer helpers.CleanAfterGock(t)

	const baseURL = "http://localhost:9999/"
	logger.AWSCloudWatchEndpoint = baseURL

	gockExpectLogStreamCreation(t, baseURL)

	err := logger.InitZerolog(
		logger.LoggingConfiguration{
			Debug:                      true,
			LogLevel:                   "debug",
			LoggingToCloudWatchEnabled: true,
		},
		logger.CloudWatchConfiguration{
			LogGroup:   "test-group",
			StreamName: "test-stream",
		},
		logger.SentryLoggingConfiguration{},
	)
	helpers.FailOnError(t, err)
}

// TestInitZerolog_LogToCloudWatch check if/how instance of zerolog is
// constructed when logging to CloudWatch is enabled.
func TestInitZerolog_LogToCloudWatch(t *testing.T) {
	defer helpers.CleanAfterGock(t)

	const baseURL = "http://localhost:9999/"
	logger.AWSCloudWatchEndpoint = baseURL

	gockExpectLogStreamCreation(t, baseURL)

	err := logger.InitZerolog(
		logger.LoggingConfiguration{
			Debug:                      false,
			LogLevel:                   "debug",
			LoggingToCloudWatchEnabled: true,
		},
		logger.CloudWatchConfiguration{
			LogGroup:   "test-group",
			StreamName: "test-stream",
		},
		logger.SentryLoggingConfiguration{},
	)
	helpers.FailOnError(t, err)
}

// TestInitZerolog_LogToCloudWatch check if/how instance of zerolog is
// constructed when logging to CloudWatch is enabled, including debug output
// for CloudWatch.
func TestInitZerolog_LogToCloudWatchWithDebug(t *testing.T) {
	defer helpers.CleanAfterGock(t)

	const baseURL = "http://localhost:9999/"
	logger.AWSCloudWatchEndpoint = baseURL

	gockExpectLogStreamCreation(t, baseURL)

	err := logger.InitZerolog(
		logger.LoggingConfiguration{
			Debug:                      false,
			LogLevel:                   "debug",
			LoggingToCloudWatchEnabled: true,
		},
		logger.CloudWatchConfiguration{
			LogGroup:   "test-group",
			StreamName: "test-stream",
			Debug:      true,
		},
		logger.SentryLoggingConfiguration{},
	)
	helpers.FailOnError(t, err)
}

func TestLoggingToCloudwatch(t *testing.T) {
	helpers.RunTestWithTimeout(t, func(t testing.TB) {
		defer helpers.CleanAfterGock(t)

		const baseURL = "http://localhost:9999/"
		logger.AWSCloudWatchEndpoint = baseURL

		expects := []RemoteLoggingExpect{
			// cloudwatchwriter2 will try to create the log stream first
			{
				http.MethodPost,
				"Logs_20140328.CreateLogStream",
				`{
					"logGroupName": "` + cloudWatchConf.LogGroup + `",
					"logStreamName": "` + cloudWatchConf.StreamName + `"
				}`,
				http.StatusBadRequest,
				`{
					"__type": "ResourceAlreadyExistsException",
					"message": "The specified log stream already exists"
				}`,
			},
			// Then it will describe the log stream to get the sequence token
			getDescribeLogStreamsEvent(cloudWatchConf.StreamName),
			// Then it will put log events
			{
				http.MethodPost,
				"Logs_20140328.PutLogEvents",
				`{
					"logEvents": [
						{
							"message": "test message text goes right here",
							"timestamp": 1
						}
					],
					"logGroupName": "` + cloudWatchConf.LogGroup + `",
					"logStreamName":"` + cloudWatchConf.StreamName + `",
					"sequenceToken":"1"
				}`,
				http.StatusOK,
				`{"nextSequenceToken":"2"}`,
			},
			{
				http.MethodPost,
				"Logs_20140328.PutLogEvents",
				`{
					"logEvents": [
						{
							"message": "second test message text goes right here",
							"timestamp": 2
						}
					],
					"logGroupName": "` + cloudWatchConf.LogGroup + `",
					"logStreamName":"` + cloudWatchConf.StreamName + `",
					"sequenceToken":"2"
				}`,
				http.StatusOK,
				`{"nextSequenceToken":"3"}`,
			},
			// Additional DescribeLogStreams calls that cloudwatchwriter2 might make
			getDescribeLogStreamsEvent(cloudWatchConf.StreamName),
			getDescribeLogStreamsEvent(cloudWatchConf.StreamName),
			getDescribeLogStreamsEvent(cloudWatchConf.StreamName),
		}

		for _, expect := range expects {
			helpers.GockExpectAPIRequest(t, baseURL, &helpers.APIRequest{
				Method:   expect.ExpectedMethod,
				Body:     expect.ExpectedBody,
				Endpoint: "",
				ExtraHeaders: http.Header{
					"X-Amz-Target": []string{expect.ExpectedTarget},
				},
			}, &helpers.APIResponse{
				StatusCode: expect.ResultStatusCode,
				Body:       expect.ResultBody,
				Headers: map[string]string{
					"Content-Type": "application/x-amz-json-1.1",
				},
			})
		}

		err := logger.InitZerolog(logger.LoggingConfiguration{
			Debug:                      false,
			LogLevel:                   "debug",
			LoggingToCloudWatchEnabled: true,
		},
			cloudWatchConf,
			logger.SentryLoggingConfiguration{},
		)
		helpers.FailOnError(t, err)

		log.Error().Msg("test message")
	}, testTimeout)
}

// TestLoggingToCloudwatch_LogStreamMissing tests log stream name behavior when the log stream is missing
// == expecting HOSTNAME as the log stream name
func TestLoggingToCloudwatch_LogStreamMissing(t *testing.T) {
	helpers.RunTestWithTimeout(t, func(t testing.TB) {
		defer helpers.CleanAfterGock(t)

		const baseURL = "http://localhost:9999/"
		logger.AWSCloudWatchEndpoint = baseURL

		hostname, err := os.Hostname()
		helpers.FailOnError(t, err)

		expects := []RemoteLoggingExpect{
			{
				http.MethodPost,
				"Logs_20140328.CreateLogStream",
				`{
					"logGroupName": "` + cloudWatchConf.LogGroup + `",
					"logStreamName": "` + hostname + `"
				}`,
				http.StatusBadRequest,
				`{
					"__type": "ResourceAlreadyExistsException",
					"message": "The specified log stream already exists"
				}`,
			},
			getDescribeLogStreamsEvent(hostname),
		}

		for _, expect := range expects {
			helpers.GockExpectAPIRequest(t, baseURL, &helpers.APIRequest{
				Method:   expect.ExpectedMethod,
				Body:     expect.ExpectedBody,
				Endpoint: "",
				ExtraHeaders: http.Header{
					"X-Amz-Target": []string{expect.ExpectedTarget},
				},
			}, &helpers.APIResponse{
				StatusCode: expect.ResultStatusCode,
				Body:       expect.ResultBody,
				Headers: map[string]string{
					"Content-Type": "application/x-amz-json-1.1",
				},
			})
		}

		cloudWatchConfCopy := &cloudWatchConf
		cloudWatchConfCopy.StreamName = "" // empty stream name == expeting HOSTNAME as stream name

		err = logger.InitZerolog(logger.LoggingConfiguration{
			Debug:                      false,
			LogLevel:                   "debug",
			LoggingToCloudWatchEnabled: true,
		},
			*cloudWatchConfCopy,
			logger.SentryLoggingConfiguration{},
		)
		helpers.FailOnError(t, err)
	}, testTimeout)
}

func TestInitZerolog_LogToSentry(t *testing.T) {
	sentryConf := logger.SentryLoggingConfiguration{
		SentryDSN:         "http://hash@localhost:9999/project/1",
		SentryEnvironment: "test_environment",
	}

	err := logger.InitZerolog(logger.LoggingConfiguration{
		Debug:                      false,
		LogLevel:                   "debug",
		LoggingToCloudWatchEnabled: false,
		LoggingToSentryEnabled:     true,
	},
		logger.CloudWatchConfiguration{},
		sentryConf,
	)
	helpers.FailOnError(t, err)
}

func TestCloseZerolog(t *testing.T) {
	sentryConf := logger.SentryLoggingConfiguration{
		SentryDSN:         "http://hash@localhost:9999/project/1",
		SentryEnvironment: "test_environment",
	}

	err := logger.InitZerolog(logger.LoggingConfiguration{
		Debug:                      false,
		LogLevel:                   "debug",
		LoggingToCloudWatchEnabled: false,
		LoggingToSentryEnabled:     true,
	},
		logger.CloudWatchConfiguration{},
		sentryConf,
	)
	helpers.FailOnError(t, err)
	logger.CloseZerolog()
}

func TestStdoutLog(t *testing.T) {
	stdOut, stdErr := helpers.CatchingOutputs(t, func() {
		_ = logger.InitZerolog(logger.LoggingConfiguration{
			Debug:                      false,
			UseStderr:                  false,
			LogLevel:                   "debug",
			LoggingToCloudWatchEnabled: false,
			LoggingToSentryEnabled:     false,
		},
			logger.CloudWatchConfiguration{},
			logger.SentryLoggingConfiguration{},
		)

		log.Info().Msg("Hello world")
	})

	assert.Contains(t, string(stdOut), `"message":"Hello world"`)
	assert.NotContains(t, string(stdErr), `"message":"Hello world"`)
}

func TestStderrLog(t *testing.T) {
	stdOut, stdErr := helpers.CatchingOutputs(t, func() {
		_ = logger.InitZerolog(logger.LoggingConfiguration{
			Debug:                      false,
			UseStderr:                  true,
			LogLevel:                   "debug",
			LoggingToCloudWatchEnabled: false,
			LoggingToSentryEnabled:     false,
		},
			logger.CloudWatchConfiguration{},
			logger.SentryLoggingConfiguration{},
		)

		log.Info().Msg("Hello world")
	})

	assert.NotContains(t, string(stdOut), `"message":"Hello world"`)
	assert.Contains(t, string(stdErr), `"message":"Hello world"`)
}

func TestKafkaLogging(t *testing.T) {
	err := logger.InitZerolog(logger.LoggingConfiguration{
		Debug:                      false,
		LogLevel:                   "debug",
		LoggingToCloudWatchEnabled: false,
		LoggingToSentryEnabled:     false,
	},
		logger.CloudWatchConfiguration{},
		logger.SentryLoggingConfiguration{},
	)
	helpers.FailOnError(t, err)
}

func TestConvertLogLevel(t *testing.T) {
	type logLevelTestStruct struct {
		Input       string
		Output      zerolog.Level
		Description string
	}

	logLevelTests := []logLevelTestStruct{
		{
			Description: "debug log level",
			Input:       "debug",
			Output:      zerolog.DebugLevel,
		},
		{
			Description: "info log level",
			Input:       "info",
			Output:      zerolog.InfoLevel,
		},
		{
			Description: "warning log level",
			Input:       "warn",
			Output:      zerolog.WarnLevel,
		},
		{
			Description: "warning log level",
			Input:       "warning",
			Output:      zerolog.WarnLevel,
		},
		{
			Description: "error log level",
			Input:       "error",
			Output:      zerolog.ErrorLevel,
		},
		{
			Description: "fatal log level",
			Input:       "fatal",
			Output:      zerolog.FatalLevel,
		},
	}

	for _, logLevelTest := range logLevelTests {
		t.Run(logLevelTest.Description, func(t *testing.T) {
			level := logger.ConvertLogLevel(logLevelTest.Input)
			assert.Equal(t, logLevelTest.Output, level)
		})
	}
}
