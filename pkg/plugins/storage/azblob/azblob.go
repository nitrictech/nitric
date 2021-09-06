package azblob_service

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/url"
	"time"

	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/nitric-dev/membrane/pkg/plugins/errors"
	"github.com/nitric-dev/membrane/pkg/plugins/errors/codes"
	"github.com/nitric-dev/membrane/pkg/plugins/storage"
	azblob_service_iface "github.com/nitric-dev/membrane/pkg/plugins/storage/azblob/iface"
	azureutils "github.com/nitric-dev/membrane/pkg/providers/azure/utils"
	"github.com/nitric-dev/membrane/pkg/utils"
)

// AzblobStorageService - Nitric membrane storage plugin implementation for Azure Storage
type AzblobStorageService struct {
	client azblob_service_iface.AzblobServiceUrlIface
}

func (a *AzblobStorageService) getBlobUrl(bucket string, key string) azblob_service_iface.AzblobBlockBlobUrlIface {
	cUrl := a.client.NewContainerURL(bucket)
	// Get a new blob for the key name
	return cUrl.NewBlockBlobURL(key)
}

func (a *AzblobStorageService) Read(bucket string, key string) ([]byte, error) {
	newErr := errors.ErrorsWithScope(
		"AzblobStorageService.Read",
		fmt.Sprintf("bucket=%s", bucket),
		fmt.Sprintf("key=%s", key),
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
		fmt.Sprintf("bucket=%s", bucket),
		fmt.Sprintf("key=%s", key),
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
		fmt.Sprintf("bucket=%s", bucket),
		fmt.Sprintf("key=%s", key),
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

const expiryBuffer = 2 * time.Minute

func tokenRefresherFromSpt(spt *adal.ServicePrincipalToken) azblob.TokenRefresher {
	return func(credential azblob.TokenCredential) time.Duration {
		if err := spt.Refresh(); err != nil {
			fmt.Println("Error refreshing token: ", err)
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
	storageAccount := utils.GetEnv("AZURE_STORAGE_ACCOUNT", "")

	if storageAccount == "" {
		return nil, fmt.Errorf("AZURE_STORAGE_ACCOUNT not configured")
	}

	spt, err := azureutils.GetServicePrincipalToken(azure.PublicCloud.ResourceIdentifiers.Storage)
	if err != nil {
		return nil, err
	}

	cTkn := azblob.NewTokenCredential(spt.Token().AccessToken, tokenRefresherFromSpt(spt))

	var accountURL *url.URL
	if accountURL, err = url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net", storageAccount)); err != nil {
		return nil, err
	}

	pipeline := azblob.NewPipeline(cTkn, azblob.PipelineOptions{})
	client := azblob.NewServiceURL(*accountURL, pipeline)

	return &AzblobStorageService{
		client: azblob_service_iface.AdaptServiceUrl(client),
	}, nil
}
