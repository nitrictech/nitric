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

package azblob_service

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"time"

	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/Azure/go-autorest/autorest/azure"

	"github.com/nitrictech/nitric/pkg/plugins/errors"
	"github.com/nitrictech/nitric/pkg/plugins/errors/codes"
	"github.com/nitrictech/nitric/pkg/plugins/storage"
	azblob_service_iface "github.com/nitrictech/nitric/pkg/plugins/storage/azblob/iface"
	azureutils "github.com/nitrictech/nitric/pkg/providers/azure/utils"
	"github.com/nitrictech/nitric/pkg/utils"
)

// AzblobStorageService - Nitric membrane storage plugin implementation for Azure Storage
type AzblobStorageService struct {
	client azblob_service_iface.AzblobServiceUrlIface
	storage.UnimplementedStoragePlugin
}

func (a *AzblobStorageService) getBlobUrl(bucket string, key string) azblob_service_iface.AzblobBlockBlobUrlIface {
	cUrl := a.client.NewContainerURL(bucket)
	// Get a new blob for the key name
	return cUrl.NewBlockBlobURL(key)
}

func (a *AzblobStorageService) Read(bucket string, key string) ([]byte, error) {
	newErr := errors.ErrorsWithScope(
		"AzblobStorageService.Read",
		map[string]interface{}{
			"bucket": bucket,
			"key":    key,
		},
	)
	// Get the bucket for this bucket name
	blob := a.getBlobUrl(bucket, key)
	//// download the blob
	r, err := blob.Download(
		context.TODO(),
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

	// TODO: Configure retries
	data := r.Body(azblob.RetryReaderOptions{MaxRetryRequests: 20})

	return ioutil.ReadAll(data)
}

func (a *AzblobStorageService) Write(bucket string, key string, object []byte) error {
	newErr := errors.ErrorsWithScope(
		"AzblobStorageService.Write",
		map[string]interface{}{
			"bucket": bucket,
			"key":    key,
		},
	)

	blob := a.getBlobUrl(bucket, key)

	if _, err := blob.Upload(
		context.TODO(),
		bytes.NewReader(object),
		azblob.BlobHTTPHeaders{},
		azblob.Metadata{},
		azblob.BlobAccessConditions{},
		azblob.DefaultAccessTier,
		nil,
		azblob.ClientProvidedKeyOptions{},
	); err != nil {
		return newErr(
			codes.Internal,
			"Unable to write blob data",
			err,
		)
	}

	return nil
}

func (a *AzblobStorageService) Delete(bucket string, key string) error {
	newErr := errors.ErrorsWithScope(
		"AzblobStorageService.Delete",
		map[string]interface{}{
			"bucket": bucket,
			"key":    key,
		},
	)

	// Get the bucket for this bucket name
	blob := a.getBlobUrl(bucket, key)

	if _, err := blob.Delete(
		context.TODO(),
		azblob.DeleteSnapshotsOptionInclude,
		azblob.BlobAccessConditions{},
	); err != nil {
		return newErr(
			codes.Internal,
			"Unable to delete blob",
			err,
		)
	}

	return nil
}

func (s *AzblobStorageService) PreSignUrl(bucket string, key string, operation storage.Operation, expiry uint32) (string, error) {
	newErr := errors.ErrorsWithScope(
		"AzblobStorageService.PreSignUrl",
		map[string]interface{}{
			"bucket":    bucket,
			"key":       key,
			"operation": operation.String(),
		},
	)

	blobUrlParts := azblob.NewBlobURLParts(s.getBlobUrl(bucket, key).Url())
	currentTime := time.Now().UTC()
	validDuration := currentTime.Add(time.Duration(expiry) * time.Second)
	cred, err := s.client.GetUserDelegationCredential(context.TODO(), azblob.NewKeyInfo(currentTime, validDuration), nil, nil)

	if err != nil {
		return "", newErr(
			codes.Internal,
			"could not get user delegation credential",
			err,
		)
	}

	sigOpts := azblob.BlobSASSignatureValues{
		Protocol:   azblob.SASProtocolHTTPS,
		ExpiryTime: validDuration,
		Permissions: azblob.BlobSASPermissions{
			Read:  operation == storage.READ,
			Write: operation == storage.WRITE,
		}.String(),
		BlobName:      key,
		ContainerName: bucket,
	}

	queryParams, err := sigOpts.NewSASQueryParameters(cred)

	if err != nil {
		return "", newErr(
			codes.Internal,
			"error signing query params for URL",
			err,
		)
	}

	blobUrlParts.SAS = queryParams
	url := blobUrlParts.URL()

	return url.String(), nil
}

const expiryBuffer = 2 * time.Minute

func tokenRefresherFromSpt(spt *adal.ServicePrincipalToken) azblob.TokenRefresher {
	return func(credential azblob.TokenCredential) time.Duration {
		if err := spt.Refresh(); err != nil {
			log.Default().Println("Error refreshing token: ", err)
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
func New() (storage.StorageService, error) {
	// TODO: Create a default storage account for the stack???
	// XXX: This will limit a membrane wrapped application
	// to accessing a single storage account
	blobEndpoint := utils.GetEnv(azureutils.AZURE_STORAGE_BLOB_ENDPOINT, "")
	if blobEndpoint == "" {
		return nil, fmt.Errorf("failed to determine Azure Storage Blob endpoint, environment variable %s not set", azureutils.AZURE_STORAGE_BLOB_ENDPOINT)
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
