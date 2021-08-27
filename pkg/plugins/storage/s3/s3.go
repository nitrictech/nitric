// Copyright 2021 Nitric Pty Ltd.
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

package s3_service

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/nitric-dev/membrane/pkg/utils"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/nitric-dev/membrane/pkg/plugins/errors"
	"github.com/nitric-dev/membrane/pkg/plugins/errors/codes"
	"github.com/nitric-dev/membrane/pkg/plugins/storage"
)

const (
	// AWS API neglects to include a constant for this error code.
	ErrCodeNoSuchTagSet = "NoSuchTagSet"
)

// S3StorageService - Is the concrete implementation of AWS S3 for the Nitric Storage Plugin
type S3StorageService struct {
	storage.UnimplementedStoragePlugin
	client s3iface.S3API
}

// getBucketByName - Finds and returns a bucket by it's Nitric name
func (s *S3StorageService) getBucketByName(bucket string) (*s3.Bucket, error) {
	out, err := s.client.ListBuckets(&s3.ListBucketsInput{})

	if err != nil {
		return nil, fmt.Errorf("Encountered an error retrieving the bucket list: %v", err)
	}

	for _, b := range out.Buckets {
		// TODO: This could be rather slow, it's interesting that they don't return this in the list buckets output
		tagout, err := s.client.GetBucketTagging(&s3.GetBucketTaggingInput{
			Bucket: b.Name,
		})

		if err != nil {
			if awsErr, ok := err.(awserr.Error); ok {
				// Table not found,  try to create and put again
				if awsErr.Code() == ErrCodeNoSuchTagSet {
					// Ignore buckets with no tags, check the next bucket
					continue
				}
				return nil, err
			}
			return nil, err
		}

		for _, tag := range tagout.TagSet {
			if *tag.Key == "x-nitric-name" && *tag.Value == bucket {
				return b, nil
			}
		}
	}

	return nil, fmt.Errorf("Unable to find bucket with name: %s", bucket)
}

// Read - Retrieves an item from a bucket
func (s *S3StorageService) Read(bucket string, key string) ([]byte, error) {
	newErr := errors.ErrorsWithScope(
		"S3StorageService.Read",
		map[string]interface{}{
			"bucket": bucket,
			"key":    key,
		},
	)

	if b, err := s.getBucketByName(bucket); err == nil {
		resp, err := s.client.GetObject(&s3.GetObjectInput{
			Bucket: b.Name,
			Key:    aws.String(key),
		})

		if err != nil {
			return nil, newErr(
				codes.NotFound,
				"error retrieving key",
				err,
			)
		}

		defer resp.Body.Close()
		//TODO: Wrap the possible error from ReadAll
		return ioutil.ReadAll(resp.Body)
	} else {
		return nil, newErr(
			codes.NotFound,
			"unable to locate bucket",
			err,
		)
	}
}

// Write - Writes an item to a bucket
func (s *S3StorageService) Write(bucket string, key string, object []byte) error {
	newErr := errors.ErrorsWithScope(
		"S3StorageService.Write",
		map[string]interface{}{
			"bucket":     bucket,
			"key":        key,
			"object.len": len(object),
		},
	)

	if b, err := s.getBucketByName(bucket); err == nil {
		contentType := http.DetectContentType(object)

		if _, err := s.client.PutObject(&s3.PutObjectInput{
			Bucket:      b.Name,
			Body:        bytes.NewReader(object),
			ContentType: &contentType,
			Key:         aws.String(key),
		}); err != nil {
			return newErr(
				codes.Internal,
				"unable to put object",
				err,
			)
		}
	} else {
		return newErr(
			codes.NotFound,
			"unable to locate bucket",
			err,
		)
	}

	return nil
}

// Delete - Deletes an item from a bucket
func (s *S3StorageService) Delete(bucket string, key string) error {
	newErr := errors.ErrorsWithScope(
		"S3StorageService.Delete",
		map[string]interface{}{
			"bucket": bucket,
			"key":    key,
		},
	)

	if b, err := s.getBucketByName(bucket); err == nil {
		// TODO: should we handle delete markers, etc.?
		if _, err := s.client.DeleteObject(&s3.DeleteObjectInput{
			Bucket: b.Name,
			Key:    aws.String(key),
		}); err != nil {
			return newErr(
				codes.Internal,
				"unable to delete object",
				err,
			)
		}
	} else {
		return newErr(
			codes.NotFound,
			"unable to locate bucket",
			err,
		)
	}

	return nil
}

// New creates a new default S3 storage plugin
func New() (storage.StorageService, error) {
	awsRegion := utils.GetEnv("AWS_REGION", "us-east-1")

	sess, sessionError := session.NewSession(&aws.Config{
		// FIXME: Use ENV configuration
		Region: aws.String(awsRegion),
	})

	if sessionError != nil {
		return nil, fmt.Errorf("Error creating new AWS session %v", sessionError)
	}

	s3Client := s3.New(sess)

	return &S3StorageService{
		client: s3Client,
	}, nil
}

// NewWithClient creates a new S3 Storage plugin and injects the given client
func NewWithClient(client s3iface.S3API) (storage.StorageService, error) {
	return &S3StorageService{
		client: client,
	}, nil
}
