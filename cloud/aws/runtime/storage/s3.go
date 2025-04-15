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
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-sdk-go-v2/otelaws"
	"google.golang.org/grpc/codes"

	"github.com/nitrictech/nitric/cloud/aws/ifaces/s3iface"
	"github.com/nitrictech/nitric/cloud/aws/runtime/env"
	"github.com/nitrictech/nitric/cloud/aws/runtime/resource"
	content "github.com/nitrictech/nitric/cloud/common/runtime/storage"
	grpc_errors "github.com/nitrictech/nitric/core/pkg/grpc/errors"
	storagepb "github.com/nitrictech/nitric/core/pkg/proto/storage/v1"
)

const (
	// ErrCodeNoSuchTagSet - AWS API neglects to include a constant for this error code.
	ErrCodeNoSuchTagSet = "NoSuchTagSet"
	ErrCodeAccessDenied = "AccessDenied"
)

// S3StorageService - an AWS S3 implementation of the Nitric Storage Service
type S3StorageService struct {
	s3Client      s3iface.S3API
	preSignClient s3iface.PreSignAPI
	resolver      resource.AwsResourceResolver
	selector      BucketSelector
}

var _ storagepb.StorageServer = (*S3StorageService)(nil)

type BucketSelector = func(nitricName string) (*string, error)

func (s *S3StorageService) getS3BucketName(ctx context.Context, bucket string) (*string, error) {
	if s.selector != nil {
		return s.selector(bucket)
	}

	buckets, err := s.resolver.GetResources(ctx, resource.AwsResource_Bucket)
	if err != nil {
		return nil, fmt.Errorf("error getting bucket list: %w", err)
	}

	if s3Bucket, ok := buckets[bucket]; ok {
		bucketName := strings.Split(s3Bucket.ARN, ":::")[1]

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

// Read and return the contents of a file in a bucket
func (s *S3StorageService) Read(ctx context.Context, req *storagepb.StorageReadRequest) (*storagepb.StorageReadResponse, error) {
	newErr := grpc_errors.ErrorsWithScope("S3StorageService.Read")

	if s3BucketName, err := s.getS3BucketName(ctx, req.BucketName); err == nil {
		resp, err := s.s3Client.GetObject(ctx, &s3.GetObjectInput{
			Bucket: s3BucketName,
			Key:    aws.String(req.Key),
		})
		if err != nil {
			if isS3AccessDeniedErr(err) {
				return nil, newErr(
					codes.PermissionDenied,
					"unable to read file, this may be due to a missing permissions request in your code.",
					err,
				)
			}

			return nil, newErr(
				codes.Unknown,
				"error reading file",
				err,
			)
		}

		defer resp.Body.Close()
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		return &storagepb.StorageReadResponse{
			Body: bodyBytes,
		}, nil
	} else {
		return nil, newErr(
			codes.NotFound,
			"error finding S3 bucket",
			err,
		)
	}
}

// Write contents to a file in a bucket
func (s *S3StorageService) Write(ctx context.Context, req *storagepb.StorageWriteRequest) (*storagepb.StorageWriteResponse, error) {
	newErr := grpc_errors.ErrorsWithScope("S3StorageService.Write")

	if b, err := s.getS3BucketName(ctx, req.BucketName); err == nil {
		contentType := content.DetectContentType(req.Key, req.Body)

		if _, err := s.s3Client.PutObject(ctx, &s3.PutObjectInput{
			Bucket:      b,
			Body:        bytes.NewReader(req.Body),
			ContentType: &contentType,
			Key:         aws.String(req.Key),
		}); err != nil {
			if isS3AccessDeniedErr(err) {
				return nil, newErr(
					codes.PermissionDenied,
					"unable to write file, this may be due to a missing permissions request in your code.",
					err,
				)
			}

			return nil, newErr(
				codes.Unknown,
				"error writing file",
				err,
			)
		}
	} else {
		return nil, newErr(
			codes.NotFound,
			"error finding S3 bucket",
			err,
		)
	}

	return &storagepb.StorageWriteResponse{}, nil
}

// Delete a file from a bucket
func (s *S3StorageService) Delete(ctx context.Context, req *storagepb.StorageDeleteRequest) (*storagepb.StorageDeleteResponse, error) {
	newErr := grpc_errors.ErrorsWithScope("S3StorageService.Delete")

	if b, err := s.getS3BucketName(ctx, req.BucketName); err == nil {
		if _, err := s.s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
			Bucket: b,
			Key:    aws.String(req.Key),
		}); err != nil {
			if isS3AccessDeniedErr(err) {
				return nil, newErr(
					codes.PermissionDenied,
					"unable to delete file, this may be due to a missing permissions request in your code.",
					err,
				)
			}

			return nil, newErr(
				codes.Unknown,
				"error deleting file",
				err,
			)
		}
	} else {
		return nil, newErr(
			codes.NotFound,
			"error finding S3 bucket",
			err,
		)
	}

	return &storagepb.StorageDeleteResponse{}, nil
}

