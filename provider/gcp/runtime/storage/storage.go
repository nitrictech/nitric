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
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/storage"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/iamcredentials/v1"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"

	"github.com/nitrictech/nitric/core/pkg/plugins/errors"
	"github.com/nitrictech/nitric/core/pkg/plugins/errors/codes"
	plugin "github.com/nitrictech/nitric/core/pkg/plugins/storage"
	ifaces_gcloud_storage "github.com/nitrictech/nitric/provider/gcp/ifaces/gcloud_storage"
)

type StorageStorageService struct {
	plugin.UnimplementedStoragePlugin
	client    ifaces_gcloud_storage.StorageClient
	projectID string
	cache     map[string]ifaces_gcloud_storage.BucketHandle
}

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

			if name, ok := b.Labels["x-nitric-name"]; ok {
				s.cache[name] = s.client.Bucket(b.Name)
			}
		}
	}

	if b, ok := s.cache[bucket]; ok {
		return b, nil
	}

	return nil, fmt.Errorf("bucket not found")
}

/**
 * Retrieves a previously stored object from a Google Cloud Storage Bucket
 */
func (s *StorageStorageService) Read(ctx context.Context, bucket string, key string) ([]byte, error) {
	newErr := errors.ErrorsWithScope(
		"StorageStorageService.Read",
		map[string]interface{}{
			"bucket": bucket,
			"key":    key,
		},
	)

	bucketHandle, err := s.getBucketByName(bucket)
	if err != nil {
		return nil, newErr(
			codes.NotFound,
			"unable to locate bucket",
			err,
		)
	}

	reader, err := bucketHandle.Object(key).NewReader(ctx)
	if err != nil {
		return nil, newErr(
			codes.Internal,
			"unable to ger reader for object",
			err,
		)
	}
	defer reader.Close()

	bytes, err := io.ReadAll(reader)
	if err != nil {
		return nil, newErr(
			codes.Internal,
			"error reading object stream",
			err,
		)
	}

	return bytes, nil
}

/**
 * Stores a new Item in a Google Cloud Storage Bucket
 */
func (s *StorageStorageService) Write(ctx context.Context, bucket string, key string, object []byte) error {
	newErr := errors.ErrorsWithScope(
		"StorageStorageService.Write",
		map[string]interface{}{
			"bucket":     bucket,
			"key":        key,
			"object.len": len(object),
		},
	)

	bucketHandle, err := s.getBucketByName(bucket)
	if err != nil {
		return newErr(
			codes.NotFound,
			"unable to locate bucket",
			err,
		)
	}

	writer := bucketHandle.Object(key).NewWriter(ctx)

	if _, err := writer.Write(object); err != nil {
		return newErr(
			codes.Internal,
			"unable to write object",
			err,
		)
	}

	if err := writer.Close(); err != nil {
		return newErr(
			codes.Internal,
			"error closing object write",
			err,
		)
	}

	return nil
}

/**
 * Delete an Item in a Google Cloud Storage Bucket
 */
func (s *StorageStorageService) Delete(ctx context.Context, bucket string, key string) error {
	newErr := errors.ErrorsWithScope(
		"StorageStorageService.Delete",
		map[string]interface{}{
			"bucket": bucket,
			"key":    key,
		},
	)

	bucketHandle, err := s.getBucketByName(bucket)
	if err != nil {
		return newErr(
			codes.NotFound,
			"unable to locate bucket",
			err,
		)
	}

	if err := bucketHandle.Object(key).Delete(ctx); err != nil {
		// ignore errors caused by the Object not existing.
		// This is to unify delete behavior between providers.
		if !errors.Is(err, storage.ErrObjectNotExist) {
			return newErr(
				codes.NotFound,
				"object does not exist",
				err,
			)
		}
	}

	return nil
}

func (s *StorageStorageService) PreSignUrl(ctx context.Context, bucket string, key string, operation plugin.Operation, expiry uint32) (string, error) {
	newErr := errors.ErrorsWithScope(
		"StorageStorageService.PreSignedUrl",
		map[string]interface{}{
			"bucket": bucket,
			"key":    key,
		},
	)

	bucketHandle, err := s.getBucketByName(bucket)
	if err != nil {
		return "", newErr(
			codes.NotFound,
			"unable to locate bucket",
			err,
		)
	}

	method := "GET"
	if operation == plugin.WRITE {
		method = "PUT"
	}

	opts := &storage.SignedURLOptions{
		Scheme:  storage.SigningSchemeV4,
		Method:  method,
		Expires: time.Now().Add(time.Duration(expiry) * time.Second),
	}

	signedUrl, err := bucketHandle.SignedURL(key, opts)
	if err != nil {
		return "", newErr(
			codes.Internal,
			"failed to create signed url",
			err,
		)
	}

	return signedUrl, nil
}

func (s *StorageStorageService) ListFiles(ctx context.Context, bucket string) ([]*plugin.FileInfo, error) {
	newErr := errors.ErrorsWithScope(
		"StorageStorageService.ListFiles",
		map[string]interface{}{
			"bucket": bucket,
		},
	)

	bucketHandle, err := s.getBucketByName(bucket)
	if err != nil {
		return nil, newErr(
			codes.NotFound,
			"unable to locate bucket",
			err,
		)
	}

	iter := bucketHandle.Objects(ctx, &storage.Query{
		Projection: storage.ProjectionNoACL,
	})

	fis := make([]*plugin.FileInfo, 0)
	for {
		obj, err := iter.Next()

		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return nil, newErr(codes.Internal, "error occurred iterating objects", err)
		}

		fis = append(fis, &plugin.FileInfo{
			Key: obj.Name,
		})
	}

	return fis, nil
}

/**
 * Creates a new Storage Plugin for use in GCP
 */
func New() (plugin.StorageService, error) {
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

func NewWithClient(client ifaces_gcloud_storage.StorageClient) (plugin.StorageService, error) {
	return &StorageStorageService{
		client: client,
	}, nil
}
