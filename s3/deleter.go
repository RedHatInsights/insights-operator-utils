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
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

// DeleteObject delete the key in the specified bucket.
func DeleteObject(client s3iface.S3API, bucket, file string) error {
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(file),
	}
	_, err := client.DeleteObject(input)
	return err
}

// DeleteObjects delete the keys in the specified bucket.
func DeleteObjects(client s3iface.S3API, bucket string, files []string) error {
	input := &s3.DeleteObjectsInput{
		Bucket: aws.String(bucket),
		Delete: &s3.Delete{
			Objects: filenamesToSliceObjects(files),
		},
	}
	_, err := client.DeleteObjects(input)
	return err
}

func filenamesToSliceObjects(filenames []string) (objects []*s3.ObjectIdentifier) {
	for _, filename := range filenames {
		objects = append(objects, &s3.ObjectIdentifier{Key: aws.String(filename)})
	}
	return objects
}
