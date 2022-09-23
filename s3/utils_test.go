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
// https://redhatinsights.github.io/insights-operator-utils/packages/s3/utils_test.html

import (
	"errors"
	"testing"

	s3util "github.com/RedHatInsights/insights-operator-utils/s3"
	s3mocks "github.com/RedHatInsights/insights-operator-utils/s3/mocks"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/stretchr/testify/assert"
)

func TestObjectExists(t *testing.T) {
	testCases := []testCase{
		{
			description:     "bucket exists and file exists",
			errorExpected:   false,
			mockErrorValue:  nil,
			shouldFileExist: true,
			file:            testFile,
		},
		{
			description:     "bucket exists and file does not exist",
			errorExpected:   false,
			mockErrorValue:  awserr.New(s3.ErrCodeNoSuchKey, "", nil),
			shouldFileExist: false,
			file:            testFile,
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
	}

	mockClient := &s3mocks.MockS3Client{
		Contents: make(s3mocks.MockContents),
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			mockClient.Err = tc.mockErrorValue
			ok, err := s3util.ObjectExists(mockClient, testBucket, tc.file)
			if tc.errorExpected {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.errorMsg)
				return // don't test the existence of the file
			} else {
				assert.NoError(t, err)
			}
			if tc.shouldFileExist {
				assert.True(t, ok)
			} else {
				assert.False(t, ok)
			}
		})
	}
}
