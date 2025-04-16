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

package ifaces_gcloud_storage

import (
	"context"
	"io"

	"cloud.google.com/go/storage"
)

type Writer interface {
	io.WriteCloser
	ObjectAttrs() *storage.ObjectAttrs
}

type Reader interface {
	io.ReadCloser
}

type ObjectHandle interface {
	NewWriter(context.Context) Writer
	NewReader(context.Context) (Reader, error)
	Delete(ctx context.Context) error
	Attrs(ctx context.Context) (*storage.ObjectAttrs, error)
}

type BucketIterator interface {
	Next() (*storage.BucketAttrs, error)
}

type ObjectIterator interface {
	Next() (*storage.ObjectAttrs, error)
}

type BucketHandle interface {
	Object(string) ObjectHandle
	Objects(context.Context, *storage.Query) ObjectIterator
	SignedURL(string, *storage.SignedURLOptions) (string, error)
}

type StorageClient interface {
	Bucket(string) BucketHandle
	Buckets(context.Context, string) BucketIterator
}
