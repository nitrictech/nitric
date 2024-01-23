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

package documents

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/smithy-go"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-sdk-go-v2/otelaws"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/nitrictech/nitric/cloud/aws/ifaces/dynamodbiface"
	"github.com/nitrictech/nitric/cloud/aws/runtime/env"
	"github.com/nitrictech/nitric/cloud/aws/runtime/resource"
	document "github.com/nitrictech/nitric/core/pkg/decorators/documents"
	grpc_errors "github.com/nitrictech/nitric/core/pkg/grpc/errors"
	documentpb "github.com/nitrictech/nitric/core/pkg/proto/documents/v1"
)

const (
	AttribPk         = "_pk"
	AttribSk         = "_sk"
	deleteQueryLimit = int32(1000)
	maxBatchWrite    = 25
)

// DynamoDocService - an AWS DynamoDB implementation of the Nitric Document Service
type DynamoDocService struct {
	client   dynamodbiface.DynamoDBAPI
	provider resource.AwsResourceProvider
}

var _ documentpb.DocumentsServer = &DynamoDocService{}

func isDynamoAccessDeniedErr(err error) bool {
	var opErr *smithy.OperationError
	if errors.As(err, &opErr) {
		return opErr.Service() == "DynamoDB" && strings.Contains(opErr.Unwrap().Error(), "AccessDenied")
	}
	return false
}

// Get a document from the DynamoDB table
func (s *DynamoDocService) Get(ctx context.Context, req *documentpb.DocumentGetRequest) (*documentpb.DocumentGetResponse, error) {
	newErr := grpc_errors.ErrorsWithScope("DynamoDocService.Get")

	err := document.ValidateKey(req.Key)
	if err != nil {
		return nil, newErr(
			codes.InvalidArgument,
			"Invalid key",
			err,
		)
	}

	keyMap := createKeyMap(req.Key)
	attributeMap, err := attributevalue.MarshalMap(keyMap)
	if err != nil {
		return nil, newErr(
			codes.InvalidArgument,
			"failed to marshal key",
			err,
		)
	}

	tableName, err := s.getTableName(ctx, req.Key.Collection)
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
			fmt.Sprintf("error retrieving key %v", req.Key),
			err,
		)
	}

	if result.Item == nil {
		return nil, newErr(
			codes.NotFound,
			fmt.Sprintf("%v not found", req.Key),
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

	return &documentpb.DocumentGetResponse{
		Document: &documentpb.Document{
			Key:     req.Key,
			Content: documentContent,
		},
	}, nil
}

// Set a document in the DynamoDB table
func (s *DynamoDocService) Set(ctx context.Context, req *documentpb.DocumentSetRequest) (*documentpb.DocumentSetResponse, error) {
	newErr := grpc_errors.ErrorsWithScope("DynamoDocService.Set")

	if err := document.ValidateKey(req.Key); err != nil {
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
	itemMap := createItemMap(req.Content.AsMap(), req.Key)
	itemAttributeMap, err := attributevalue.MarshalMap(itemMap)
	if err != nil {
		return nil, newErr(
			codes.NotFound,
			"failed to marshal content",
			err,
		)
	}

	tableName, err := s.getTableName(ctx, req.Key.Collection)
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

	return &documentpb.DocumentSetResponse{}, nil
}

// Delete a document from the DynamoDB table
func (s *DynamoDocService) Delete(ctx context.Context, req *documentpb.DocumentDeleteRequest) (*documentpb.DocumentDeleteResponse, error) {
	newErr := grpc_errors.ErrorsWithScope("DynamoDocService.Delete")

	if err := document.ValidateKey(req.Key); err != nil {
		return nil, newErr(
			codes.InvalidArgument,
			"invalid key",
			err,
		)
	}

	keyMap := createKeyMap(req.Key)
	attributeMap, err := attributevalue.MarshalMap(keyMap)
	if err != nil {
		return nil, newErr(
			codes.InvalidArgument,
			fmt.Sprintf("failed to marshal keys: %v", req.Key),
			err,
		)
	}

	tableName, err := s.getTableName(ctx, req.Key.Collection)
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
			fmt.Sprintf("error deleting %v item %v : %v", req.Key.Collection, req.Key.Id, err),
			err,
		)
	}

	// Delete sub collection items
	if req.Key.Collection.Parent == nil {
		var lastEvaluatedKey map[string]types.AttributeValue
		for {
			queryInput := createDeleteQuery(tableName, req.Key, lastEvaluatedKey)
			resp, err := s.client.Query(ctx, queryInput)
			if err != nil {
				return nil, newErr(
					codes.Internal,
					"error performing delete in table",
					err,
				)
			}

			lastEvaluatedKey = resp.LastEvaluatedKey

			err = s.processDeleteQuery(ctx, *tableName, resp)
			if err != nil {
				return nil, newErr(
					codes.Internal,
					"error performing delete",
					err,
				)
			}

			if len(lastEvaluatedKey) == 0 {
				break
			}
		}
	}

	return &documentpb.DocumentDeleteResponse{}, nil
}

