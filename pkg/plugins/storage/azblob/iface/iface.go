package azblob_service_iface

import (
	"context"
	"io"

	"github.com/Azure/azure-storage-blob-go/azblob"
)

// AzblobServiceUrlIface - Mockable client interface
// for azblob.ServiceUrl
type AzblobServiceUrlIface interface {
	NewContainerURL(string) AzblobContainerUrlIface
}

// AzblobContainerUrlIface - Mockable client interface
// for azblob.ContainerUrl
type AzblobContainerUrlIface interface {
	NewBlockBlobURL(string) AzblobBlockBlobUrlIface
}

// AzblobBlockBlobUrlIface - Mockable client interface
// for azblob.BlockBlobUrl
type AzblobBlockBlobUrlIface interface {
	Download(context.Context, int64, int64, azblob.BlobAccessConditions, bool, azblob.ClientProvidedKeyOptions) (AzblobDownloadResponse, error)
	Upload(context.Context, io.ReadSeeker, azblob.BlobHTTPHeaders, azblob.Metadata, azblob.BlobAccessConditions, azblob.AccessTierType, azblob.BlobTagsMap, azblob.ClientProvidedKeyOptions) (*azblob.BlockBlobUploadResponse, error)
	Delete(context.Context, azblob.DeleteSnapshotsOptionType, azblob.BlobAccessConditions) (*azblob.BlobDeleteResponse, error)
}

// AzblobDownloadResponse - Mockable client interface
// for azblob.DownloadResponse
type AzblobDownloadResponse interface {
	Body(azblob.RetryReaderOptions) io.ReadCloser
}
