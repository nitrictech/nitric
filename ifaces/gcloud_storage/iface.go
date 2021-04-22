package ifaces_gcloud_storage

import (
	"context"
	"io"

	"cloud.google.com/go/storage"
)

type Writer interface {
	io.WriteCloser
	// embedToIncludeNewMethods()
}

type Reader interface {
	io.ReadCloser
}

type ObjectHandle interface {
	NewWriter(context.Context) Writer
	NewReader(context.Context) (Reader, error)
	Delete(ctx context.Context) error

	// embedToIncludeNewMethods()
}

type BucketIterator interface {
	Next() (*storage.BucketAttrs, error)

	// embedToIncludeNewMethods()
}

type BucketHandle interface {
	Object(string) ObjectHandle

	// embedToIncludeNewMethods()
}

type StorageClient interface {
	Bucket(string) BucketHandle
	Buckets(context.Context, string) BucketIterator

	// embedToIncludeNewMethods()
}
