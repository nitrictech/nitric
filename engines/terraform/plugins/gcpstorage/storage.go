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
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"cloud.google.com/go/storage"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/iamcredentials/v1"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	storagepb "github.com/nitrictech/nitric/proto/storage/v2"
	nitricStorage "github.com/nitrictech/nitric/server/runtime/storage"
)

type cloudStorage struct {
	storagepb.UnimplementedStorageServer
	nitricStackId string
	projectID     string
	cache         map[string]*storage.BucketHandle
	client        *storage.Client
}

var _ storagepb.StorageServer = &cloudStorage{}

func (s *cloudStorage) getBucketByName(bucket string) (*storage.BucketHandle, error) {
	if s.cache == nil {
		buckets := s.client.Buckets(context.Background(), s.projectID)
		s.cache = make(map[string]*storage.BucketHandle)
		for {
			b, err := buckets.Next()
			if errors.Is(err, iterator.Done) {
				break
			}
			if err != nil {
				return nil, fmt.Errorf("an error occurred finding bucket: %s; %w", bucket, err)
			}

			if name, ok := b.Labels["x-nitric-name"]; ok && name == bucket {
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

// detectContentType detects the content type of a file based on its extension or content
func detectContentType(filename string, content []byte) string {
	contentType := mime.TypeByExtension(filepath.Ext(filename))
	if contentType != "" {
		return contentType
	}

	return http.DetectContentType(content)
}

/**
 * Retrieves a previously stored object from a Google Cloud Storage Bucket
 */
func (s *cloudStorage) Read(ctx context.Context, req *storagepb.StorageReadRequest) (*storagepb.StorageReadResponse, error) {
	bucketHandle, err := s.getBucketByName(req.BucketName)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "unable to locate bucket")
	}

	reader, err := bucketHandle.Object(req.Key).NewReader(ctx)
	if err != nil {
		if isPermissionDenied(err) {
			return nil, status.Errorf(codes.PermissionDenied, "unable to read file, this may be due to a missing permissions request in your code.")
		}

		return nil, status.Errorf(codes.Unknown, "error reading file: %v", err)
	}

	defer reader.Close()

	bytes, err := io.ReadAll(reader)
	if err != nil {
		if isPermissionDenied(err) {
			return nil, status.Errorf(codes.PermissionDenied, "unable to read file, this may be due to a missing permissions request in your code.")
		}

		return nil, status.Errorf(codes.Unknown, "error reading file: %v", err)
	}

	return &storagepb.StorageReadResponse{
		Body: bytes,
	}, nil
}

/**
 * Stores a new Item in a Google Cloud Storage Bucket
 */
func (s *cloudStorage) Write(ctx context.Context, req *storagepb.StorageWriteRequest) (*storagepb.StorageWriteResponse, error) {
	bucketHandle, err := s.getBucketByName(req.BucketName)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "unable to locate bucket")
	}

	writer := bucketHandle.Object(req.Key).NewWriter(ctx)
	writer.ObjectAttrs.ContentType = detectContentType(req.Key, req.Body)

	if _, err := writer.Write(req.Body); err != nil {
		if isPermissionDenied(err) {
			return nil, status.Errorf(codes.PermissionDenied, "unable to write file, this may be due to a missing permissions request in your code.")
		}

		return nil, status.Errorf(codes.Unknown, "error writing file: %v", err)
	}

	if err := writer.Close(); err != nil {
		if isPermissionDenied(err) {
			return nil, status.Errorf(codes.PermissionDenied, "unable to write file, this may be due to a missing permissions request in your code.")
		}

		return nil, status.Errorf(codes.Unknown, "error closing file: %v", err)
	}

	return &storagepb.StorageWriteResponse{}, nil
}

/**
 * Delete an Item in a Google Cloud Storage Bucket
 */
