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
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/nitrictech/nitric/cloud/common/deploy/tags"
	"github.com/nitrictech/nitric/cloud/gcp/runtime/env"

	"cloud.google.com/go/storage"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/iamcredentials/v1"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"google.golang.org/grpc/codes"

	content "github.com/nitrictech/nitric/cloud/common/runtime/storage"
	ifaces_gcloud_storage "github.com/nitrictech/nitric/cloud/gcp/ifaces/gcloud_storage"
	grpc_errors "github.com/nitrictech/nitric/core/pkg/grpc/errors"
	storagePb "github.com/nitrictech/nitric/core/pkg/proto/storage/v1"
)

type StorageStorageService struct {
	// plugin.UnimplementedStoragePlugin
	client    ifaces_gcloud_storage.StorageClient
	projectID string
	cache     map[string]ifaces_gcloud_storage.BucketHandle
}

var _ storagePb.StorageServer = &StorageStorageService{}

func (s *StorageStorageService) getBucketByName(bucket string) (ifaces_gcloud_storage.BucketHandle, error) {
	if s.cache == nil {
		buckets := s.client.Buckets(context.Background(), s.projectID)
		s.cache = make(map[string]ifaces_gcloud_storage.BucketHandle)
		for {
			b, err := buckets.Next()
			if errors.Is(err, iterator.Done) {
				break
			}
			if err != nil {
				return nil, fmt.Errorf("an error occurred finding bucket: %s; %w", bucket, err)
			}

			if name, ok := b.Labels[tags.GetResourceNameKey(env.GetNitricStackID())]; ok && name == bucket {
				s.cache[name] = s.client.Bucket(b.Name)
			}
		}
	}

	if b, ok := s.cache[bucket]; ok {
		return b, nil
	}

	return nil, fmt.Errorf("bucket not found")
}

// isPermissionDenied returns true if the given error is a Google API error with a StatusForbidden code
func isPermissionDenied(err error) bool {
	var ee *googleapi.Error
	if errors.As(err, &ee) {
		return ee.Code == http.StatusForbidden
	}
	return false
}

/**
 * Retrieves a previously stored object from a Google Cloud Storage Bucket
 */
func (s *StorageStorageService) Read(ctx context.Context, req *storagePb.StorageReadRequest) (*storagePb.StorageReadResponse, error) {
	newErr := grpc_errors.ErrorsWithScope("StorageStorageService.Read")

	bucketHandle, err := s.getBucketByName(req.BucketName)
	if err != nil {
		return nil, newErr(
			codes.NotFound,
			"unable to locate bucket",
			err,
		)
	}

	reader, err := bucketHandle.Object(req.Key).NewReader(ctx)
	if err != nil {
		if isPermissionDenied(err) {
			return nil, newErr(
				codes.PermissionDenied,
				"unable to read file, have you requested access to this bucket?",
				err,
			)
		}

		return nil, newErr(
			codes.Internal,
			"unable to ger reader for object",
			err,
		)
	}
	defer reader.Close()

	bytes, err := io.ReadAll(reader)
	if err != nil {
		if isPermissionDenied(err) {
			return nil, newErr(
				codes.PermissionDenied,
				"unable to read file, have you requested access to this bucket?",
				err,
			)
		}

		return nil, newErr(
			codes.Internal,
			"error reading object stream",
			err,
		)
	}

	return &storagePb.StorageReadResponse{
		Body: bytes,
	}, nil
}

/**
 * Stores a new Item in a Google Cloud Storage Bucket
 */
func (s *StorageStorageService) Write(ctx context.Context, req *storagePb.StorageWriteRequest) (*storagePb.StorageWriteResponse, error) {
	newErr := grpc_errors.ErrorsWithScope("StorageStorageService.Write")

	bucketHandle, err := s.getBucketByName(req.BucketName)
	if err != nil {
		return nil, newErr(
			codes.NotFound,
			"unable to locate bucket",
			err,
		)
	}

	writer := bucketHandle.Object(req.Key).NewWriter(ctx)
	writer.ObjectAttrs().ContentType = content.DetectContentType(req.Key, req.Body)

	if _, err := writer.Write(req.Body); err != nil {
		if isPermissionDenied(err) {
			return nil, newErr(
				codes.PermissionDenied,
				"unable to write to file, have you requested access to this bucket?",
				err,
			)
		}

		return nil, newErr(
			codes.Internal,
			"unable to write object",
			err,
		)
	}

	if err := writer.Close(); err != nil {
		if isPermissionDenied(err) {
			return nil, newErr(
				codes.PermissionDenied,
				"unable to write to file, have you requested access to this bucket?",
				err,
			)
		}

		return nil, newErr(
			codes.Internal,
			"error closing object write",
			err,
		)
	}

	return &storagePb.StorageWriteResponse{}, nil
}

/**
 * Delete an Item in a Google Cloud Storage Bucket
 */
