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

// Documentation in literate-programming-style is available at:
// https://redhatinsights.github.io/insights-operator-utils/packages/s3/mocks/mock.html

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/service/s3"
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
func (m *MockS3Client) HeadObject(ctx context.Context, input *s3.HeadObjectInput, opts ...func(*s3.Options)) (*s3.HeadObjectOutput, error) {
	return &s3.HeadObjectOutput{}, m.Err
}
