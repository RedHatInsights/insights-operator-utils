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

// Documentation in literate-programming-style is available at:
// https://redhatinsights.github.io/insights-operator-utils/packages/s3/uploader_test.html

import (
	"errors"
	"testing"

	s3util "github.com/RedHatInsights/insights-operator-utils/s3"
	s3mocks "github.com/RedHatInsights/insights-operator-utils/s3/mocks"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/stretchr/testify/assert"
)

var testBody = []byte("test")

func TestUploadObject(t *testing.T) {
	testCases := []testCase{
		{
			description:    "bucket exists",
			errorExpected:  false,
			mockErrorValue: nil,
			file:           testFile,
		},
		{
			description:    "bucket doesn't exist",
			errorExpected:  true,
			mockErrorValue: awserr.New(s3.ErrCodeNoSuchBucket, "", nil),
			errorMsg:       s3.ErrCodeNoSuchBucket,
			file:           testFile,
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
			description:    "empty key",
			errorExpected:  true,
			mockErrorValue: nil,
			errorMsg:       s3mocks.EmptyKeyError,
			file:           "",
		},
		{
			description:    "empty object",
			errorExpected:  true,
			mockErrorValue: nil,
			errorMsg:       s3util.CannotReadError,
			file:           testFile,
			body:           []byte(""),
		},
	}

	mockClient := &s3mocks.MockS3Client{
		Contents: make(s3mocks.MockContents),
	}

	var body []byte
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			mockClient.Err = tc.mockErrorValue
			if tc.body == nil {
				body = testBody
			} else {
				body = tc.body
			}
			err := s3util.UploadObject(mockClient, testBucket, tc.file, body)
			if tc.errorExpected {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.errorMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
