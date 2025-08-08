// Copyright 2021, 2022 Red Hat, Inc
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

package s3util

// Documentation in literate-programming-style is available at:
// https://redhatinsights.github.io/insights-operator-utils/packages/s3/downloader.html

import (
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// CannotReadError is used in the tests.
const CannotReadError = "cannot read remote object"

// DownloadObject downloads a file from an S3 bucket given the bucket and key. The
// response is in slice of byte format.
func DownloadObject(ctx context.Context, client GetObjectAPIClient, bucket, src string) ([]byte, error) {
	result, err := client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket), // Bucket to be used
		Key:    aws.String(src),    // Name of the file to be downloaded
	})

	if err != nil {
		return nil, err
	}
	b, err := io.ReadAll(result.Body)
	if err != nil {
		return nil, fmt.Errorf("%s: %s", CannotReadError, err.Error())
	}

	return b, err
}
