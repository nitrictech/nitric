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
	"errors"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/smithy-go"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-sdk-go-v2/otelaws"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/nitrictech/nitric/cloud/aws/ifaces/dynamodbiface"
	"github.com/nitrictech/nitric/cloud/aws/runtime/env"
	"github.com/nitrictech/nitric/cloud/aws/runtime/resource"
	document "github.com/nitrictech/nitric/core/pkg/decorators/keyvalue"
	grpc_errors "github.com/nitrictech/nitric/core/pkg/grpc/errors"
	keyvaluepb "github.com/nitrictech/nitric/core/pkg/proto/keyvalue/v1"
)

const (
	AttribPk         = "_pk"
	AttribSk         = "_sk"
	deleteQueryLimit = int32(1000)
	maxBatchWrite    = 25
)

// DynamoKeyValueService - an AWS DynamoDB implementation of the Nitric Document Service
type DynamoKeyValueService struct {
	client   dynamodbiface.DynamoDBAPI
	provider resource.AwsResourceProvider
}

var _ keyvaluepb.KeyValueServer = &DynamoKeyValueService{}

func isDynamoAccessDeniedErr(err error) bool {
	var opErr *smithy.OperationError
	if errors.As(err, &opErr) {
		return opErr.Service() == "DynamoDB" && strings.Contains(opErr.Unwrap().Error(), "AccessDenied")
	}
	return false
}

// Get a document from the DynamoDB table
func (s *DynamoKeyValueService) Get(ctx context.Context, req *keyvaluepb.KeyValueGetRequest) (*keyvaluepb.KeyValueGetResponse, error) {
	newErr := grpc_errors.ErrorsWithScope("DynamoDocService.Get")

	err := document.ValidateValueRef(req.Ref)
	if err != nil {
		return nil, newErr(
			codes.InvalidArgument,
			"Invalid key",
			err,
		)
	}

	keyMap := createKeyMap(req.Ref)
	attributeMap, err := attributevalue.MarshalMap(keyMap)
	if err != nil {
		return nil, newErr(
			codes.InvalidArgument,
			"failed to marshal key",
			err,
		)
	}

	tableName, err := s.getTableName(ctx, req.Ref.Store)
	if err != nil {
		return nil, err
	}

	input := &dynamodb.GetItemInput{
		Key:       attributeMap,
		TableName: tableName,
	}

	result, err := s.client.GetItem(ctx, input)
	if err != nil {
		if isDynamoAccessDeniedErr(err) {
			return nil, newErr(
				codes.PermissionDenied,
				"unable to get document value, this may be due to a missing permissions request in your code.",
				err,
			)
		}

		return nil, newErr(
			codes.Internal,
			fmt.Sprintf("error retrieving value with key %s from store %s", req.Ref.Key, req.Ref.Store),
			err,
		)
	}

	if result.Item == nil {
		return nil, newErr(
			codes.NotFound,
			fmt.Sprintf("%v not found", req.Ref),
			err,
		)
	}

	var itemMap map[string]interface{}
	err = attributevalue.UnmarshalMap(result.Item, &itemMap)
	if err != nil {
		return nil, newErr(
			codes.Internal,
			"error unmarshalling item",
			err,
		)
	}

	delete(itemMap, AttribPk)
	delete(itemMap, AttribSk)

	documentContent, err := structpb.NewStruct(itemMap)
	if err != nil {
		return nil, newErr(
			codes.Internal,
			"error converting returned document to struct",
			err,
		)
	}

	return &keyvaluepb.KeyValueGetResponse{
		Value: &keyvaluepb.Value{
			Ref:     req.Ref,
			Content: documentContent,
		},
	}, nil
}

// Set a document in the DynamoDB table
func (s *DynamoKeyValueService) Set(ctx context.Context, req *keyvaluepb.KeyValueSetRequest) (*keyvaluepb.KeyValueSetResponse, error) {
	newErr := grpc_errors.ErrorsWithScope("DynamoDocService.Set")

	if err := document.ValidateValueRef(req.Ref); err != nil {
		return nil, newErr(
			codes.InvalidArgument,
			"invalid key",
			err,
		)
	}

	if req.Content == nil {
		return nil, newErr(
			codes.InvalidArgument,
			"document content must not be nil",
			nil,
		)
	}

	// Construct DynamoDB attribute value object
	itemMap := createItemMap(req.Content.AsMap(), req.Ref)
	itemAttributeMap, err := attributevalue.MarshalMap(itemMap)
	if err != nil {
		return nil, newErr(
			codes.NotFound,
			"failed to marshal content",
			err,
		)
	}

	tableName, err := s.getTableName(ctx, req.Ref.Store)
	if err != nil {
		return nil, newErr(
			codes.NotFound,
			"unable to find table",
			err,
		)
	}

	input := &dynamodb.PutItemInput{
		Item:      itemAttributeMap,
		TableName: tableName,
	}

	_, err = s.client.PutItem(ctx, input)
	if err != nil {
		if isDynamoAccessDeniedErr(err) {
			return nil, newErr(
				codes.PermissionDenied,
				"unable to set document value, this may be due to a missing permissions request in your code.",
				err,
			)
		}

		return nil, newErr(
			codes.Unknown,
			"unable to set document value",
			err,
		)
	}

	return &keyvaluepb.KeyValueSetResponse{}, nil
}