func (s *cloudStorage) Delete(ctx context.Context, req *storagepb.StorageDeleteRequest) (*storagepb.StorageDeleteResponse, error) {
	bucketHandle, err := s.getBucketByName(req.BucketName)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "unable to locate bucket")
	}

	if err := bucketHandle.Object(req.Key).Delete(ctx); err != nil {
		if isPermissionDenied(err) {
			return nil, status.Errorf(codes.PermissionDenied, "unable to delete file, this may be due to a missing permissions request in your code.")
		}

		// ignore errors caused by the Object not existing.
		// This is to unify delete behavior between providers.
		if !errors.Is(err, storage.ErrObjectNotExist) {
			return nil, status.Errorf(codes.Unknown, "error deleting file: %v", err)
		}
	}

	return &storagepb.StorageDeleteResponse{}, nil
}

func (s *cloudStorage) PreSignUrl(ctx context.Context, req *storagepb.StoragePreSignUrlRequest) (*storagepb.StoragePreSignUrlResponse, error) {
	bucketHandle, err := s.getBucketByName(req.BucketName)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "unable to locate bucket")
	}

	method := "GET"
	if req.Operation == storagepb.StoragePreSignUrlRequest_WRITE {
		method = "PUT"
	}

	opts := &storage.SignedURLOptions{
		Scheme:  storage.SigningSchemeV4,
		Method:  method,
		Expires: time.Now().Add(req.Expiry.AsDuration()),
	}

	signedUrl, err := bucketHandle.SignedURL(req.Key, opts)
	if err != nil {
		if isPermissionDenied(err) {
			return nil, status.Errorf(codes.PermissionDenied, "unable to create signed url, this may be due to a missing permissions request in your code.")
		}

		return nil, status.Errorf(codes.Unknown, "error creating signed url: %v", err)
	}

	return &storagepb.StoragePreSignUrlResponse{
		Url: signedUrl,
	}, nil
}

func (s *cloudStorage) ListBlobs(ctx context.Context, req *storagepb.StorageListBlobsRequest) (*storagepb.StorageListBlobsResponse, error) {
	bucketHandle, err := s.getBucketByName(req.BucketName)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "unable to locate bucket")
	}

	iter := bucketHandle.Objects(ctx, &storage.Query{
		Projection: storage.ProjectionNoACL,
		Prefix:     req.Prefix,
	})

	fis := make([]*storagepb.Blob, 0)
	for {
		obj, err := iter.Next()

		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			if isPermissionDenied(err) {
				return nil, status.Errorf(codes.PermissionDenied, "unable to list blobs, this may be due to a missing permissions request in your code.")
			}

			return nil, status.Errorf(codes.Unknown, "error listing blobs: %v", err)
		}

		fis = append(fis, &storagepb.Blob{
			Key: obj.Name,
		})
	}

	return &storagepb.StorageListBlobsResponse{
		Blobs: fis,
	}, nil
}

func (s *cloudStorage) Exists(ctx context.Context, req *storagepb.StorageExistsRequest) (*storagepb.StorageExistsResponse, error) {
	bucketHandle, err := s.getBucketByName(req.BucketName)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "unable to locate bucket")
	}

	_, err = bucketHandle.Object(req.Key).Attrs(ctx)
	if errors.Is(err, storage.ErrObjectNotExist) {
		return &storagepb.StorageExistsResponse{
			Exists: false,
		}, nil
	}
	if err != nil {
		if isPermissionDenied(err) {
			return nil, status.Errorf(codes.PermissionDenied, "unable to check if file exists, this may be due to a missing permissions request in your code.")
		}

		return nil, status.Errorf(codes.Unknown, "error checking if file exists: %v", err)
	}

	return &storagepb.StorageExistsResponse{
		Exists: true,
	}, nil
}

/**
 * Creates a new Storage Plugin for use in GCP
 */
func Plugin() (nitricStorage.Storage, error) {
	nitricStackId := os.Getenv("NITRIC_STACK_ID")
	if nitricStackId == "" {
		return nil, fmt.Errorf("NITRIC_STACK_ID is not set")
	}

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

	return &cloudStorage{
		client:    client,
		projectID: credentials.ProjectID,
	}, nil
}
