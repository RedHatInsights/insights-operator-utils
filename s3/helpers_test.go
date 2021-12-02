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

package s3util_test

import s3mocks "github.com/RedHatInsights/insights-operator-utils/s3/mocks"

var (
	testBucket             = "my_bucket"
	testFile               = "my_file"
	fileContent            = "some content"
	randomError            = "an error"
	testFileInvalidContent = "invalid_content"
	mockFiles              = s3mocks.MockContents{
		testFile:               []byte(fileContent),
		testFileInvalidContent: []byte(""),
	}
)

type testCase struct {
	description     string
	errorExpected   bool
	mockErrorValue  error
	shouldFileExist bool
	errorMsg        string
	mockContents    s3mocks.MockContents
	file            string
	body            []byte
	downloadError   error
	wantFiles       []string
	lastKey         string
	maxCalls        int
}
