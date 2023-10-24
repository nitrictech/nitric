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

package storage

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aws/smithy-go"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-sdk-go-v2/otelaws"

	"github.com/nitrictech/nitric/cloud/aws/ifaces/s3iface"
	"github.com/nitrictech/nitric/cloud/aws/runtime/core"
	"github.com/nitrictech/nitric/core/pkg/plugins/errors"
	"github.com/nitrictech/nitric/core/pkg/plugins/errors/codes"
	"github.com/nitrictech/nitric/core/pkg/plugins/storage"
	"github.com/nitrictech/nitric/core/pkg/utils"
)

const (
	// ErrCodeNoSuchTagSet - AWS API neglects to include a constant for this error code.
	ErrCodeNoSuchTagSet = "NoSuchTagSet"
	ErrCodeAccessDenied = "AccessDenied"
)

// S3StorageService - Is the concrete implementation of AWS S3 for the Nitric Storage Plugin
type S3StorageService struct {
	// storage.UnimplementedStoragePlugin
	client        s3iface.S3API
	preSignClient s3iface.PreSignAPI
	provider      core.AwsProvider
	selector      BucketSelector
	storage.UnimplementedStoragePlugin
}

type BucketSelector = func(nitricName string) (*string, error)

func (s *S3StorageService) getBucketName(ctx context.Context, bucket string) (*string, error) {
	if s.selector != nil {
		return s.selector(bucket)
	}

	buckets, err := s.provider.GetResources(ctx, core.AwsResource_Bucket)
	if err != nil {
		return nil, fmt.Errorf("error getting bucket list: %w", err)
	}

	if bucketArn, ok := buckets[bucket]; ok {
		bucketName := strings.Split(bucketArn, ":::")[1]

		return aws.String(bucketName), nil
	}

	return nil, fmt.Errorf("bucket %s does not exist", bucket)
}

func isS3AccessDeniedErr(err error) bool {
	var opErr *smithy.OperationError
	if errors.As(err, &opErr) {
		return opErr.Service() == "S3" && strings.Contains(opErr.Unwrap().Error(), "AccessDenied")
	}
	return false
}

// Read - Retrieves an item from a bucket
func (s *S3StorageService) Read(ctx context.Context, bucket string, key string) ([]byte, error) {
	newErr := errors.ErrorsWithScope(
		"S3StorageService.Read",
		map[string]interface{}{
			"bucket": bucket,
			"key":    key,
		},
	)

	if b, err := s.getBucketName(ctx, bucket); err == nil {
		resp, err := s.client.GetObject(ctx, &s3.GetObjectInput{
			Bucket: b,
			Key:    aws.String(key),
		})
		if err != nil {

			if isS3AccessDeniedErr(err) {
				return nil, newErr(
					codes.PermissionDenied,
					"unable to read file, have you requested access to this bucket?",
					err,
				)
			}

			return nil, newErr(
				codes.NotFound,
				"error retrieving key",
				err,
			)
		}

		defer resp.Body.Close()
		// TODO: Wrap the possible error from ReadAll
		return io.ReadAll(resp.Body)
	} else {
		return nil, newErr(
			codes.NotFound,
			"unable to locate bucket",
			err,
		)
	}
}

