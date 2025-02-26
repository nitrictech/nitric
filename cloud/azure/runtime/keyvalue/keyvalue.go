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
	"net/http"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/data/aztables"
	"github.com/nitrictech/nitric/cloud/azure/runtime/env"
	document "github.com/nitrictech/nitric/core/pkg/decorators/keyvalue"
	grpc_errors "github.com/nitrictech/nitric/core/pkg/grpc/errors"
	kvstorepb "github.com/nitrictech/nitric/core/pkg/proto/kvstore/v1"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"
)

// AzureStorageTableKeyValueService - an Azure Storage Table implementation of the Nitric Key/Value Service
type AzureStorageTableKeyValueService struct {
	clientFactory AzureStorageClientFactory
}

// Ensure AzureStorageTableKeyValueService implements the KeyValueServer interface
var _ kvstorepb.KvStoreServer = (*AzureStorageTableKeyValueService)(nil)

type AztableEntity struct {
	aztables.Entity

	Content aztables.EDMBinary
}

// Get a value from the Azure Storage table
func (s *AzureStorageTableKeyValueService) GetValue(ctx context.Context, req *kvstorepb.KvStoreGetValueRequest) (*kvstorepb.KvStoreGetValueResponse, error) {
	newErr := grpc_errors.ErrorsWithScope("AzureStorageTableKeyValueService.GetKey")
	client, err := s.clientFactory(req.Ref.Store)
	if err != nil {
		return nil, newErr(
			codes.Internal,
			"Unable to create client",
			err,
		)
	}

	err = document.ValidateValueRef(req.Ref)
	if err != nil {
		return nil, newErr(
			codes.InvalidArgument,
			"Invalid key",
			err,
		)
	}

	response, err := client.GetEntity(ctx, req.Ref.Store, req.Ref.Key, nil)
	if err != nil {
		var respErr *azcore.ResponseError
		if errors.As(err, &respErr) {
			switch respErr.StatusCode {
			case http.StatusNotFound:
				// Handle not found error
				return nil, newErr(
					codes.NotFound,
					fmt.Sprintf("key %s not found in store %s", req.Ref.Key, req.Ref.Store),
					err,
				)
			case http.StatusForbidden:
				// Handle forbidden error
				return nil, newErr(
					codes.PermissionDenied,
					"unable to get value, this may be due to a missing permissions request in your code.",
					err,
				)
			}
		}

		return nil, newErr(
			codes.Unknown,
			"failed to call aztables:GetEntity",
			err,
		)
	}

	var entity AztableEntity
	err = json.Unmarshal(response.Value, &entity)
	if err != nil {
		return nil, newErr(
			codes.Internal,
			"Unable to convert value to AztableEntity",
			err,
		)
	}

	var structContent structpb.Struct
	err = proto.Unmarshal(entity.Content, &structContent)
	if err != nil {
		return nil, newErr(
			codes.Internal,
			"Unable to convert value to pb struct",
			err,
		)
	}

	return &kvstorepb.KvStoreGetValueResponse{
		Value: &kvstorepb.Value{
			Ref:     req.Ref,
			Content: &structContent,
		},
	}, nil
}

// Set a value in the Azure Storage table
func (s *AzureStorageTableKeyValueService) SetValue(ctx context.Context, req *kvstorepb.KvStoreSetValueRequest) (*kvstorepb.KvStoreSetValueResponse, error) {
	newErr := grpc_errors.ErrorsWithScope("AzureStorageTableKeyValueService.SetKeys")
	client, err := s.clientFactory(req.Ref.Store)
	if err != nil {
		return nil, newErr(
			codes.Internal,
			"Unable to create client",
			err,
		)
	}

	if err := document.ValidateValueRef(req.Ref); err != nil {
		return nil, newErr(
			codes.InvalidArgument,
			"invalid key",
			err,
		)
	}

	content, err := proto.Marshal(req.Content)
	if err != nil {
		return nil, newErr(
			codes.Internal,
			"Unable to marshal content",
			err,
		)
	}

	entity := AztableEntity{
		Entity: aztables.Entity{
			PartitionKey: req.Ref.Store,
			RowKey:       req.Ref.Key,
			Timestamp:    aztables.EDMDateTime(time.Now()),
		},
		Content: content,
	}

	entityJson, err := json.Marshal(entity)
	if err != nil {
		return nil, newErr(
			codes.Internal,
			"unable to convert struct to json",
			err,
		)
	}

	_, err = client.UpsertEntity(ctx, entityJson, &aztables.UpsertEntityOptions{
		UpdateMode: aztables.UpdateModeReplace,
	})
	if err != nil {
		var respErr *azcore.ResponseError
		if errors.As(err, &respErr) {
			switch respErr.StatusCode {
			case http.StatusForbidden:
				// Handle forbidden error
				return nil, newErr(
					codes.PermissionDenied,
					"unable to set value, this may be due to a missing permissions request in your code.",
					err,
				)
			}
		}

		return nil, newErr(
			codes.Unknown,
			"unable to call aztables.UpsertEntity",
			err,
		)
	}

	return &kvstorepb.KvStoreSetValueResponse{}, nil
}

