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
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/rs/zerolog/log"
)

// ObjectExists returns a boolean of whether the key exists in the bucket.
func ObjectExists(client s3iface.S3API, bucket, file string) (bool, error) {
	input := &s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(file),
	}

	_, err := client.HeadObject(input)
	return checkAwsErr(err)
}

func checkAwsErr(err error) (bool, error) {
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchKey:
				log.Debug().
					Err(aerr).
					Str("awserr", s3.ErrCodeNoSuchKey).
					Msg("File does not exist")
				return false, nil
			case s3.ErrCodeNoSuchBucket:
				log.Debug().
					Err(aerr).
					Str("awserr", s3.ErrCodeNoSuchBucket).
					Msg("Bucket does not exist")
				return false, errors.New(s3.ErrCodeNoSuchBucket)
			default:
				log.Debug().Err(aerr).Msg("Unknown AWS error")
				return false, aerr
			}
		}
		return false, err
	}
	return true, nil
}
