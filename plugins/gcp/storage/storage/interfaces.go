package storage_plugin

import (
	"context"
	"io"

	"cloud.google.com/go/storage"
)

type Writer interface {
	io.WriteCloser

	// embedToIncludeNewMethods()
}

type ObjectHandle interface {
	NewWriter(context.Context) Writer

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

type Client interface {
	Bucket(string) BucketHandle
	Buckets(context.Context, string) BucketIterator

	// embedToIncludeNewMethods()
}
