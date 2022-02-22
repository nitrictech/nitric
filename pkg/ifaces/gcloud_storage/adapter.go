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

	"cloud.google.com/go/storage"
)

// AdaptClientStorageClient wraps a storage.Client so that it satisfies the Client
func AdaptStorageClient(c *storage.Client) StorageClient {
	return client{c}
}

type (
	client         struct{ *storage.Client }
	bucketHandle   struct{ *storage.BucketHandle }
	objectHandle   struct{ *storage.ObjectHandle }
	bucketIterator struct{ *storage.BucketIterator }
	writer         struct{ *storage.Writer }
	reader         struct{ *storage.Reader }
)

// func (client) embedToIncludeNewMethods()         {}
// func (bucketHandle) embedToIncludeNewMethods()   {}
// func (objectHandle) embedToIncludeNewMethods()   {}
// func (bucketIterator) embedToIncludeNewMethods() {}
// func (writer) embedToIncludeNewMethods()         {}

func (c client) Bucket(name string) BucketHandle {
	return bucketHandle{c.Client.Bucket(name)}
}

func (c client) Buckets(ctx context.Context, projectID string) BucketIterator {
	return bucketIterator{c.Client.Buckets(ctx, projectID)}
}

func (b bucketHandle) Object(name string) ObjectHandle {
	return objectHandle{b.BucketHandle.Object(name)}
}

func (b bucketHandle) SignedURL(object string, opts *storage.SignedURLOptions) (string, error) {
	return b.BucketHandle.SignedURL(object, opts)
}

func (o objectHandle) Key(encryptionKey []byte) ObjectHandle {
	return objectHandle{o.ObjectHandle.Key(encryptionKey)}
}

func (o objectHandle) NewWriter(ctx context.Context) Writer {
	return writer{o.ObjectHandle.NewWriter(ctx)}
}

func (o objectHandle) NewReader(ctx context.Context) (Reader, error) {
	newReader, err := o.ObjectHandle.NewReader(ctx)
	return reader{newReader}, err
}

func (o objectHandle) Delete(ctx context.Context) error {
	return o.ObjectHandle.Delete(ctx)
}