// PreSignUrl generates a signed URL which can be used to perform direct operations on a file
// useful for large file uploads/downloads so they can bypass application code and work directly with S3
func (s *S3StorageService) PreSignUrl(ctx context.Context, req *storagepb.StoragePreSignUrlRequest) (*storagepb.StoragePreSignUrlResponse, error) {
	newErr := grpc_errors.ErrorsWithScope("S3StorageService.PreSignUrl")

	if b, err := s.getS3BucketName(ctx, req.BucketName); err == nil {
		switch req.Operation {
		case storagepb.StoragePreSignUrlRequest_READ:
			response, err := s.preSignClient.PresignGetObject(ctx, &s3.GetObjectInput{
				Bucket: b,
				Key:    aws.String(req.Key),
			}, s3.WithPresignExpires(req.Expiry.AsDuration()))
			if err != nil {
				return nil, newErr(
					codes.Internal,
					"failed to generate signed READ URL",
					err,
				)
			}
			return &storagepb.StoragePreSignUrlResponse{
				Url: response.URL,
			}, err
		case storagepb.StoragePreSignUrlRequest_WRITE:
			req, err := s.preSignClient.PresignPutObject(ctx, &s3.PutObjectInput{
				Bucket: b,
				Key:    aws.String(req.Key),
			}, s3.WithPresignExpires(req.Expiry.AsDuration()))
			if err != nil {
				return nil, newErr(
					codes.Internal,
					"failed to generate signed WRITE URL",
					err,
				)
			}
			return &storagepb.StoragePreSignUrlResponse{
				Url: req.URL,
			}, err
		default:
			return nil, newErr(codes.Unimplemented, "requested operation not supported for pre-signed AWS S3 URLs", nil)
		}
	} else {
		return nil, newErr(
			codes.NotFound,
			"error finding S3 bucket",
			err,
		)
	}
}

// ListFiles lists all files in a bucket
func (s *S3StorageService) ListBlobs(ctx context.Context, req *storagepb.StorageListBlobsRequest) (*storagepb.StorageListBlobsResponse, error) {
	newErr := grpc_errors.ErrorsWithScope("S3StorageService.ListFiles")

	var prefix *string = nil
	if req.Prefix != "" {
		// Only apply if prefix isn't default
		prefix = &req.Prefix
	}

	if b, err := s.getS3BucketName(ctx, req.BucketName); err == nil {
		objects, err := s.s3Client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
			Bucket: b,
			Prefix: prefix,
		})
		if err != nil {
			if isS3AccessDeniedErr(err) {
				return nil, newErr(
					codes.PermissionDenied,
					"unable to list files, this may be due to a missing permissions request in your code.",
					err,
				)
			}

			return nil, newErr(
				codes.Unknown,
				"error listing files",
				err,
			)
		}

		files := make([]*storagepb.Blob, 0, len(objects.Contents))
		for _, o := range objects.Contents {
			files = append(files, &storagepb.Blob{
				Key: *o.Key,
			})
		}

		return &storagepb.StorageListBlobsResponse{
			Blobs: files,
		}, nil
	} else {
		return nil, newErr(
			codes.NotFound,
			"error finding S3 bucket",
			err,
		)
	}
}

func (s *S3StorageService) Exists(ctx context.Context, req *storagepb.StorageExistsRequest) (*storagepb.StorageExistsResponse, error) {
	newErr := grpc_errors.ErrorsWithScope("S3StorageService.Exists")

	b, err := s.getS3BucketName(ctx, req.BucketName)
	if err != nil {
		return nil, newErr(codes.NotFound, "error finding S3 bucket", err)
	}

	_, err = s.s3Client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: b,
		Key:    aws.String(req.Key),
	})
	if err != nil {
		if isS3AccessDeniedErr(err) {
			return nil, newErr(
				codes.PermissionDenied,
				"unable to check if file exists, this may be due to a missing permissions request in your code.",
				err,
			)
		}

		return &storagepb.StorageExistsResponse{
			Exists: false,
		}, nil
	}

	return &storagepb.StorageExistsResponse{
		Exists: true,
	}, nil
}

// New creates a new default S3 storage plugin
func New(resolver resource.AwsResourceResolver) (*S3StorageService, error) {
	awsRegion := env.AWS_REGION.String()

	cfg, sessionError := config.LoadDefaultConfig(context.TODO(), config.WithRegion(awsRegion))
	if sessionError != nil {
		return nil, fmt.Errorf("error creating new AWS session %w", sessionError)
	}

	otelaws.AppendMiddlewares(&cfg.APIOptions)

	s3Client := s3.NewFromConfig(cfg)

	return &S3StorageService{
		s3Client:      s3Client,
		preSignClient: s3.NewPresignClient(s3Client),
		resolver:      resolver,
	}, nil
}

// NewWithClient creates a new S3 Storage plugin and injects the given client
func NewWithClient(provider resource.AwsResourceResolver, client s3iface.S3API, preSignClient s3iface.PreSignAPI, opts ...S3StorageServiceOption) (*S3StorageService, error) {
	s3Client := &S3StorageService{
		s3Client:      client,
		preSignClient: preSignClient,
		resolver:      provider,
	}

	for _, o := range opts {
		o.Apply(s3Client)
	}

	return s3Client, nil
}
