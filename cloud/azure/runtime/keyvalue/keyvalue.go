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

package keyvalue

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/data/aztables"
	"github.com/nitrictech/nitric/cloud/aws/runtime/resource"
	"github.com/nitrictech/nitric/cloud/azure/runtime/env"
	document "github.com/nitrictech/nitric/core/pkg/decorators/keyvalue"
	grpc_errors "github.com/nitrictech/nitric/core/pkg/grpc/errors"
	keyvaluepb "github.com/nitrictech/nitric/core/pkg/proto/keyvalue/v1"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/structpb"
)

// DynamoKeyValueService - an AWS DynamoDB implementation of the Nitric Document Service
type AzureStorageTableKeyValueService struct {
	clientFactory AzureStorageClientFactory
}

var _ keyvaluepb.KeyValueServer = &AzureStorageTableKeyValueService{}

type AztableEntity struct {
	aztables.Entity

	Content *structpb.Struct
}

// Get a document from the DynamoDB table
func (s *AzureStorageTableKeyValueService) Get(ctx context.Context, req *keyvaluepb.KeyValueGetRequest) (*keyvaluepb.KeyValueGetResponse, error) {
	newErr := grpc_errors.ErrorsWithScope("AzureStorageTableKeyValueService.Get")
	client, err := s.clientFactory(req.Key.Store)
	if err != nil {
		return nil, newErr(
			codes.Internal,
			"Unable to create client",
			err,
		)
	}

	err = document.ValidateKey(req.Key)
	if err != nil {
		return nil, newErr(
			codes.InvalidArgument,
			"Invalid key",
			err,
		)
	}

	response, err := client.GetEntity(ctx, req.Key.Store, req.Key.Key, nil)
	if err != nil {
		return nil, newErr(
			codes.InvalidArgument,
			"failed to call aztables:GetEntity",
			err,
		)
	}

	var entity AztableEntity
	err = json.Unmarshal(response.Value, &entity)
	if err != nil {
		return nil, newErr(
			codes.Internal,
			"Unable to convert value to pb struct",
			err,
		)
	}

	return &keyvaluepb.KeyValueGetResponse{
		Value: &keyvaluepb.Value{
			Key:     req.Key,
			Content: entity.Content,
		},
	}, nil
}

// Set a document in the DynamoDB table
func (s *AzureStorageTableKeyValueService) Set(ctx context.Context, req *keyvaluepb.KeyValueSetRequest) (*keyvaluepb.KeyValueSetResponse, error) {
	newErr := grpc_errors.ErrorsWithScope("AzureStorageTableKeyValueService.Set")
	client, err := s.clientFactory(req.Key.Store)
	if err != nil {
		return nil, newErr(
			codes.Internal,
			"Unable to create client",
			err,
		)
	}

	if err := document.ValidateKey(req.Key); err != nil {
		return nil, newErr(
			codes.InvalidArgument,
			"invalid key",
			err,
		)
	}

	entity := AztableEntity{
		Entity: aztables.Entity{
			PartitionKey: req.Key.Store,
			RowKey:       req.Key.Key,
		},
		Content: req.Content,
	}

	entityJson, err := json.Marshal(entity)
	if err != nil {
		return nil, newErr(
			codes.Internal,
			"Unable to struct to json",
			err,
		)
	}

	_, err = client.UpsertEntity(ctx, entityJson, nil)
	if err != nil {
		return nil, newErr(
			codes.Internal,
			"Unable to call aztables.UpsertEntity",
			err,
		)
	}

	return &keyvaluepb.KeyValueSetResponse{}, nil
}

// Delete a document from the DynamoDB table
func (s *AzureStorageTableKeyValueService) Delete(ctx context.Context, req *keyvaluepb.KeyValueDeleteRequest) (*keyvaluepb.KeyValueDeleteResponse, error) {
	newErr := grpc_errors.ErrorsWithScope("AzureStorageTableKeyValueService.Delete")
	client, err := s.clientFactory(req.Key.Store)
	if err != nil {
		return nil, newErr(
			codes.Internal,
			"Unable to create client",
			err,
		)
	}

	if err := document.ValidateKey(req.Key); err != nil {
		return nil, newErr(
			codes.InvalidArgument,
			"invalid key",
			err,
		)
	}

	_, err = client.DeleteEntity(ctx, req.Key.Store, req.Key.Key, nil)
	if err != nil {
		return nil, newErr(
			codes.Internal,
			"failed to call aztables.DeleteEntity",
			err,
		)
	}

	return &keyvaluepb.KeyValueDeleteResponse{}, nil
}

type AzureStorageClientFactory func(tableName string) (*aztables.Client, error)

func newStorageTablesClientFactory(creds *azidentity.DefaultAzureCredential, storageAccountName string) AzureStorageClientFactory {
	return func(tableName string) (*aztables.Client, error) {
		serviceURL := fmt.Sprintf("https://%s.table.core.windows.net/%s", storageAccountName, tableName)
		return aztables.NewClient(serviceURL, creds, nil)
	}
}

// New creates a new AWS DynamoDB implementation of a DocumentServiceServer
func New() (*AzureStorageTableKeyValueService, error) {
	storageAccountName := env.AZURE_STORAGE_ACCOUNT_NAME.String()
	if storageAccountName == "" {
		return nil, fmt.Errorf("failed to determine Azure Storage Account name, environment variable not set")
	}

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to locate default azure credential")
	}

	return &AzureStorageTableKeyValueService{
		clientFactory: newStorageTablesClientFactory(cred, storageAccountName),
	}, nil
}

// NewWithClient creates a DocumentServiceServer with an given DynamoDB client instance.
//
//	Primarily used for testing
func NewWithClient(provider resource.AwsResourceProvider, clientFactory AzureStorageClientFactory) (*AzureStorageTableKeyValueService, error) {
	return &AzureStorageTableKeyValueService{
		// storageAccountName: storageAccountName,
		// defaultCredential:  cred,
		clientFactory: clientFactory,
	}, nil
}