func (s *DynamoDocService) query(ctx context.Context, collection *documentpb.Collection, expressions []*documentpb.Expression, limit int32, pagingToken map[string]string) (*documentpb.DocumentQueryResponse, error) {
	queryResult := &documentpb.DocumentQueryResponse{
		Documents: make([]*documentpb.Document, 0),
	}

	var resFunc resultRetriever = s.performQuery
	if collection.Parent == nil || collection.Parent.Id == "" {
		resFunc = s.performScan
	}

	if res, err := resFunc(ctx, collection, expressions, limit, pagingToken); err != nil {
		return nil, err
	} else {
		queryResult.Documents = append(queryResult.Documents, res.Documents...)
		queryResult.PagingToken = res.PagingToken
	}

	return queryResult, nil
}

// Query documents from the DynamoDB table with pagination
func (s *DynamoDocService) Query(ctx context.Context, req *documentpb.DocumentQueryRequest) (*documentpb.DocumentQueryResponse, error) {
	newErr := grpc_errors.ErrorsWithScope("DynamoDocService.Query")

	if err := document.ValidateQueryCollection(req.Collection); err != nil {
		return nil, newErr(
			codes.InvalidArgument,
			"invalid collection",
			err,
		)
	}

	if err := document.ValidateExpressions(req.Expressions); err != nil {
		return nil, newErr(
			codes.InvalidArgument,
			"invalid expressions",
			err,
		)
	}

	queryResult, err := s.query(ctx, req.Collection, req.Expressions, req.Limit, req.PagingToken)
	if err != nil {
		if isDynamoAccessDeniedErr(err) {
			return nil, newErr(
				codes.PermissionDenied,
				"unable to query document values, this may be due to a missing permissions request in your code.",
				err,
			)
		}

		return nil, newErr(
			codes.Internal,
			"query error",
			err,
		)
	}

	remainingLimit := req.Limit - int32(len(queryResult.Documents))

	// If more results available, perform additional queries
	for remainingLimit > 0 &&
		(queryResult.PagingToken != nil && len(queryResult.PagingToken) > 0) {
		if res, err := s.query(ctx, req.Collection, req.Expressions, remainingLimit, queryResult.PagingToken); err != nil {
			return nil, newErr(
				codes.Internal,
				"query error",
				err,
			)
		} else {
			queryResult.Documents = append(queryResult.Documents, res.Documents...)
			queryResult.PagingToken = res.PagingToken
		}

		remainingLimit = req.Limit - int32(len(queryResult.Documents))
	}

	return queryResult, nil
}

