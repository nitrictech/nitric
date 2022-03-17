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
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"

	"github.com/nitrictech/nitric/pkg/plugins/errors"
	"github.com/nitrictech/nitric/pkg/plugins/errors/codes"
	"github.com/nitrictech/nitric/pkg/plugins/storage"
	"github.com/nitrictech/nitric/pkg/providers/aws/core"
	"github.com/nitrictech/nitric/pkg/utils"
)

const (
	// ErrCodeNoSuchTagSet - AWS API neglects to include a constant for this error code.
	ErrCodeNoSuchTagSet = "NoSuchTagSet"
	ErrCodeAccessDenied = "AccessDenied"
)

// S3StorageService - Is the concrete implementation of AWS S3 for the Nitric Storage Plugin
type S3StorageService struct {
	//storage.UnimplementedStoragePlugin
	client   s3iface.S3API
	provider core.AwsProvider
	selector BucketSelector
	storage.UnimplementedStoragePlugin
}

type BucketSelector = func(nitricName string) (*string, error)

func (s *S3StorageService) getBucketName(bucket string) (*string, error) {
	if s.selector != nil {
		return s.selector(bucket)
	}

	buckets, err := s.provider.GetResources(core.AwsResource_Bucket)

	if err != nil {
		return nil, fmt.Errorf("error getting bucket list: %v", err)
	}

	if bucketArn, ok := buckets[bucket]; ok {
		bucketName := strings.Split(bucketArn, ":::")[1]

		return aws.String(bucketName), nil
	}

	return nil, fmt.Errorf("bucket %s does not exist", bucket)
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

	if b, err := s.getBucketName(bucket); err == nil {
		resp, err := s.client.GetObject(&s3.GetObjectInput{
			Bucket: b,
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

	if b, err := s.getBucketName(bucket); err == nil {
		contentType := http.DetectContentType(object)

		if _, err := s.client.PutObject(&s3.PutObjectInput{
			Bucket:      b,
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

	if b, err := s.getBucketName(bucket); err == nil {
		// TODO: should we handle delete markers, etc.?
		if _, err := s.client.DeleteObject(&s3.DeleteObjectInput{
			Bucket: b,
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

// PreSignUrl - generates a signed URL which can be used to perform direct operations on a file
// useful for large file uploads/downloads so they can bypass application code and work directly with S3
func (s *S3StorageService) PreSignUrl(bucket string, key string, operation storage.Operation, expiry uint32) (string, error) {
	newErr := errors.ErrorsWithScope(
		"S3StorageService.PreSignUrl",
		map[string]interface{}{
			"bucket":    bucket,
			"key":       key,
			"operation": operation.String(),
		},
	)

	if b, err := s.getBucketName(bucket); err == nil {
		switch operation {
		case storage.READ:
			req, _ := s.client.GetObjectRequest(&s3.GetObjectInput{
				Bucket: b,
				Key:    aws.String(key),
			})
			url, err := req.Presign(time.Duration(expiry) * time.Second)
			if err != nil {
				return "", newErr(
					codes.Internal,
					"failed to generate pre-signed READ URL",
					err,
				)
			}
			return url, err
		case storage.WRITE:
			req, _ := s.client.PutObjectRequest(&s3.PutObjectInput{
				Bucket: b,
				Key:    aws.String(key),
			})
			url, err := req.Presign(time.Duration(expiry) * time.Second)
			if err != nil {
				return "", newErr(
					codes.Internal,
					"failed to generate pre-signed WRITE URL",
					err,
				)
			}
			return url, err
		default:
			return "", fmt.Errorf("requested operation not supported for pre-signed AWS S3 urls")
		}
	} else {
		return "", newErr(
			codes.NotFound,
			"unable to locate bucket",
			err,
		)
	}
}

// New creates a new default S3 storage plugin
func New(provider core.AwsProvider) (storage.StorageService, error) {
	awsRegion := utils.GetEnv("AWS_REGION", "us-east-1")

	sess, sessionError := session.NewSession(&aws.Config{
		// FIXME: Use ENV configuration
		Region: aws.String(awsRegion),
	})

	if sessionError != nil {
		return nil, fmt.Errorf("error creating new AWS session %v", sessionError)
	}

	s3Client := s3.New(sess)

	return &S3StorageService{
		client:   s3Client,
		provider: provider,
	}, nil
}

// NewWithClient creates a new S3 Storage plugin and injects the given client
func NewWithClient(provider core.AwsProvider, client s3iface.S3API, opts ...S3StorageServiceOption) (storage.StorageService, error) {
	s3Client := &S3StorageService{
		client:   client,
		provider: provider,
	}

	for _, o := range opts {
		o.Apply(s3Client)
	}

	return s3Client, nil
}
