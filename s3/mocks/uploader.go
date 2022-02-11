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
	"fmt"

	"github.com/aws/aws-sdk-go/service/s3"

	s3util "github.com/RedHatInsights/insights-operator-utils/s3"
)

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
