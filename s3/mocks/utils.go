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
// https://redhatinsights.github.io/insights-operator-utils/packages/s3/mocks/utils.html

import (
	"sort"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
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

func sliceOfStringToS3(files []string) (contents []types.Object) {
	for i := range files {
		contents = append(contents, types.Object{Key: aws.String(files[i])})
	}
	return contents
}

func sliceOfStringToS3Folders(folders []string) (contents []types.CommonPrefix) {
	for i := range folders {
		contents = append(contents, types.CommonPrefix{Prefix: aws.String(folders[i])})
	}
	return contents
}
