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