// QuerySteam queries documents from the DynamoDB table as a stream
func (s *DynamoDocService) QueryStream(req *documentpb.DocumentQueryStreamRequest, srv documentpb.Documents_QueryStreamServer) error {
	newErr := grpc_errors.ErrorsWithScope("DynamoDocService.QueryStream")

	colErr := document.ValidateQueryCollection(req.Collection)
	if colErr != nil {
		return newErr(
			codes.InvalidArgument,
			"invalid arguments",
			fmt.Errorf("collection error: %w", colErr),
		)
	}

	expErr := document.ValidateExpressions(req.Expressions)
	if expErr != nil {
		return newErr(
			codes.InvalidArgument,
			"invalid arguments",
			fmt.Errorf("expression error: %w", expErr),
		)
	}

	var pagingToken map[string]string
	numReturned := int32(0)

	for numReturned < req.Limit {
		res, fetchErr := s.query(srv.Context(), req.Collection, req.Expressions, req.Limit-numReturned, pagingToken)
		pagingToken = res.PagingToken

		if fetchErr != nil {
			return newErr(
				codes.Internal,
				"query error",
				fetchErr,
			)
		}

		// no more results to return
		if len(res.Documents) == 0 {
			return nil
		}

		for _, doc := range res.Documents {
			if err := srv.Send(&documentpb.DocumentQueryStreamResponse{
				Document: doc,
			}); err != nil {
				return newErr(
					codes.Internal,
					"error returning document",
					err,
				)
			}

			numReturned++
			if numReturned >= req.Limit {
				break
			}
		}
	}

	return nil
}

// New creates a new AWS DynamoDB implementation of a DocumentServiceServer
func New(provider resource.AwsResourceProvider) (*DynamoDocService, error) {
	awsRegion := env.AWS_REGION.String()

	// Create a new AWS session
	cfg, sessionError := config.LoadDefaultConfig(context.TODO(), config.WithRegion(awsRegion))
	if sessionError != nil {
		return nil, fmt.Errorf("error creating new AWS session %w", sessionError)
	}

	otelaws.AppendMiddlewares(&cfg.APIOptions)

	dynamoClient := dynamodb.NewFromConfig(cfg)

	return &DynamoDocService{
		client:   dynamoClient,
		provider: provider,
	}, nil
}

// NewWithClient creates a DocumentServiceServer with an given DynamoDB client instance.
//
//	Primarily used for testing
func NewWithClient(provider resource.AwsResourceProvider, client *dynamodb.Client) (*DynamoDocService, error) {
	return &DynamoDocService{
		provider: provider,
		client:   client,
	}, nil
}

func createKeyMap(key *documentpb.Key) map[string]string {
	keyMap := make(map[string]string)

	parentKey := key.Collection.Parent

	if parentKey == nil {
		keyMap[AttribPk] = key.Id
		keyMap[AttribSk] = key.Collection.Name + "#"
	} else {
		keyMap[AttribPk] = parentKey.Id
		keyMap[AttribSk] = key.Collection.Name + "#" + key.Id
	}

	return keyMap
}

func createItemMap(source map[string]interface{}, key *documentpb.Key) map[string]interface{} {
	// Copy map
	newMap := make(map[string]interface{})
	for key, value := range source {
		newMap[key] = value
	}

	keyMap := createKeyMap(key)

	// Add key attributes
	newMap[AttribPk] = keyMap[AttribPk]
	newMap[AttribSk] = keyMap[AttribSk]

	return newMap
}

type resultRetriever = func(
	ctx context.Context,
	collection *documentpb.Collection,
	expressions []*documentpb.Expression,
	limit int32,
	pagingToken map[string]string,
) (*documentpb.DocumentQueryResponse, error)

