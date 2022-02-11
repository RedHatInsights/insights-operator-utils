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
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
)

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
