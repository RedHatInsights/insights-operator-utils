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
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

// ListNObjectsInBucket returns a slice of the files at an S3 bucket at a given prefix. It starts
// listing the files from `lastKey` and returns up to `maxKeys`.
func ListNObjectsInBucket(client s3iface.S3API, bucket, prefix, lastKey, delimiter string, maxKeys int64) (objects []string, isTruncated bool, err error) {
	input := s3.ListObjectsV2Input{
		Bucket:     aws.String(bucket),
		Prefix:     aws.String(prefix),
		MaxKeys:    &maxKeys,
		StartAfter: aws.String(lastKey),
		Delimiter:  aws.String(delimiter),
	}
	result, err := client.ListObjectsV2(&input)

	if err != nil {
		return nil, false, err
	}

	if delimiter == "" {
		objects = sliceObjectsToSliceString(result.Contents)
	} else {
		objects = sliceCommonPrefixToSliceString(result.CommonPrefixes)
	}
	return objects, *result.IsTruncated, nil
}

// ListBucket returns a slice of the files at an S3 bucket at a given prefix. It lists the
// objects using maxKeys in each iteration, but returns all of the objects, which may be
// higher than maxKeys.
func ListBucket(client s3iface.S3API, bucket, prefix, lastKey string, maxKeys int64) ([]string, error) {
	return listWithDelimiter(client, bucket, prefix, lastKey, "", maxKeys)
}

func sliceObjectsToSliceString(input []*s3.Object) (output []string) {
	for i := range input {
		output = append(output, *input[i].Key)
	}
	return
}

// ListFolders returns the folders stored at `prefix` in the bucket `bucket`.
func ListFolders(client s3iface.S3API, bucket, prefix, lastKey string, maxKeys int64) ([]string, error) {
	return listWithDelimiter(client, bucket, prefix, lastKey, "/", maxKeys)
}

func listWithDelimiter(client s3iface.S3API, bucket, prefix, lastKey, delimiter string, maxKeys int64) ([]string, error) {
	output, isTruncated, err := ListNObjectsInBucket(client, bucket, prefix, lastKey, delimiter, maxKeys)
	if err != nil {
		return []string{}, err
	}

	if isTruncated {
		lastKey = output[len(output)-1]
		newOutput, err := listWithDelimiter(client, bucket, prefix, lastKey, delimiter, maxKeys)
		if err != nil {
			return []string{}, err
		}
		output = append(output, newOutput...)
	}
	return output, nil
}

func sliceCommonPrefixToSliceString(input []*s3.CommonPrefix) (output []string) {
	for i := range input {
		output = append(output, *input[i].Prefix)
	}
	return
}
