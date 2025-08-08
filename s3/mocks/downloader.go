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
// https://redhatinsights.github.io/insights-operator-utils/packages/s3/mocks/downloader.html

import (
	"bytes"
	"context"
	"io"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type mockReadCloser struct {
	err     error
	content []byte
}

// Read will return an error if mockReadCloser.generateError is true,
// otherwise it will Read the contents as normally.
func (m mockReadCloser) Read(p []byte) (int, error) {
	if m.err != nil {
		return 0, m.err
	}
	n, _ := bytes.NewReader(m.content).Read(p)
	_ = m.Close()
	return n, io.EOF
}

// Close won't return an error.
func (m mockReadCloser) Close() error {
	return nil
}

// GetObject returns a GetObjectOutput object with the value of MockS3Client.Contents corresponding
// to that key. If the mock client Err field is not nil, returns an error.
func (m *MockS3Client) GetObject(ctx context.Context, input *s3.GetObjectInput, opts ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
	if m.Err != nil {
		return nil, m.Err
	}

	if *input.Key == "" {
		return nil, &types.NoSuchKey{}
	}
	fileContent, exists := m.Contents[*input.Key]
	if !exists {
		return nil, &types.NoSuchKey{}
	}

	b := mockReadCloser{
		m.DownloadError,
		fileContent,
	}

	size := int64(len(fileContent))

	return &s3.GetObjectOutput{
		Body:          b,
		ContentLength: &size,
	}, m.Err
}