func (s *StorageStorageService) Delete(ctx context.Context, req *storagePb.StorageDeleteRequest) (*storagePb.StorageDeleteResponse, error) {
	newErr := grpc_errors.ErrorsWithScope("StorageStorageService.Delete")

	bucketHandle, err := s.getBucketByName(req.BucketName)
	if err != nil {
		return nil, newErr(
			codes.NotFound,
			"unable to locate bucket",
			err,
		)
	}

	if err := bucketHandle.Object(req.Key).Delete(ctx); err != nil {
		if isPermissionDenied(err) {
			return nil, newErr(
				codes.PermissionDenied,
				"unable to delete to file, have you requested access to this bucket?",
				err,
			)
		}

		// ignore errors caused by the Object not existing.
		// This is to unify delete behavior between providers.
		if !errors.Is(err, storage.ErrObjectNotExist) {
			return nil, newErr(
				codes.NotFound,
				"object does not exist",
				err,
			)
		}
	}

	return &storagePb.StorageDeleteResponse{}, nil
}

func (s *StorageStorageService) PreSignUrl(ctx context.Context, req *storagePb.StoragePreSignUrlRequest) (*storagePb.StoragePreSignUrlResponse, error) {
	newErr := grpc_errors.ErrorsWithScope("StorageStorageService.PreSignedUrl")

	bucketHandle, err := s.getBucketByName(req.BucketName)
	if err != nil {
		return nil, newErr(
			codes.NotFound,
			"unable to locate bucket",
			err,
		)
	}

	method := "GET"
	if req.Operation == storagePb.StoragePreSignUrlRequest_WRITE {
		method = "PUT"
	}

	opts := &storage.SignedURLOptions{
		Scheme:  storage.SigningSchemeV4,
		Method:  method,
		Expires: time.Now().Add(req.Expiry.AsDuration()),
	}

	signedUrl, err := bucketHandle.SignedURL(req.Key, opts)
	if err != nil {
		return nil, newErr(
			codes.Internal,
			"failed to create signed url",
			err,
		)
	}

	return &storagePb.StoragePreSignUrlResponse{
		Url: signedUrl,
	}, nil
}

func (s *StorageStorageService) ListBlobs(ctx context.Context, req *storagePb.StorageListBlobsRequest) (*storagePb.StorageListBlobsResponse, error) {
	newErr := grpc_errors.ErrorsWithScope("StorageStorageService.ListFiles")

	bucketHandle, err := s.getBucketByName(req.BucketName)
	if err != nil {
		return nil, newErr(
			codes.NotFound,
			"unable to locate bucket",
			err,
		)
	}

	iter := bucketHandle.Objects(ctx, &storage.Query{
		Projection: storage.ProjectionNoACL,
		Prefix:     req.Prefix,
	})

	fis := make([]*storagePb.Blob, 0)
	for {
		obj, err := iter.Next()

		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return nil, newErr(codes.Internal, "error occurred iterating objects", err)
		}

		fis = append(fis, &storagePb.Blob{
			Key: obj.Name,
		})
	}

	return &storagePb.StorageListBlobsResponse{
		Blobs: fis,
	}, nil
}

func (s *StorageStorageService) Exists(ctx context.Context, req *storagePb.StorageExistsRequest) (*storagePb.StorageExistsResponse, error) {
	newErr := grpc_errors.ErrorsWithScope("StorageStorageService.Exists")

	bucketHandle, err := s.getBucketByName(req.BucketName)
	if err != nil {
		return nil, newErr(
			codes.NotFound,
			"unable to locate bucket",
			err,
		)
	}

	_, err = bucketHandle.Object(req.Key).Attrs(ctx)
	if errors.Is(err, storage.ErrObjectNotExist) {
		return &storagePb.StorageExistsResponse{
			Exists: false,
		}, nil
	}
	if err != nil {
		return nil, newErr(codes.Internal, "error calling object.Attrs", err)
	}

	return &storagePb.StorageExistsResponse{
		Exists: true,
	}, nil
}

/**
 * Creates a new Storage Plugin for use in GCP
 */
func New() (*StorageStorageService, error) {
	ctx := context.Background()

	credentials, credentialsError := google.FindDefaultCredentials(ctx,
		storage.ScopeReadWrite,
		// required for signing blob urls
		iamcredentials.CloudPlatformScope,
	)
	if credentialsError != nil {
		return nil, fmt.Errorf("GCP credentials error: %w", credentialsError)
	}

	// Get the client credentials
	client, err := storage.NewClient(ctx, option.WithCredentials(credentials))
	if err != nil {
		return nil, fmt.Errorf("storage client error: %w", err)
	}

	return &StorageStorageService{
		client:    ifaces_gcloud_storage.AdaptStorageClient(client),
		projectID: credentials.ProjectID,
	}, nil
}

func NewWithClient(client ifaces_gcloud_storage.StorageClient) (*StorageStorageService, error) {
	return &StorageStorageService{
		client: client,
	}, nil
}
