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
// https://redhatinsights.github.io/insights-operator-utils/packages/s3/deleter.html

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// DeleteObject delete the key in the specified bucket.
func DeleteObject(ctx context.Context, client DeleteObjectsAPIClient, bucket, file string) error {
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(file),
	}
	_, err := client.DeleteObject(ctx, input)
	return err
}

// DeleteObjects delete the keys in the specified bucket.
func DeleteObjects(ctx context.Context, client DeleteObjectsAPIClient, bucket string, files []string) error {
	input := &s3.DeleteObjectsInput{
		Bucket: aws.String(bucket),
		Delete: &types.Delete{
			Objects: filenamesToSliceObjects(files),
		},
	}
	_, err := client.DeleteObjects(ctx, input)
	return err
}

func filenamesToSliceObjects(filenames []string) (objects []types.ObjectIdentifier) {
	for _, filename := range filenames {
		objects = append(objects, types.ObjectIdentifier{Key: aws.String(filename)})
	}
	return objects
}
