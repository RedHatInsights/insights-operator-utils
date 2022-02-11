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
	"sort"

	"github.com/aws/aws-sdk-go/service/s3"
)

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
