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

package storage_service

import (
	"context"
	"fmt"
	"io/ioutil"

	"cloud.google.com/go/storage"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"

	ifaces_gcloud_storage "github.com/nitrictech/nitric/pkg/ifaces/gcloud_storage"
	"github.com/nitrictech/nitric/pkg/plugins/errors"
	"github.com/nitrictech/nitric/pkg/plugins/errors/codes"
	plugin "github.com/nitrictech/nitric/pkg/plugins/storage"
)

type StorageStorageService struct {
	plugin.UnimplementedStoragePlugin
	client    ifaces_gcloud_storage.StorageClient
	projectID string
}

func (s *StorageStorageService) getBucketByName(bucket string) (ifaces_gcloud_storage.BucketHandle, error) {
	buckets := s.client.Buckets(context.Background(), s.projectID)
	for {
		b, err := buckets.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("an error occurred finding bucket: %s; %v", bucket, err)
		}
		// We'll label the buckets by their name in the nitric.yaml file and use this...
		if b.Labels["x-nitric-name"] == bucket {
			bucketHandle := s.client.Bucket(b.Name)
			return bucketHandle, nil
		}
	}
	return nil, fmt.Errorf("bucket not found")
}

/**
 * Retrieves a previously stored object from a Google Cloud Storage Bucket
 */
func (s *StorageStorageService) Read(bucket string, key string) ([]byte, error) {
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

	reader, err := bucketHandle.Object(key).NewReader(context.Background())
	if err != nil {
		return nil, newErr(
			codes.Internal,
			"unable to ger reader for object",
			err,
		)
	}
	defer reader.Close()

	bytes, err := ioutil.ReadAll(reader)
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
func (s *StorageStorageService) Write(bucket string, key string, object []byte) error {
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

	writer := bucketHandle.Object(key).NewWriter(context.Background())

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
func (s *StorageStorageService) Delete(bucket string, key string) error {
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

	if err := bucketHandle.Object(key).Delete(context.Background()); err != nil {
		// ignore errors caused by the Object not existing.
		// This is to unify delete behavior between providers.
		if err != storage.ErrObjectNotExist {
			return newErr(
				codes.NotFound,
				"object does not exist",
				err,
			)
		}
	}

	return nil
}

/**
 * Creates a new Storage Plugin for use in GCP
 */
func New() (plugin.StorageService, error) {
	ctx := context.Background()

	credentials, credentialsError := google.FindDefaultCredentials(ctx, storage.ScopeReadWrite)
	if credentialsError != nil {
		return nil, fmt.Errorf("GCP credentials error: %v", credentialsError)
	}
	// Get the
	client, err := storage.NewClient(ctx, option.WithCredentials(credentials))

	if err != nil {
		return nil, fmt.Errorf("storage client error: %v", err)
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
