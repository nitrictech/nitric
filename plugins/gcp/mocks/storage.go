package mocks

import (
	"context"
	"fmt"

	"cloud.google.com/go/storage"
	storage_plugin "github.com/nitric-dev/membrane/plugins/gcp/storage/storage"
)

type MockStorageClient struct {
	storage_plugin.Client
	buckets []string
	storage *map[string]map[string][]byte
}

func (s *MockStorageClient) Bucket(name string) storage_plugin.BucketHandle {
	return &MockBucketHandle{
		name:   name,
		client: s,
	}
}

func (s *MockStorageClient) Buckets(ctx context.Context, projectID string) storage_plugin.BucketIterator {
	return &MockBucketIterator{
		idx:    0,
		client: s,
	}
}

type MockBucketIterator struct {
	storage_plugin.BucketIterator
	client *MockStorageClient
	idx    int
}

func (s *MockBucketIterator) Next() (*storage.BucketAttrs, error) {
	if s.idx < len(s.client.buckets) {
		s.idx++
		return &storage.BucketAttrs{
			Name: s.client.buckets[s.idx-1],
			Labels: map[string]string{
				"x-nitric-name": s.client.buckets[s.idx-1],
			},
		}, nil
	}

	return nil, fmt.Errorf("end of the line")
}

type MockBucketHandle struct {
	storage_plugin.BucketHandle
	name   string
	client *MockStorageClient
}

func (s *MockBucketHandle) Object(name string) storage_plugin.ObjectHandle {
	return &MockObjectHandle{
		name:   name,
		bucket: s.name,
		client: s.client,
	}
}

type MockObjectHandle struct {
	storage_plugin.ObjectHandle
	bucket string
	name   string
	client *MockStorageClient
}

func (s *MockObjectHandle) NewWriter(ctx context.Context) storage_plugin.Writer {
	return &MockWriter{
		bucket: s.bucket,
		key:    s.name,
		client: s.client,
	}
}

type MockWriter struct {
	bucket string
	key    string
	client *MockStorageClient
}

func (s *MockWriter) Write(p []byte) (n int, err error) {
	for _, b := range s.client.buckets {
		if s.bucket == b {
			store := *s.client.storage
			// Continue
			// if s.client.storage == nil {
			// 	s.client.storage = make(map[string]map[string][]byte)
			// }

			if store[s.bucket] == nil {
				store[s.bucket] = make(map[string][]byte)
			}
			// Store the item...
			store[s.bucket][s.key] = p
			return len(p), nil
		}
	}

	return -1, fmt.Errorf("Cannot not write to bucket that does not exist")

}

func (s *MockWriter) Close() error {
	return nil
}

func NewStorageClient(buckets []string, storage *map[string]map[string][]byte) storage_plugin.Client {
	return &MockStorageClient{
		buckets: buckets,
		storage: storage,
	}
}