// Delete a key/value pair from the Azure Storage table
func (s *AzureStorageTableKeyValueService) DeleteKey(ctx context.Context, req *kvstorepb.KvStoreDeleteKeyRequest) (*kvstorepb.KvStoreDeleteKeyResponse, error) {
	newErr := grpc_errors.ErrorsWithScope("AzureStorageTableKeyValueService.DeleteKey")
	client, err := s.clientFactory(req.Ref.Store)
	if err != nil {
		return nil, newErr(
			codes.Internal,
			"Unable to create client",
			err,
		)
	}

	if err := document.ValidateValueRef(req.Ref); err != nil {
		return nil, newErr(
			codes.InvalidArgument,
			"invalid key",
			err,
		)
	}

	_, err = client.DeleteEntity(ctx, req.Ref.Store, req.Ref.Key, nil)
	if err != nil {
		var respErr *azcore.ResponseError
		if errors.As(err, &respErr) {
			switch respErr.StatusCode {
			case http.StatusNotFound:
				// not found isn't an error for delete
				return &kvstorepb.KvStoreDeleteKeyResponse{}, nil
			case http.StatusForbidden:
				// Handle forbidden error
				return nil, newErr(
					codes.PermissionDenied,
					"unable to delete value, this may be due to a missing permissions request in your code.",
					err,
				)
			}
		}

		return nil, newErr(
			codes.Internal,
			"failed to call aztables.DeleteEntity",
			err,
		)
	}

	return &kvstorepb.KvStoreDeleteKeyResponse{}, nil
}

// Return all keys in the Azure Storage table for a key/value store
func (s *AzureStorageTableKeyValueService) ScanKeys(req *kvstorepb.KvStoreScanKeysRequest, stream kvstorepb.KvStore_ScanKeysServer) error {
	newErr := grpc_errors.ErrorsWithScope("AzureStorageTableKeyValueService.Keys")
	storeName := req.GetStore().GetName()

	if storeName == "" {
		return newErr(
			codes.InvalidArgument,
			"store name is required",
			nil,
		)
	}

	client, err := s.clientFactory(storeName)
	if err != nil {
		return newErr(
			codes.Internal,
			"unable to create client",
			err,
		)
	}

	// ge "GreaterThanOrEqual" is used for string prefix filtering (https://learn.microsoft.com/en-us/rest/api/storageservices/querying-tables-and-entities#filtering-on-string-properties)
	keyFilter := fmt.Sprintf("PartitionKey eq '%s' and RowKey ge '%s'", storeName, req.GetPrefix())

	pager := client.NewListEntitiesPager(
		&aztables.ListEntitiesOptions{
			Filter: &keyFilter,
		},
	)

	for pager.More() {
		response, err := pager.NextPage(context.TODO())
		if err != nil {
			var respErr *azcore.ResponseError
			if errors.As(err, &respErr) {
				switch respErr.StatusCode {
				case http.StatusForbidden:
					// Handle forbidden error
					return newErr(
						codes.PermissionDenied,
						"unable to list keys, this may be due to a missing permissions request in your code.",
						err,
					)
				}
			}
			return newErr(
				codes.Unknown,
				"failed to call aztables.ListEntities",
				err,
			)
		}

		for _, entityBytes := range response.Entities {
			var entity AztableEntity
			err = json.Unmarshal(entityBytes, &entity)
			if err != nil {
				return newErr(
					codes.Internal,
					"Unable to convert value to AztableEntity",
					err,
				)
			}

			if err := stream.Send(&kvstorepb.KvStoreScanKeysResponse{
				Key: entity.RowKey,
			}); err != nil {
				return newErr(
					codes.Internal,
					"failed to send response",
					err,
				)
			}
		}
	}

	return nil
}

type AzureStorageClientFactory func(tableName string) (*aztables.Client, error)

func newStorageTablesClientFactory(creds *azidentity.DefaultAzureCredential, storageAccountName string) AzureStorageClientFactory {
	return func(tableName string) (*aztables.Client, error) {
		if tableName == "" {
			return nil, fmt.Errorf("table name is required")
		}
		// Hyphens not supported by azure storage tables
		normalizedTableName := strings.Replace(tableName, "-", "", -1)

		serviceURL := fmt.Sprintf("https://%s.table.core.windows.net/%s", storageAccountName, normalizedTableName)
		return aztables.NewClient(serviceURL, creds, nil)
	}
}

// New creates a new Azure Storage Table implementation of a KeyValueServer
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
func NewWithClient(clientFactory AzureStorageClientFactory) (*AzureStorageTableKeyValueService, error) {
	return &AzureStorageTableKeyValueService{
		// storageAccountName: storageAccountName,
		// defaultCredential:  cred,
		clientFactory: clientFactory,
	}, nil
}