// Delete a document from the DynamoDB table
func (s *DynamoKeyValueService) Delete(ctx context.Context, req *keyvaluepb.KeyValueDeleteRequest) (*keyvaluepb.KeyValueDeleteResponse, error) {
	newErr := grpc_errors.ErrorsWithScope("DynamoDocService.Delete")

	if err := document.ValidateValueRef(req.Ref); err != nil {
		return nil, newErr(
			codes.InvalidArgument,
			"invalid key",
			err,
		)
	}

	keyMap := createKeyMap(req.Ref)
	attributeMap, err := attributevalue.MarshalMap(keyMap)
	if err != nil {
		return nil, newErr(
			codes.InvalidArgument,
			fmt.Sprintf("failed to marshal keys: %v", req.Ref),
			err,
		)
	}

	tableName, err := s.getTableName(ctx, req.Ref.Store)
	if err != nil {
		return nil, newErr(
			codes.NotFound,
			"unable to find table",
			err,
		)
	}

	deleteInput := &dynamodb.DeleteItemInput{
		Key:       attributeMap,
		TableName: tableName,
	}

	_, err = s.client.DeleteItem(ctx, deleteInput)
	if err != nil {
		if isDynamoAccessDeniedErr(err) {
			return nil, newErr(
				codes.PermissionDenied,
				"unable to delete document, this may be due to a missing permissions request in your code.",
				err,
			)
		}

		return nil, newErr(
			codes.Internal,
			fmt.Sprintf("error deleting %v item %v : %v", req.Ref.Store, req.Ref.Key, err),
			err,
		)
	}

	return &keyvaluepb.KeyValueDeleteResponse{}, nil
}

// New creates a new AWS DynamoDB implementation of a DocumentServiceServer
func New(provider resource.AwsResourceProvider) (*DynamoKeyValueService, error) {
	awsRegion := env.AWS_REGION.String()

	// Create a new AWS session
	cfg, sessionError := config.LoadDefaultConfig(context.TODO(), config.WithRegion(awsRegion))
	if sessionError != nil {
		return nil, fmt.Errorf("error creating new AWS session %w", sessionError)
	}

	otelaws.AppendMiddlewares(&cfg.APIOptions)

	dynamoClient := dynamodb.NewFromConfig(cfg)

	return &DynamoKeyValueService{
		client:   dynamoClient,
		provider: provider,
	}, nil
}

// NewWithClient creates a DocumentServiceServer with an given DynamoDB client instance.
//
//	Primarily used for testing
func NewWithClient(provider resource.AwsResourceProvider, client *dynamodb.Client) (*DynamoKeyValueService, error) {
	return &DynamoKeyValueService{
		provider: provider,
		client:   client,
	}, nil
}

func createKeyMap(ref *keyvaluepb.ValueRef) map[string]string {
	keyMap := make(map[string]string)

	keyMap[AttribPk] = ref.Key
	keyMap[AttribSk] = ref.Store + "#"

	return keyMap
}

func createItemMap(source map[string]interface{}, ref *keyvaluepb.ValueRef) map[string]interface{} {
	// Copy map
	newMap := make(map[string]interface{})
	for key, value := range source {
		newMap[key] = value
	}

	keyMap := createKeyMap(ref)

	// Add key attributes
	newMap[AttribPk] = keyMap[AttribPk]
	newMap[AttribSk] = keyMap[AttribSk]

	return newMap
}

func (s *DynamoKeyValueService) getTableName(ctx context.Context, store string) (*string, error) {
	tables, err := s.provider.GetResources(ctx, resource.AwsResource_Collection)
	if err != nil {
		return nil, fmt.Errorf("encountered an error retrieving the table list: %w", err)
	}

	if table, ok := tables[store]; ok {
		tableName := strings.Split(table.ARN, "/")[1]

		// split the table arn to get the name
		return aws.String(tableName), nil
	}

	return nil, fmt.Errorf("store %s does not exist", store)
}
