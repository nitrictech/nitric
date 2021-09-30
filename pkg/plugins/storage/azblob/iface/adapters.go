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

func AdaptServiceUrl(c azblob.ServiceURL) AzblobServiceUrlIface {
	return serviceUrl{c}
}

func AdaptContainerUrl(c azblob.ContainerURL) AzblobContainerUrlIface {
	return containerUrl{c}
}

func AdaptBlobUrl(c azblob.BlockBlobURL) AzblobBlockBlobUrlIface {
	return blobUrl{c}
}

type (
	serviceUrl   struct{ c azblob.ServiceURL }
	containerUrl struct{ c azblob.ContainerURL }
	blobUrl      struct{ c azblob.BlockBlobURL }
)

func (c serviceUrl) NewContainerURL(bucket string) AzblobContainerUrlIface {
	return AdaptContainerUrl(c.c.NewContainerURL(bucket))
}

func (c serviceUrl) GetUserDelegationCredential(ctx context.Context, info azblob.KeyInfo, timeout *int32, requestID *string) (azblob.UserDelegationCredential, error) {
	return c.c.GetUserDelegationCredential(ctx, info, timeout, requestID)
}

func (c containerUrl) NewBlockBlobURL(blob string) AzblobBlockBlobUrlIface {
	return AdaptBlobUrl(c.c.NewBlockBlobURL(blob))
}

func (c blobUrl) Download(ctx context.Context, offset int64, count int64, bac azblob.BlobAccessConditions, f bool, cpk azblob.ClientProvidedKeyOptions) (AzblobDownloadResponse, error) {
	return c.c.Download(ctx, offset, count, bac, f, cpk)
}

func (c blobUrl) Upload(ctx context.Context, r io.ReadSeeker, h azblob.BlobHTTPHeaders, m azblob.Metadata, bac azblob.BlobAccessConditions, att azblob.AccessTierType, btm azblob.BlobTagsMap, cpk azblob.ClientProvidedKeyOptions) (*azblob.BlockBlobUploadResponse, error) {
	return c.c.Upload(ctx, r, h, m, bac, att, btm, cpk)
}

func (c blobUrl) Delete(ctx context.Context, dot azblob.DeleteSnapshotsOptionType, bac azblob.BlobAccessConditions) (*azblob.BlobDeleteResponse, error) {
	return c.c.Delete(ctx, dot, bac)
}
