package storage_plugin

import (
	"context"
	"fmt"

	"cloud.google.com/go/storage"
	"github.com/nitric-dev/membrane/plugins/sdk"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

type StoragePlugin struct {
	sdk.UnimplementedStoragePlugin
	client *storage.Client
}

func (s *StoragePlugin) Get(bucket string, key string) ([]byte, error) {
	return nil, fmt.Errorf("UNIMPLEMENTED")
}

func (s *StoragePlugin) Put(bucket string, key string, object []byte) error {
	return fmt.Errorf("UNIMPLEMENTED")
}

func New() (sdk.StoragePlugin, error) {
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

	return &StoragePlugin{
		client: client,
	}, nil
}
