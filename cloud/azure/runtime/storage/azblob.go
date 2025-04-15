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

package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/url"
	"time"

	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/Azure/go-autorest/autorest/azure"
	"google.golang.org/grpc/codes"

	"github.com/nitrictech/nitric/cloud/azure/runtime/env"
	azblob_service_iface "github.com/nitrictech/nitric/cloud/azure/runtime/storage/iface"
	azureutils "github.com/nitrictech/nitric/cloud/azure/runtime/utils"
	content "github.com/nitrictech/nitric/cloud/common/runtime/storage"
	grpc_errors "github.com/nitrictech/nitric/core/pkg/grpc/errors"
	"github.com/nitrictech/nitric/core/pkg/logger"
	storagepb "github.com/nitrictech/nitric/core/pkg/proto/storage/v1"
)

// AzblobStorageService - Nitric storage plugin implementation for Azure Storage
type AzblobStorageService struct {
	client azblob_service_iface.AzblobServiceUrlIface
}

var _ storagepb.StorageServer = &AzblobStorageService{}

func (a *AzblobStorageService) getContainerUrl(bucket string) azblob_service_iface.AzblobContainerUrlIface {
	return a.client.NewContainerURL(bucket)
}

func (a *AzblobStorageService) getBlobUrl(bucket string, key string) azblob_service_iface.AzblobBlockBlobUrlIface {
	return a.getContainerUrl(bucket).NewBlockBlobURL(key)
}

func (a *AzblobStorageService) Read(ctx context.Context, req *storagepb.StorageReadRequest) (*storagepb.StorageReadResponse, error) {
	newErr := grpc_errors.ErrorsWithScope("AzblobStorageService.Read")
	// Get the bucket for this bucket name
	blob := a.getBlobUrl(req.BucketName, req.Key)
	//// download the blob
	r, err := blob.Download(
		ctx,
		0,
		azblob.CountToEnd,
		azblob.BlobAccessConditions{},
		false,
		azblob.ClientProvidedKeyOptions{},
	)
	if err != nil {
		return nil, newErr(
			codes.Internal,
			"Unable to download blob",
			err,
		)
	}

	data := r.Body(azblob.RetryReaderOptions{MaxRetryRequests: 20})

	body, err := io.ReadAll(data)
	if err != nil {
		return nil, newErr(
			codes.Internal,
			"Error reading blob response",
			err,
		)
	}

	return &storagepb.StorageReadResponse{
		Body: body,
	}, nil
}

func (a *AzblobStorageService) Write(ctx context.Context, req *storagepb.StorageWriteRequest) (*storagepb.StorageWriteResponse, error) {
	newErr := grpc_errors.ErrorsWithScope("AzblobStorageService.Write")

	contentType := content.DetectContentType(req.Key, req.Body)

	blob := a.getBlobUrl(req.BucketName, req.Key)

	if _, err := blob.Upload(
		ctx,
		bytes.NewReader(req.Body),
		azblob.BlobHTTPHeaders{
			ContentType: contentType,
		},
		azblob.Metadata{},
		azblob.BlobAccessConditions{},
		azblob.DefaultAccessTier,
		nil,
		azblob.ClientProvidedKeyOptions{},
	); err != nil {
		return nil, newErr(
			codes.Internal,
			"Unable to write blob data",
			err,
		)
	}

	return &storagepb.StorageWriteResponse{}, nil
}

func (a *AzblobStorageService) Delete(ctx context.Context, req *storagepb.StorageDeleteRequest) (*storagepb.StorageDeleteResponse, error) {
	newErr := grpc_errors.ErrorsWithScope("AzblobStorageService.Delete")

	// Get the bucket for this bucket name
	blob := a.getBlobUrl(req.BucketName, req.Key)

	if _, err := blob.Delete(
		context.TODO(),
		azblob.DeleteSnapshotsOptionInclude,
		azblob.BlobAccessConditions{},
	); err != nil {
		return nil, newErr(
			codes.Internal,
			"Unable to delete blob",
			err,
		)
	}

	return &storagepb.StorageDeleteResponse{}, nil
}