// Write - Writes an item to a bucket
func (s *S3StorageService) Write(ctx context.Context, bucket string, key string, object []byte) error {
	newErr := errors.ErrorsWithScope(
		"S3StorageService.Write",
		map[string]interface{}{
			"bucket":     bucket,
			"key":        key,
			"object.len": len(object),
		},
	)

	if b, err := s.getBucketName(ctx, bucket); err == nil {
		contentType := http.DetectContentType(object)

		if _, err := s.client.PutObject(ctx, &s3.PutObjectInput{
			Bucket:      b,
			Body:        bytes.NewReader(object),
			ContentType: &contentType,
			Key:         aws.String(key),
		}); err != nil {
			if isS3AccessDeniedErr(err) {
				return newErr(
					codes.PermissionDenied,
					"unable to write file, have you requested access to this bucket?",
					err,
				)
			}

			return newErr(
				codes.Internal,
				"unable to put object"+fmt.Sprintf("Error: %v, Type: %T\n", err, err),
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
func (s *S3StorageService) Delete(ctx context.Context, bucket string, key string) error {
	newErr := errors.ErrorsWithScope(
		"S3StorageService.Delete",
		map[string]interface{}{
			"bucket": bucket,
			"key":    key,
		},
	)

	if b, err := s.getBucketName(ctx, bucket); err == nil {
		// TODO: should we handle delete markers, etc.?
		if _, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
			Bucket: b,
			Key:    aws.String(key),
		}); err != nil {
			if isS3AccessDeniedErr(err) {
				return newErr(
					codes.PermissionDenied,
					"unable to delete file, have you requested access to this bucket?",
					err,
				)
			}

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
func (s *S3StorageService) PreSignUrl(ctx context.Context, bucket string, key string, operation storage.Operation, expiry uint32) (string, error) {
	newErr := errors.ErrorsWithScope(
		"S3StorageService.PreSignUrl",
		map[string]interface{}{
			"bucket":    bucket,
			"key":       key,
			"operation": operation.String(),
		},
	)

	if b, err := s.getBucketName(ctx, bucket); err == nil {
		switch operation {
		case storage.READ:
			req, err := s.preSignClient.PresignGetObject(ctx, &s3.GetObjectInput{
				Bucket: b,
				Key:    aws.String(key),
			}, s3.WithPresignExpires(time.Duration(expiry)*time.Second))
			if err != nil {
				return "", newErr(
					codes.Internal,
					"failed to generate pre-signed READ URL",
					err,
				)
			}
			return req.URL, err
		case storage.WRITE:
			req, err := s.preSignClient.PresignPutObject(ctx, &s3.PutObjectInput{
				Bucket: b,
				Key:    aws.String(key),
			}, s3.WithPresignExpires(time.Duration(expiry)*time.Second))
			if err != nil {
				return "", newErr(
					codes.Internal,
					"failed to generate pre-signed WRITE URL",
					err,
				)
			}
			return req.URL, err
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

func (s *S3StorageService) ListFiles(ctx context.Context, bucket string, options *storage.ListFileOptions) ([]*storage.FileInfo, error) {
	newErr := errors.ErrorsWithScope(
		"S3StorageService.ListFiles",
		map[string]interface{}{
			"bucket": bucket,
		},
	)

	var prefix *string = nil
	if options != nil {
		// Only apply if prefix isn't default
		if options.Prefix != "" {
			prefix = aws.String(options.Prefix)
		}
	}

	if b, err := s.getBucketName(ctx, bucket); err == nil {
		objects, err := s.client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
			Bucket: b,
			Prefix: prefix,
		})
		if err != nil {
			return nil, newErr(
				codes.Internal,
				"unable to fetch file list",
				err,
			)
		}

		files := make([]*storage.FileInfo, 0, len(objects.Contents))
		for _, o := range objects.Contents {
			files = append(files, &storage.FileInfo{
				Key: *o.Key,
			})
		}

		return files, nil
	} else {
		return nil, newErr(
			codes.NotFound,
			"unable to locate bucket",
			err,
		)
	}
}

func (s *S3StorageService) Exists(ctx context.Context, bucket string, key string) (bool, error) {
	newErr := errors.ErrorsWithScope(
		"S3StorageService.Exists",
		map[string]interface{}{
			"bucket": bucket,
			"key":    key,
		},
	)

	b, err := s.getBucketName(ctx, bucket)
	if err != nil {
		return false, newErr(codes.Internal, "unable to locate bucket", err)
	}

	_, err = s.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: b,
		Key:    aws.String(key),
	})

	// TODO: Handle specific error types
	if err != nil {
		return false, nil
	}

	return true, nil
}

// New creates a new default S3 storage plugin
func New(provider core.AwsProvider) (storage.StorageService, error) {
	awsRegion := utils.GetEnv("AWS_REGION", "us-east-1")

	cfg, sessionError := config.LoadDefaultConfig(context.TODO(), config.WithRegion(awsRegion))
	if sessionError != nil {
		return nil, fmt.Errorf("error creating new AWS session %w", sessionError)
	}

	otelaws.AppendMiddlewares(&cfg.APIOptions)

	s3Client := s3.NewFromConfig(cfg)

	return &S3StorageService{
		client:        s3Client,
		preSignClient: s3.NewPresignClient(s3Client),
		provider:      provider,
	}, nil
}

// NewWithClient creates a new S3 Storage plugin and injects the given client
func NewWithClient(provider core.AwsProvider, client s3iface.S3API, preSignClient s3iface.PreSignAPI, opts ...S3StorageServiceOption) (storage.StorageService, error) {
	s3Client := &S3StorageService{
		client:        client,
		preSignClient: preSignClient,
		provider:      provider,
	}

	for _, o := range opts {
		o.Apply(s3Client)
	}

	return s3Client, nil
}
