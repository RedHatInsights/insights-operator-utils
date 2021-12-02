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

package s3util_test

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	s3util "github.com/RedHatInsights/insights-operator-utils/s3"
	s3mocks "github.com/RedHatInsights/insights-operator-utils/s3/mocks"
)

func TestDownloadObject(t *testing.T) {
	testCases := []testCase{
		{
			description:   "bucket exists",
			errorExpected: false,
			mockContents:  mockFiles,
			file:          testFile,
		},
		{
			description:    "bucket doesn't exist",
			errorExpected:  true,
			mockErrorValue: awserr.New(s3.ErrCodeNoSuchBucket, "", nil),
			errorMsg:       s3.ErrCodeNoSuchBucket,
			file:           testFile,
		},
		{
			description:    "key doesn't exist",
			errorExpected:  true,
			mockErrorValue: nil,
			errorMsg:       s3.ErrCodeNoSuchKey,
			mockContents:   make(s3mocks.MockContents),
			file:           "doesn't exists",
		},
		{
			description:    "empty key",
			errorExpected:  true,
			mockErrorValue: nil,
			errorMsg:       s3.ErrCodeNoSuchKey,
			mockContents:   make(s3mocks.MockContents),
			file:           "",
		},
		{
			description:    "unknown aws error",
			errorExpected:  true,
			mockErrorValue: awserr.New(randomError, "", nil),
			errorMsg:       randomError,
			file:           testFile,
		},
		{
			description:    "unknown error",
			errorExpected:  true,
			mockErrorValue: errors.New(randomError),
			errorMsg:       randomError,
			file:           testFile,
		},
		{
			description:    "bad content",
			errorExpected:  true,
			mockErrorValue: awserr.New(s3util.CannotReadError, "", nil),
			errorMsg:       s3util.CannotReadError,
			file:           testFileInvalidContent,
			mockContents:   mockFiles,
		},
		{
			description:   "error reading contents",
			errorExpected: true,
			errorMsg:      s3util.CannotReadError,
			file:          testFile,
			mockContents:  mockFiles,
			downloadError: errors.New(s3util.CannotReadError),
		},
	}

	mockClient := &s3mocks.MockS3Client{}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			mockClient.Err = tc.mockErrorValue
			mockClient.Contents = tc.mockContents
			mockClient.DownloadError = tc.downloadError
			b, err := s3util.DownloadObject(mockClient, testBucket, tc.file)
			if tc.errorExpected {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.errorMsg)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, b)
				assert.Equal(t, []byte(fileContent), b)
			}
		})
	}
}