func (s *DynamoDocService) performQuery(
	ctx context.Context,
	collection *documentpb.Collection,
	expressions []*documentpb.Expression,
	limit int32,
	pagingToken map[string]string,
) (*documentpb.DocumentQueryResponse, error) {
	if collection.Parent == nil {
		// Should never occur
		return nil, fmt.Errorf("cannot perform query without partition key defined")
	}

	// Sort expressions to help map where "A >= %1 AND A <= %2" to DynamoDB expression "A BETWEEN %1 AND %2"
	sort.Sort(document.ExpsSort(expressions))

	tableName, err := s.getTableName(ctx, collection)
	if err != nil {
		return nil, err
	}

	input := &dynamodb.QueryInput{
		TableName: tableName,
	}

	// Configure KeyConditionExpression
	keyExp := "#pk = :pk AND begins_with(#sk, :sk)"
	input.KeyConditionExpression = aws.String(keyExp)

	// Configure FilterExpression
	filterExp := createFilterExpression(expressions)
	if filterExp != "" {
		input.FilterExpression = aws.String(filterExp)
	}

	// Configure ExpressionAttributeNames
	input.ExpressionAttributeNames = make(map[string]string)
	input.ExpressionAttributeNames["#pk"] = "_pk"
	input.ExpressionAttributeNames["#sk"] = "_sk"
	for _, exp := range expressions {
		input.ExpressionAttributeNames["#"+exp.Operand] = exp.Operand
	}

	// Configure ExpressionAttributeValues
	input.ExpressionAttributeValues = make(map[string]types.AttributeValue)
	input.ExpressionAttributeValues[":pk"] = &types.AttributeValueMemberS{
		Value: collection.Parent.Id,
	}
	input.ExpressionAttributeValues[":sk"] = &types.AttributeValueMemberS{
		Value: collection.Name + "#",
	}
	for i, exp := range expressions {
		expKey := fmt.Sprintf(":%v%v", exp.Operand, i)
		valAttrib, err := attributevalue.Marshal(exp.Value)
		if err != nil {
			return nil, fmt.Errorf("error marshalling %v: %v", exp.Operand, exp.Value)
		}
		input.ExpressionAttributeValues[expKey] = valAttrib
	}

	// Configure fetch Limit
	if limit > 0 {
		limit64 := limit
		input.Limit = &(limit64)

		if len(pagingToken) > 0 {
			startKey, err := attributevalue.MarshalMap(pagingToken)
			if err != nil {
				return nil, fmt.Errorf("error performing query %v: %w", input, err)
			}
			input.ExclusiveStartKey = startKey
		}
	}

	// Perform query
	resp, err := s.client.Query(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("error performing query %v: %w", input, err)
	}

	return marshalQueryResult(collection, resp.Items, resp.LastEvaluatedKey)
}

func (s *DynamoDocService) performScan(
	ctx context.Context,
	collection *documentpb.Collection,
	expressions []*documentpb.Expression,
	limit int32,
	pagingToken map[string]string,
) (*documentpb.DocumentQueryResponse, error) {
	// Sort expressions to help map where "A >= %1 AND A <= %2" to DynamoDB expression "A BETWEEN %1 AND %2"
	sort.Sort(document.ExpsSort(expressions))

	tableName, err := s.getTableName(ctx, collection)
	if err != nil {
		return nil, err
	}

	input := &dynamodb.ScanInput{
		TableName: tableName,
	}

	// Filter on SK collection name or sub-collection name
	filterExp := "#sk = :sk"
	if collection.Parent != nil {
		filterExp = "begins_with(#sk, :sk)"
	}

	expFilters := createFilterExpression(expressions)
	if expFilters != "" {
		filterExp += " AND " + expFilters
	}

	// Configure FilterExpression
	input.FilterExpression = aws.String(filterExp)

	// Configure ExpressionAttributeNames
	input.ExpressionAttributeNames = make(map[string]string)
	input.ExpressionAttributeNames["#sk"] = "_sk"

	for _, exp := range expressions {
		input.ExpressionAttributeNames["#"+exp.Operand] = exp.Operand
	}

	// Configure ExpressionAttributeValues
	input.ExpressionAttributeValues = make(map[string]types.AttributeValue)
	keyAttrib := &types.AttributeValueMemberS{Value: collection.Name + "#"}

	input.ExpressionAttributeValues[":sk"] = keyAttrib
	for i, exp := range expressions {
		expKey := fmt.Sprintf(":%v%v", exp.Operand, i)
		valAttrib, err := attributevalue.Marshal(exp.Value)
		if err != nil {
			return nil, fmt.Errorf("error marshalling %v: %v", exp.Operand, exp.Value)
		}
		input.ExpressionAttributeValues[expKey] = valAttrib
	}

	// Configure fetch Limit
	if limit > 0 {
		// Account for parent record in fetch limit
		limit32 := limit
		input.Limit = &(limit32)

		if len(pagingToken) > 0 {
			startKey, err := attributevalue.MarshalMap(pagingToken)
			if err != nil {
				return nil, fmt.Errorf("error performing scan %v: %w", input, err)
			}
			input.ExclusiveStartKey = startKey
		}
	}

	resp, err := s.client.Scan(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("error performing scan %v: %w", input, err)
	}

	return marshalQueryResult(collection, resp.Items, resp.LastEvaluatedKey)
}

