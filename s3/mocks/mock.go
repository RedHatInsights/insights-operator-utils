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
	"bytes"
	"errors"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"

	s3util "github.com/RedHatInsights/insights-operator-utils/s3"
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

// PutObject returns an empty PutObjectOutput and the mock client Err field, if not nil.
// It also updates MockS3Client.Contents with the new input, if no Err is specified.
func (m *MockS3Client) PutObject(input *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
	if m.Err != nil {
		return &s3.PutObjectOutput{}, m.Err
	}
	if *input.Key == "" {
		return nil, errors.New(EmptyKeyError)
	}
	b := make([]byte, int(*input.ContentLength))
	_, err := input.Body.Read(b)
	if err != nil {
		return nil, fmt.Errorf("%s: %s", s3util.CannotReadError, err.Error())
	}
	m.Contents[*input.Key] = b

	return &s3.PutObjectOutput{}, nil
}

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
func (m *MockS3Client) GetObject(input *s3.GetObjectInput) (*s3.GetObjectOutput, error) {
	if m.Err != nil {
		return nil, m.Err
	}

	if *input.Key == "" {
		return nil, awserr.New(s3.ErrCodeNoSuchKey, "", nil)
	}
	fileContent, exists := m.Contents[*input.Key]
	if !exists {
		return nil, awserr.New(s3.ErrCodeNoSuchKey, "", nil)
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

// DeleteObject returns an empty DeleteObjectOutput object. If the mock client Err field is not nil, returns an error.
func (m *MockS3Client) DeleteObject(input *s3.DeleteObjectInput) (*s3.DeleteObjectOutput, error) {
	if m.Err != nil {
		return nil, m.Err
	}

	if *input.Key == "" {
		return nil, awserr.New(s3.ErrCodeNoSuchKey, "", nil)
	}
	_, exists := m.Contents[*input.Key]
	if !exists {
		return nil, awserr.New(s3.ErrCodeNoSuchKey, "", nil)
	}

	delete(m.Contents, *input.Key)

	return &s3.DeleteObjectOutput{}, m.Err
}

// CopyObject returns an empty CopyObjectOutput object. If the mock client Err field is not nil, returns an error.
func (m *MockS3Client) CopyObject(input *s3.CopyObjectInput) (*s3.CopyObjectOutput, error) {
	if m.Err != nil {
		return nil, m.Err
	}

	inputKey := strings.Replace(*input.CopySource, fmt.Sprintf("/%s/", *input.Bucket), "", 1)

	_, exists := m.Contents[inputKey]
	if !exists {
		return nil, awserr.New(s3.ErrCodeNoSuchKey, "", nil)
	}

	m.Contents[*input.Key] = m.Contents[inputKey]

	return &s3.CopyObjectOutput{}, m.Err
}

// ListObjectsV2 return the list of files it has using ListObjectsV2Output type
func (m *MockS3Client) ListObjectsV2(input *s3.ListObjectsV2Input) (*s3.ListObjectsV2Output, error) {
	m.calls++
	var (
		output       *s3.ListObjectsV2Output
		indexedFiles []string
		isTruncated  bool
		err          error
	)

	if m.MaxCalls != 0 && m.calls == m.MaxCalls {
		return output, ErrMaxKey
	}

	if *input.Delimiter != "" {
		contents := sliceOfStringToS3Folders(m.Folders)
		output = &s3.ListObjectsV2Output{
			CommonPrefixes: contents,
			IsTruncated:    &isTruncated,
		}
	} else {
		indexedFiles, isTruncated, err = listObjects(m.Contents, *input.StartAfter, int(*input.MaxKeys))
		if err != nil {
			return output, err
		}
		contents := sliceOfStringToS3(indexedFiles)
		output = &s3.ListObjectsV2Output{
			Contents:    contents,
			IsTruncated: &isTruncated,
		}
	}

	return output, m.Err
}

func listObjects(contents MockContents, startAfterKey string, maxKeys int) (output []string, isTRuncated bool, err error) {
	var indexOfKey int
	files := contentsToSlice(contents)
	if startAfterKey == "" {
		output = files
	} else {
		indexOfKey, err = findKeyIndex(files, startAfterKey) // err if not found
		if err != nil {
			return
		}
		output = files[indexOfKey+1:]
	}

	if maxKeys == 0 || len(output) < maxKeys {
		return output, false, nil
	}

	return output[:maxKeys], true, nil
}

func contentsToSlice(contents MockContents) []string {
	keys := make([]string, len(contents))

	i := 0
	for k := range contents {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	return keys
}

func sliceOfStringToS3(files []string) (contents []*s3.Object) {
	for i := range files {
		contents = append(contents, &s3.Object{Key: &files[i]})
	}
	return contents
}

func sliceOfStringToS3Folders(folders []string) (contents []*s3.CommonPrefix) {
	for i := range folders {
		contents = append(contents, &s3.CommonPrefix{Prefix: &folders[i]})
	}
	return contents
}

func findKeyIndex(files []string, key string) (int, error) {
	for i, file := range files {
		if file == key {
			return i, nil
		}
	}

	return 0, ErrKeyNotFound
}
