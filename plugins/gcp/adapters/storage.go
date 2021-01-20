package adapters

import (
	"context"

	"cloud.google.com/go/storage"
	"github.com/nitric-dev/membrane/plugins/gcp/ifaces"
)

// AdaptClientStorageClient wraps a storage.Client so that it satisfies the Client
func AdaptStorageClient(c *storage.Client) ifaces.StorageClient {
	return client{c}
}

type (
	client         struct{ *storage.Client }
	bucketHandle   struct{ *storage.BucketHandle }
	objectHandle   struct{ *storage.ObjectHandle }
	bucketIterator struct{ *storage.BucketIterator }
	writer         struct{ *storage.Writer }
)

// func (client) embedToIncludeNewMethods()         {}
// func (bucketHandle) embedToIncludeNewMethods()   {}
// func (objectHandle) embedToIncludeNewMethods()   {}
// func (bucketIterator) embedToIncludeNewMethods() {}
// func (writer) embedToIncludeNewMethods()         {}

func (c client) Bucket(name string) ifaces.BucketHandle {
	return bucketHandle{c.Client.Bucket(name)}
}

func (c client) Buckets(ctx context.Context, projectID string) ifaces.BucketIterator {
	return bucketIterator{c.Client.Buckets(ctx, projectID)}
}

func (b bucketHandle) Object(name string) ifaces.ObjectHandle {
	return objectHandle{b.BucketHandle.Object(name)}
}

func (o objectHandle) Key(encryptionKey []byte) ifaces.ObjectHandle {
	return objectHandle{o.ObjectHandle.Key(encryptionKey)}
}

func (o objectHandle) NewWriter(ctx context.Context) ifaces.Writer {
	return writer{o.ObjectHandle.NewWriter(ctx)}
}
