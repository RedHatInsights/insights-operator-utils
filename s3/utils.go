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
// https://redhatinsights.github.io/insights-operator-utils/packages/s3/utils.html

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/rs/zerolog/log"
)

// ObjectExists returns a boolean of whether the key exists in the bucket.
func ObjectExists(ctx context.Context, client HeadObjectAPIClient, bucket, file string) (bool, error) {
	input := &s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(file),
	}

	_, err := client.HeadObject(ctx, input)
	return checkAwsErr(err)
}

func checkAwsErr(err error) (bool, error) {
	if err != nil {
		var noSuchKey *types.NoSuchKey
		var noSuchBucket *types.NoSuchBucket

		if errors.As(err, &noSuchKey) {
			log.Debug().
				Err(err).
				Str("awserr", "NoSuchKey").
				Msg("File does not exist")
			return false, nil
		} else if errors.As(err, &noSuchBucket) {
			log.Debug().
				Err(err).
				Str("awserr", "NoSuchBucket").
				Msg("Bucket does not exist")
			return false, errors.New("NoSuchBucket")
		} else {
			log.Debug().Err(err).Msg("Unknown AWS error")
			return false, err
		}
	}
	return true, nil
}
