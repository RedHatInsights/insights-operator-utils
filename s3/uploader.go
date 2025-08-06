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

package s3util

// Documentation in literate-programming-style is available at:
// https://redhatinsights.github.io/insights-operator-utils/packages/s3/uploader.html

import (
	"bytes"
	"context"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// UploadObject uploads a slice of bytes to a specific path and S3 bucket.
func UploadObject(ctx context.Context, client PutObjectAPIClient, bucket, dst string, src []byte) error {
	_, err := client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:        aws.String(bucket),                      // Bucket to be used
		Key:           aws.String(dst),                         // Name of the file to be saved
		Body:          bytes.NewReader(src),                    // File content
		ContentLength: aws.Int64(int64(len(src))),              // File size
		ContentType:   aws.String(http.DetectContentType(src)), // File content
	})

	return err
}
