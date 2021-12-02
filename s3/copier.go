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

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

// CopyObject copies the `inputKey` in the inputBucket to `outputKey` in the outputBucket.
func CopyObject(client s3iface.S3API, inputBucket, inputKey, outputBucket, outputKey string) error {
	input := &s3.CopyObjectInput{
		Bucket:     aws.String(outputBucket),
		Key:        aws.String(outputKey),
		CopySource: aws.String(fmt.Sprintf("/%s/%s", inputBucket, inputKey)),
	}
	_, err := client.CopyObject(input)
	return err
}

// RenameObject renames the `inputKey` in the bucket to `outputKey`.
func RenameObject(client s3iface.S3API, bucket, inputKey, outputKey string) error {
	if err := CopyObject(client, bucket, inputKey, bucket, outputKey); err != nil {
		return err
	}

	return DeleteObject(client, bucket, inputKey)
}
