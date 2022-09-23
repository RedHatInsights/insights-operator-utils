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
// https://redhatinsights.github.io/insights-operator-utils/packages/s3/mocks/lister.html

import (
	collections "github.com/RedHatInsights/insights-operator-utils/collections"
	"github.com/aws/aws-sdk-go/service/s3"
)

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
	files := contentsToSlice(contents)
	if startAfterKey == "" {
		output = files
	} else {
		indexOfKey, found := collections.Index(startAfterKey, files)
		if !found {
			err = ErrKeyNotFound
			return
		}
		output = files[indexOfKey+1:]
	}

	if maxKeys == 0 || len(output) < maxKeys {
		return output, false, nil
	}

	return output[:maxKeys], true, nil
}
