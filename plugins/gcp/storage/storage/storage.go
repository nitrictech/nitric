package storage_service

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"github.com/nitric-dev/membrane/plugins/gcp/adapters"
	"github.com/nitric-dev/membrane/plugins/gcp/ifaces"
	"github.com/nitric-dev/membrane/plugins/sdk"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"io/ioutil"
)

type StorageStorageService struct {
	//sdk.UnimplementedStoragePlugin
	client    ifaces.StorageClient
	projectID string
}

func (s *StorageStorageService) getBucketByName(bucket string) (ifaces.BucketHandle, error) {
	buckets := s.client.Buckets(context.Background(), s.projectID)
	for {
		b, err := buckets.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			fmt.Println(err.Error())
			return nil, fmt.Errorf("an error occurred finding bucket: %s; %v", bucket, err)
		}
		// We'll label the buckets by their name in the nitric.yaml file and use this...
		if b.Labels["x-nitric-name"] == bucket {
			bucketHandle := s.client.Bucket(b.Name)
			return bucketHandle, nil
		}
	}
	return nil, fmt.Errorf("bucket not found")
}

/**
 * Retrieves a previously stored object from a Google Cloud Storage Bucket
 */
func (s *StorageStorageService) Get(bucket string, key string) ([]byte, error) {
	bucketHandle, err := s.getBucketByName(bucket)
	if err != nil {
		return nil, err
	}

	reader, err := bucketHandle.Object(key).NewReader(context.Background())
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

/**
 * Stores a new Item in a Google Cloud Storage Bucket
 */
func (s *StorageStorageService) Put(bucket string, key string, object []byte) error {
	bucketHandle, err := s.getBucketByName(bucket)

	if err != nil {
		return err
	}

	writer := bucketHandle.Object(key).NewWriter(context.Background())

	if _, err := writer.Write(object); err != nil {
		return err
	}

	if err := writer.Close(); err != nil {
		return err
	}

	return nil
}

/**
 * Delete an Item in a Google Cloud Storage Bucket
 */
func (s *StorageStorageService) Delete(bucket string, key string) error {
	bucketHandle, err := s.getBucketByName(bucket)

	if err != nil {
		return err
	}

	if err := bucketHandle.Object(key).Delete(context.Background()); err != nil {
		// ignore errors caused by the Object not existing.
		// This is to unify delete behavior between providers.
		if err != storage.ErrObjectNotExist {
			return err
		}
	}

	storage.ErrObjectNotExist.Error()

	return nil
}

/**
 * Creates a new Storage Plugin for use in GCP
 */
func New() (sdk.StorageService, error) {
	ctx := context.Background()

	credentials, credentialsError := google.FindDefaultCredentials(ctx, storage.ScopeReadWrite)
	if credentialsError != nil {
		return nil, fmt.Errorf("GCP credentials error: %v", credentialsError)
	}
	// Get the
	client, err := storage.NewClient(ctx, option.WithCredentials(credentials))

	if err != nil {
		return nil, fmt.Errorf("storage client error: %v", err)
	}

	return &StorageStorageService{
		client: adapters.AdaptStorageClient(client),
		projectID: credentials.ProjectID,
	}, nil
}

func NewWithClient(client ifaces.StorageClient) (sdk.StorageService, error) {
	return &StorageStorageService{
		client: client,
	}, nil
}
