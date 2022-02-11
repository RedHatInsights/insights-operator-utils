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

package mocks

import (
	"errors"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

const (
	// EmptyKeyError is used in the tests
	EmptyKeyError = "empty key"
)

var (
	// ErrMaxKey is used in the tests
	ErrMaxKey = errors.New("max key error")
	// ErrKeyNotFound is used in the tests
	ErrKeyNotFound = errors.New("key not found")
)

// MockS3Client can be used in tests to mock an S3 client.
type MockS3Client struct {
	s3iface.S3API
	Err           error
	Contents      MockContents
	DownloadError error
	MaxCalls      int
	calls         int
	Folders       []string
}

// MockContents stores the file inside the mocked S3 bucket.
type MockContents map[string][]byte

// HeadObject returns an empty HeadObjectOutput and the mock client Err field.
func (m *MockS3Client) HeadObject(input *s3.HeadObjectInput) (*s3.HeadObjectOutput, error) {
	return &s3.HeadObjectOutput{}, m.Err
}