func marshalQueryResult(collection *documentpb.Collection, items []map[string]types.AttributeValue, lastEvaluatedKey map[string]types.AttributeValue) (*documentpb.DocumentQueryResponse, error) {
	// Unmarshal Dynamo response items
	var pTkn map[string]string = nil
	var valueMaps []map[string]interface{}
	if err := attributevalue.UnmarshalListOfMaps(items, &valueMaps); err != nil {
		return nil, fmt.Errorf("error unmarshalling query response: %w", err)
	}

	docs := make([]*documentpb.Document, 0, len(valueMaps))

	// Strip keys & append results
	for _, m := range valueMaps {
		// Retrieve the original ID on the result
		var id string
		var c *documentpb.Collection
		if collection.Parent == nil {
			// We know this is a root document so its key will be located in PK
			pk, _ := m[AttribPk].(string)
			id = pk
			c = collection
		} else {
			// We know this is a child document so its key will be located in the SK
			pk, _ := m[AttribPk].(string)
			sk, _ := m[AttribSk].(string)
			idStr := strings.Split(sk, "#")
			id = idStr[len(idStr)-1]
			c = &documentpb.Collection{
				Name: collection.Name,
				Parent: &documentpb.Key{
					Collection: &documentpb.Collection{
						Name: collection.Parent.Collection.Name,
					},
					Id: pk,
				},
			}
		}

		// Split out sort key value
		delete(m, AttribPk)
		delete(m, AttribSk)

		structContent, err := structpb.NewStruct(m)
		if err != nil {
			return nil, err
		}

		sdkDoc := &documentpb.Document{
			Key: &documentpb.Key{
				Collection: c,
				Id:         id,
			},
			Content: structContent,
		}
		docs = append(docs, sdkDoc)
	}

	// Unmarshal lastEvalutedKey
	var resultPagingToken map[string]string
	if len(lastEvaluatedKey) > 0 {
		if err := attributevalue.UnmarshalMap(lastEvaluatedKey, &resultPagingToken); err != nil {
			return nil, fmt.Errorf("error unmarshalling query lastEvaluatedKey: %w", err)
		}
		pTkn = resultPagingToken
	}

	return &documentpb.DocumentQueryResponse{
		Documents:   docs,
		PagingToken: pTkn,
	}, nil
}

