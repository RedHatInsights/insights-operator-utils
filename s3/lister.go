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
// https://redhatinsights.github.io/insights-operator-utils/packages/s3/lister.html

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// ListNObjectsInBucket returns a slice of the files at an S3 bucket at a given prefix. It starts
// listing the files from `lastKey` and returns up to `maxKeys`.
func ListNObjectsInBucket(ctx context.Context, client ListObjectsV2APIClient, bucket, prefix, lastKey, delimiter string, maxKeys int32) (objects []string, isTruncated bool, err error) {
	input := s3.ListObjectsV2Input{
		Bucket:     aws.String(bucket),
		Prefix:     aws.String(prefix),
		MaxKeys:    &maxKeys,
		StartAfter: aws.String(lastKey),
		Delimiter:  aws.String(delimiter),
	}
	result, err := client.ListObjectsV2(ctx, &input)

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
func ListBucket(ctx context.Context, client ListObjectsV2APIClient, bucket, prefix, lastKey string, maxKeys int32) ([]string, error) {
	return listWithDelimiter(ctx, client, bucket, prefix, lastKey, "", maxKeys)
}

func sliceObjectsToSliceString(input []types.Object) (output []string) {
	for i := range input {
		output = append(output, *input[i].Key)
	}
	return
}

// ListFolders returns the folders stored at `prefix` in the bucket `bucket`.
func ListFolders(ctx context.Context, client ListObjectsV2APIClient, bucket, prefix, lastKey string, maxKeys int32) ([]string, error) {
	return listWithDelimiter(ctx, client, bucket, prefix, lastKey, "/", maxKeys)
}

func listWithDelimiter(ctx context.Context, client ListObjectsV2APIClient, bucket, prefix, lastKey, delimiter string, maxKeys int32) ([]string, error) {
	output, isTruncated, err := ListNObjectsInBucket(ctx, client, bucket, prefix, lastKey, delimiter, maxKeys)
	if err != nil {
		return []string{}, err
	}

	if isTruncated {
		lastKey = output[len(output)-1]
		newOutput, err := listWithDelimiter(ctx, client, bucket, prefix, lastKey, delimiter, maxKeys)
		if err != nil {
			return []string{}, err
		}
		output = append(output, newOutput...)
	}
	return output, nil
}

func sliceCommonPrefixToSliceString(input []types.CommonPrefix) (output []string) {
	for i := range input {
		output = append(output, *input[i].Prefix)
	}
	return
}
