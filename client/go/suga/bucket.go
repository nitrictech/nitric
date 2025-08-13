package suga

import (
	"context"
	"fmt"
	"time"

	storagepb "github.com/nitrictech/nitric/proto/storage/v2"
	"google.golang.org/protobuf/types/known/durationpb"
)

// Bucket provides methods for interacting with cloud storage buckets
type Bucket struct {
	name          string
	storageClient storagepb.StorageClient
}

// Read a file from the bucket
func (c *Bucket) Read(key string) ([]byte, error) {
	ctx := context.Background()

	req := &storagepb.StorageReadRequest{
		BucketName: c.name,
		Key:        key,
	}

	res, err := c.storageClient.Read(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to read file from the %s bucket: %w", c.name, err)
	}

	return res.Body, nil
}

// Write a file to the bucket
func (c *Bucket) Write(key string, data []byte) error {
	ctx := context.Background()

	req := &storagepb.StorageWriteRequest{
		BucketName: c.name,
		Key:        key,
		Body:       data,
	}

	_, err := c.storageClient.Write(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to write file to bucket: %w", err)
	}

	return nil
}

// Delete a file from the bucket
func (c *Bucket) Delete(key string) error {
	ctx := context.Background()

	req := &storagepb.StorageDeleteRequest{
		BucketName: c.name,
		Key:        key,
	}

	_, err := c.storageClient.Delete(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to delete file from bucket: %w", err)
	}

	return nil
}

// List files in the bucket with a given prefix
func (c *Bucket) List(prefix string) ([]string, error) {
	ctx := context.Background()

	req := &storagepb.StorageListBlobsRequest{
		BucketName: c.name,
		Prefix:     prefix,
	}

	res, err := c.storageClient.ListBlobs(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to list files in bucket: %w", err)
	}

	keys := make([]string, len(res.Blobs))
	for i, blob := range res.Blobs {
		keys[i] = blob.Key
	}

	return keys, nil
}

// Exists checks if a file exists in the bucket
func (c *Bucket) Exists(key string) (bool, error) {
	ctx := context.Background()

	req := &storagepb.StorageExistsRequest{
		BucketName: c.name,
		Key:        key,
	}

	res, err := c.storageClient.Exists(ctx, req)
	if err != nil {
		return false, fmt.Errorf("failed to check if file exists in bucket: %w", err)
	}

	return res.Exists, nil
}

type Mode int

const (
	ModeRead Mode = iota
	ModeWrite
)

type presignUrlOptions struct {
	mode   Mode
	expiry time.Duration
}

type PresignUrlOption func(opts *presignUrlOptions)

// WithPresignUrlExpiry sets the expiry duration for presigned URLs
func WithPresignUrlExpiry(expiry time.Duration) PresignUrlOption {
	return func(opts *presignUrlOptions) {
		opts.expiry = expiry
	}
}

func getPresignUrlOpts(mode Mode, opts ...PresignUrlOption) *presignUrlOptions {
	defaultOpts := &presignUrlOptions{
		mode:   mode,
		expiry: time.Minute * 5,
	}

	for _, opt := range opts {
		opt(defaultOpts)
	}

	return defaultOpts
}

func (c *Bucket) preSignUrl(key string, opts *presignUrlOptions) (string, error) {
	ctx := context.Background()

	op := storagepb.StoragePreSignUrlRequest_READ

	if opts.mode == ModeWrite {
		op = storagepb.StoragePreSignUrlRequest_WRITE
	}

	req := &storagepb.StoragePreSignUrlRequest{
		BucketName: c.name,
		Key:        key,
		Operation:  op,
		Expiry:     durationpb.New(opts.expiry),
	}

	res, err := c.storageClient.PreSignUrl(ctx, req)
	if err != nil {
		return "", fmt.Errorf("failed to get presigned URL for file: %w", err)
	}

	return res.Url, nil
}

// GetDownloadURL returns a presigned URL for downloading a file from the bucket, with a limited expiry time
func (c *Bucket) GetDownloadURL(key string, opts ...PresignUrlOption) (string, error) {
	optsWithDefaults := getPresignUrlOpts(ModeRead, opts...)

	return c.preSignUrl(key, optsWithDefaults)
}

// GetUploadURL returns a presigned URL for uploading a file to the bucket, with a limited expiry time
func (c *Bucket) GetUploadURL(key string, opts ...PresignUrlOption) (string, error) {
	optsWithDefaults := getPresignUrlOpts(ModeWrite, opts...)

	return c.preSignUrl(key, optsWithDefaults)
}

// NewBucket creates a new client interactive with a named bucket
func NewBucket(storageClient storagepb.StorageClient, bucketName string) *Bucket {
	return &Bucket{
		storageClient: storageClient,
		name:          bucketName,
	}
}