func createFilterExpression(expressions []*documentpb.Expression) string {
	keyExp := ""
	for i, exp := range expressions {
		if keyExp != "" {
			keyExp += " AND "
		}

		if isBetweenStart(i, expressions) {
			// #{exp.operand} BETWEEN :{exp.operand}{exp.index})
			keyExp += fmt.Sprintf("#%v BETWEEN :%s%d", exp.Operand, exp.Operand, i)
		} else if isBetweenEnd(i, expressions) {
			// AND :{exp.operand}{exp.index})
			keyExp += fmt.Sprintf(":%s%d", exp.Operand, i)
		} else if exp.Operator == "startsWith" {
			// begins_with(#{exp.operand}, :{exp.operand}{exp.index})
			keyExp += fmt.Sprintf("begins_with(#%s, :%s%d)", exp.Operand, exp.Operand, i)
		} else if exp.Operator == "==" {
			// #{exp.operand} = :{exp.operand}{exp.index}
			keyExp += fmt.Sprintf("#%s = :%s%d", exp.Operand, exp.Operand, i)
		} else {
			// #{exp.operand} {exp.operator} :{exp.operand}{exp.index}
			keyExp += fmt.Sprintf("#%s %s :%s%d", exp.Operand, exp.Operator, exp.Operand, i)
		}
	}

	return keyExp
}

func isBetweenStart(index int, exps []*documentpb.Expression) bool {
	if index < (len(exps) - 1) {
		if exps[index].Operand == exps[index+1].Operand &&
			exps[index].Operator == ">=" &&
			exps[index+1].Operator == "<=" {
			return true
		}
	}
	return false
}

func isBetweenEnd(index int, exps []*documentpb.Expression) bool {
	if index > 0 && index < len(exps) {
		if exps[index-1].Operand == exps[index].Operand &&
			exps[index-1].Operator == ">=" &&
			exps[index].Operator == "<=" {
			return true
		}
	}
	return false
}

func (s *DynamoDocService) getTableName(ctx context.Context, collection *documentpb.Collection) (*string, error) {
	tables, err := s.provider.GetResources(ctx, resource.AwsResource_Collection)
	if err != nil {
		return nil, fmt.Errorf("encountered an error retrieving the table list: %w", err)
	}

	coll := collection
	for coll.Parent != nil {
		coll = coll.Parent.Collection
	}

	if table, ok := tables[coll.Name]; ok {
		tableName := strings.Split(table.ARN, "/")[1]

		// split the table arn to get the name
		return aws.String(tableName), nil
	}

	return nil, fmt.Errorf("collection %s does not exist", coll.Name)
}

func createDeleteQuery(table *string, key *documentpb.Key, startKey map[string]types.AttributeValue) *dynamodb.QueryInput {
	limit := deleteQueryLimit

	return &dynamodb.QueryInput{
		TableName:              table,
		Limit:                  &(limit),
		Select:                 types.SelectSpecificAttributes,
		ProjectionExpression:   aws.String("#pk, #sk"),
		KeyConditionExpression: aws.String("#pk = :pk"),
		ExpressionAttributeNames: map[string]string{
			"#pk": "_pk",
			"#sk": "_sk",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{
				Value: key.Id,
			},
		},
		ExclusiveStartKey: startKey,
	}
}

func (s *DynamoDocService) processDeleteQuery(ctx context.Context, table string, resp *dynamodb.QueryOutput) error {
	itemIndex := 0
	for itemIndex < len(resp.Items) {
		batchInput := &dynamodb.BatchWriteItemInput{}
		batchInput.RequestItems = make(map[string][]types.WriteRequest)
		writeRequests := make([]types.WriteRequest, 0, maxBatchWrite)

		batchCount := 0
		for batchCount < maxBatchWrite && itemIndex < len(resp.Items) {
			item := resp.Items[itemIndex]
			itemIndex += 1

			writeRequest := types.WriteRequest{}

			writeRequest.DeleteRequest = &types.DeleteRequest{
				Key: map[string]types.AttributeValue{
					AttribPk: item[AttribPk],
					AttribSk: item[AttribSk],
				},
			}
			writeRequests = append(writeRequests, writeRequest)

			batchCount += 1
		}

		batchInput.RequestItems = make(map[string][]types.WriteRequest)
		batchInput.RequestItems[table] = writeRequests

		_, err := s.client.BatchWriteItem(ctx, batchInput)
		if err != nil {
			return err
		}
	}

	return nil
}
