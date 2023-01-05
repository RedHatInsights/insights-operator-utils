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
// https://redhatinsights.github.io/insights-operator-utils/packages/s3/lister_test.html

import (
	"testing"

	s3util "github.com/RedHatInsights/insights-operator-utils/s3"
	s3mocks "github.com/RedHatInsights/insights-operator-utils/s3/mocks"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	maxKeys      int64 = 2
	mockContents       = s3mocks.MockContents{
		"folder1/key1": []byte(""),
		"folder1/key2": []byte(""),
		"folder2/key1": []byte(""),
		"folder3/key1": []byte(""),
		"folder3/key2": []byte(""),
	}
	wantKeys = []string{
		"folder1/key1",
		"folder1/key2",
		"folder2/key1",
		"folder3/key1",
		"folder3/key2",
	}

	wantFolders = []string{
		"folder1",
		"folder2",
		"folder3",
	}
)

func TestListNObjectsInBucket(t *testing.T) {

	testCases := []testCase{
		{
			description:  "bucket exists",
			mockContents: mockContents,
			wantFiles:    wantKeys[:maxKeys],
		},
		{
			description:    "bucket doesn't exist",
			errorExpected:  true,
			mockErrorValue: awserr.New(s3.ErrCodeNoSuchBucket, "", nil),
			errorMsg:       s3.ErrCodeNoSuchBucket,
		},
		{
			description:  "use a last key",
			mockContents: mockContents,
			wantFiles:    wantKeys[1 : 1+maxKeys],
			lastKey:      wantKeys[0],
		},
		{
			description:   "startKey isn't in the bucket",
			errorExpected: true,
			errorMsg:      s3mocks.ErrKeyNotFound.Error(),
			lastKey:       "not present",
		},
	}

	mockClient := &s3mocks.MockS3Client{}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			mockClient.Err = tc.mockErrorValue
			mockClient.Contents = tc.mockContents
			got, _, err := s3util.ListNObjectsInBucket(mockClient, testBucket, "", tc.lastKey, "", maxKeys)
			if tc.errorExpected {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.errorMsg)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				assert.Equal(t, tc.wantFiles, got)
			}
		})
	}
}

func TestListBucket(t *testing.T) {
	testCases := []testCase{
		{
			description:   "bucket exists",
			errorExpected: false,
			mockContents:  mockContents,
		},
		{
			description:    "bucket doesn't exist",
			errorExpected:  true,
			mockErrorValue: awserr.New(s3.ErrCodeNoSuchBucket, "", nil),
			errorMsg:       s3.ErrCodeNoSuchBucket,
		},
		{
			description:   "error in bucket in the second iteration",
			errorExpected: true,
			errorMsg:      s3mocks.ErrMaxKey.Error(),
			maxCalls:      2,
			mockContents:  mockContents,
		},
	}

	for _, tc := range testCases {
		mockClient := &s3mocks.MockS3Client{}
		mockClient.Err = tc.mockErrorValue
		mockClient.Contents = tc.mockContents
		mockClient.MaxCalls = tc.maxCalls
		t.Run(tc.description, func(t *testing.T) {
			got, err := s3util.ListBucket(mockClient, testBucket, "", "", maxKeys)
			if tc.errorExpected {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.errorMsg)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				assert.Equal(t, wantKeys, got)
			}
		})
	}
}

func TestListFolders(t *testing.T) {
	testCases := []testCase{
		{
			description:   "bucket exists",
			errorExpected: false,
		},
		{
			description:    "bucket doesn't exist",
			errorExpected:  true,
			mockErrorValue: awserr.New(s3.ErrCodeNoSuchBucket, "", nil),
			errorMsg:       s3.ErrCodeNoSuchBucket,
		},
	}

	for _, tc := range testCases {
		mockClient := &s3mocks.MockS3Client{}
		mockClient.Err = tc.mockErrorValue
		mockClient.Folders = wantFolders
		mockClient.MaxCalls = tc.maxCalls
		t.Run(tc.description, func(t *testing.T) {
			got, err := s3util.ListFolders(mockClient, testBucket, "", "", maxKeys)
			if tc.errorExpected {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.errorMsg)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				assert.Equal(t, wantFolders, got)
			}
		})
	}
}
