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
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
)

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