func (s *AzblobStorageService) PreSignUrl(ctx context.Context, req *storagepb.StoragePreSignUrlRequest) (*storagepb.StoragePreSignUrlResponse, error) {
	newErr := grpc_errors.ErrorsWithScope("AzblobStorageService.PreSignUrl")

	blobUrlParts := azblob.NewBlobURLParts(s.getBlobUrl(req.BucketName, req.Key).Url())
	currentTime := time.Now().UTC()
	validDuration := currentTime.Add(req.Expiry.AsDuration())
	cred, err := s.client.GetUserDelegationCredential(ctx, azblob.NewKeyInfo(currentTime, validDuration), nil, nil)
	if err != nil {
		return nil, newErr(
			codes.Internal,
			"could not get user delegation credential",
			err,
		)
	}

	sigOpts := azblob.BlobSASSignatureValues{
		Protocol:   azblob.SASProtocolHTTPS,
		ExpiryTime: validDuration,
		Permissions: azblob.BlobSASPermissions{
			Read:  req.Operation == storagepb.StoragePreSignUrlRequest_READ,
			Write: req.Operation == storagepb.StoragePreSignUrlRequest_WRITE,
		}.String(),
		BlobName:      req.Key,
		ContainerName: req.BucketName,
	}

	queryParams, err := sigOpts.NewSASQueryParameters(cred)
	if err != nil {
		return nil, newErr(
			codes.Internal,
			"error signing query params for URL",
			err,
		)
	}

	blobUrlParts.SAS = queryParams
	url := blobUrlParts.URL()

	return &storagepb.StoragePreSignUrlResponse{
		Url: url.String(),
	}, nil
}

func (s *AzblobStorageService) ListBlobs(ctx context.Context, req *storagepb.StorageListBlobsRequest) (*storagepb.StorageListBlobsResponse, error) {
	newErr := grpc_errors.ErrorsWithScope("AzblobStorageService.ListFiles")

	cUrl := s.getContainerUrl(req.BucketName)
	files := make([]*storagepb.Blob, 0)

	// List the blob(s) in our container; since a container may hold millions of blobs, this is done 1 segment at a time.
	for marker := (azblob.Marker{}); marker.NotDone(); { // The parens around Marker{} are required to avoid compiler error.
		// Get a result segment starting with the blob indicated by the current Marker.
		listBlob, err := cUrl.ListBlobsFlatSegment(ctx, marker, azblob.ListBlobsSegmentOptions{
			Prefix: req.Prefix,
		})
		if err != nil {
			return nil, newErr(codes.Internal, "error listing files", err)
		}
		// IMPORTANT: ListBlobs returns the start of the next segment; you MUST use this to get
		// the next segment (after processing the current result segment).
		marker = listBlob.NextMarker

		// Process the blobs returned in this result segment (if the segment is empty, the loop body won't execute)
		for _, blobInfo := range listBlob.Segment.BlobItems {
			files = append(files, &storagepb.Blob{
				Key: blobInfo.Name,
			})
		}
	}

	return &storagepb.StorageListBlobsResponse{
		Blobs: files,
	}, nil
}

func (s *AzblobStorageService) Exists(ctx context.Context, req *storagepb.StorageExistsRequest) (*storagepb.StorageExistsResponse, error) {
	newErr := grpc_errors.ErrorsWithScope("AzblobStorageService.Exists")

	bUrl := s.getBlobUrl(req.BucketName, req.Key)

	// Call get properties and use error to determine existence
	_, err := bUrl.GetProperties(context.TODO(), azblob.BlobAccessConditions{}, azblob.ClientProvidedKeyOptions{})

	//nolint:all
	if storageErr, ok := err.(azblob.StorageError); ok {
		if storageErr.ServiceCode() == azblob.ServiceCodeBlobNotFound {
			return &storagepb.StorageExistsResponse{
				Exists: false,
			}, nil
		}
	}

	if err != nil {
		return nil, newErr(codes.Internal, "error getting blob properties", err)
	}

	return &storagepb.StorageExistsResponse{
		Exists: true,
	}, nil
}

const expiryBuffer = 2 * time.Minute

func tokenRefresherFromSpt(spt *adal.ServicePrincipalToken) azblob.TokenRefresher {
	return func(credential azblob.TokenCredential) time.Duration {
		if err := spt.Refresh(); err != nil {
			logger.Errorf("Error refreshing token: %s", err)
		} else {
			tkn := spt.Token()
			credential.SetToken(tkn.AccessToken)

			return tkn.Expires().Sub(time.Now().Add(expiryBuffer))
		}

		// Mark the token as already expired
		return time.Duration(0)
	}
}

// New - Creates a new instance of the AzblobStorageService
func New() (*AzblobStorageService, error) {
	blobEndpoint := env.AZURE_STORAGE_BLOB_ENDPOINT.String()
	if blobEndpoint == "" {
		return nil, fmt.Errorf("failed to determine Azure Storage Blob endpoint, environment variable not set")
	}

	spt, err := azureutils.GetServicePrincipalToken(azure.PublicCloud.ResourceIdentifiers.Storage)
	if err != nil {
		return nil, err
	}

	cTkn := azblob.NewTokenCredential(spt.Token().AccessToken, tokenRefresherFromSpt(spt))

	var accountURL *url.URL
	if accountURL, err = url.Parse(blobEndpoint); err != nil {
		return nil, err
	}

	pipeline := azblob.NewPipeline(cTkn, azblob.PipelineOptions{})
	client := azblob.NewServiceURL(*accountURL, pipeline)

	return &AzblobStorageService{
		client: azblob_service_iface.AdaptServiceUrl(client),
	}, nil
}
