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
	"fmt"
	"github.com/aws/smithy-go"
	"io"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-sdk-go-v2/otelaws"

	"github.com/nitrictech/nitric/cloud/aws/ifaces/dynamodbiface"
	"github.com/nitrictech/nitric/cloud/aws/runtime/core"
	"github.com/nitrictech/nitric/core/pkg/plugins/document"
	"github.com/nitrictech/nitric/core/pkg/plugins/errors"
	"github.com/nitrictech/nitric/core/pkg/plugins/errors/codes"
	"github.com/nitrictech/nitric/core/pkg/utils"
)

const (
	AttribPk         = "_pk"
	AttribSk         = "_sk"
	deleteQueryLimit = int32(1000)
	maxBatchWrite    = 25
)

// DynamoDocService - AWS DynamoDB AWS Nitric Document service
type DynamoDocService struct {
	document.UnimplementedDocumentPlugin
	client   dynamodbiface.DynamoDBAPI
	provider core.AwsProvider
}

func isDynamoAccessDeniedErr(err error) bool {
	var opErr *smithy.OperationError
	if errors.As(err, &opErr) {
		return opErr.Service() == "DynamoDB" && strings.Contains(opErr.Unwrap().Error(), "AccessDenied")
	}
	return false
}

func (s *DynamoDocService) Get(ctx context.Context, key *document.Key) (*document.Document, error) {
	newErr := errors.ErrorsWithScope(
		"DynamoDocService.Get",
		map[string]interface{}{
			"key": key,
		},
	)

	err := document.ValidateKey(key)
	if err != nil {
		return nil, newErr(
			codes.InvalidArgument,
			"Invalid key",
			err,
		)
	}

	keyMap := createKeyMap(key)
	attributeMap, err := attributevalue.MarshalMap(keyMap)
	if err != nil {
		return nil, newErr(
			codes.InvalidArgument,
			"failed to marshal key",
			err,
		)
	}

	tableName, err := s.getTableName(ctx, *key.Collection)
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
				"unable to get document value, have you requested access to this collection?",
				err,
			)
		}

		return nil, newErr(
			codes.Internal,
			fmt.Sprintf("error retrieving key %v", key),
			err,
		)
	}

	if result.Item == nil {
		return nil, newErr(
			codes.NotFound,
			fmt.Sprintf("%v not found", key),
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

	return &document.Document{
		Key:     key,
		Content: itemMap,
	}, nil
}

func (s *DynamoDocService) Set(ctx context.Context, key *document.Key, value map[string]interface{}) error {
	newErr := errors.ErrorsWithScope(
		"DynamoDocService.Set",
		map[string]interface{}{
			"key": key,
		},
	)

	if err := document.ValidateKey(key); err != nil {
		return newErr(
			codes.InvalidArgument,
			"invalid key",
			err,
		)
	}

	if value == nil {
		return newErr(
			codes.InvalidArgument,
			"provide non-nil value",
			nil,
		)
	}

	// Construct DynamoDB attribute value object
	itemMap := createItemMap(value, key)
	itemAttributeMap, err := attributevalue.MarshalMap(itemMap)
	if err != nil {
		return fmt.Errorf("failed to marshal value")
	}

	tableName, err := s.getTableName(ctx, *key.Collection)
	if err != nil {
		return newErr(
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
			return newErr(
				codes.PermissionDenied,
				"unable to set document value, have you requested access to this collection?",
				err,
			)
		}

		return newErr(
			codes.Internal,
			"error putting item",
			err,
		)
	}

	return nil
}

func (s *DynamoDocService) Delete(ctx context.Context, key *document.Key) error {
	newErr := errors.ErrorsWithScope(
		"DynamoDocService.Delete",
		map[string]interface{}{
			"key": key,
		},
	)

	if err := document.ValidateKey(key); err != nil {
		return newErr(
			codes.InvalidArgument,
			"invalid key",
			err,
		)
	}

	keyMap := createKeyMap(key)
	attributeMap, err := attributevalue.MarshalMap(keyMap)
	if err != nil {
		return newErr(
			codes.InvalidArgument,
			fmt.Sprintf("failed to marshal keys: %v", key),
			err,
		)
	}

	tableName, err := s.getTableName(ctx, *key.Collection)
	if err != nil {
		return newErr(
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
			return newErr(
				codes.PermissionDenied,
				"unable to delete document, have you requested access to this collection?",
				err,
			)
		}

		return newErr(
			codes.Internal,
			fmt.Sprintf("error deleting %v item %v : %v", key.Collection, key.Id, err),
			err,
		)
	}

	// Delete sub collection items
	if key.Collection.Parent == nil {
		var lastEvaluatedKey map[string]types.AttributeValue
		for {
			queryInput := createDeleteQuery(tableName, key, lastEvaluatedKey)
			resp, err := s.client.Query(ctx, queryInput)
			if err != nil {
				return newErr(
					codes.Internal,
					"error performing delete in table",
					err,
				)
			}

			lastEvaluatedKey = resp.LastEvaluatedKey

			err = s.processDeleteQuery(ctx, *tableName, resp)
			if err != nil {
				return newErr(
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

	return nil
}

func (s *DynamoDocService) query(ctx context.Context, collection *document.Collection, expressions []document.QueryExpression, limit int, pagingToken map[string]string) (*document.QueryResult, error) {
	queryResult := &document.QueryResult{
		Documents: make([]document.Document, 0),
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

func (s *DynamoDocService) Query(ctx context.Context, collection *document.Collection, expressions []document.QueryExpression, limit int, pagingToken map[string]string) (*document.QueryResult, error) {
	newErr := errors.ErrorsWithScope(
		"DynamoDocService.Query",
		map[string]interface{}{
			"collection": collection,
		},
	)

	if err := document.ValidateQueryCollection(collection); err != nil {
		return nil, newErr(
			codes.InvalidArgument,
			"invalid collection",
			err,
		)
	}

	if err := document.ValidateExpressions(expressions); err != nil {
		return nil, newErr(
			codes.InvalidArgument,
			"invalid expressions",
			err,
		)
	}

	queryResult, err := s.query(ctx, collection, expressions, limit, pagingToken)
	if err != nil {
		if isDynamoAccessDeniedErr(err) {
			return nil, newErr(
				codes.PermissionDenied,
				"unable to query document values, have you requested access to this collection?",
				err,
			)
		}

		return nil, newErr(
			codes.Internal,
			"query error",
			err,
		)
	}

	remainingLimit := limit - len(queryResult.Documents)

	// If more results available, perform additional queries
	for remainingLimit > 0 &&
		(queryResult.PagingToken != nil && len(queryResult.PagingToken) > 0) {
		if res, err := s.query(ctx, collection, expressions, remainingLimit, queryResult.PagingToken); err != nil {
			return nil, newErr(
				codes.Internal,
				"query error",
				err,
			)
		} else {
			queryResult.Documents = append(queryResult.Documents, res.Documents...)
			queryResult.PagingToken = res.PagingToken
		}

		remainingLimit = limit - len(queryResult.Documents)
	}

	return queryResult, nil
}

func (s *DynamoDocService) QueryStream(ctx context.Context, collection *document.Collection, expressions []document.QueryExpression, limit int) document.DocumentIterator {
	newErr := errors.ErrorsWithScope(
		"DynamoDocService.QueryStream",
		map[string]interface{}{
			"collection": collection,
		},
	)

	colErr := document.ValidateQueryCollection(collection)
	expErr := document.ValidateExpressions(expressions)

	if colErr != nil || expErr != nil {
		// Return an error only iterator
		return func() (*document.Document, error) {
			return nil, newErr(
				codes.InvalidArgument,
				"invalid arguments",
				fmt.Errorf("collection error: %w, expression error: %w", colErr, expErr),
			)
		}
	}

	tmpLimit := limit
	var documents []document.Document
	var pagingToken map[string]string

	// Initial fetch
	res, fetchErr := s.query(ctx, collection, expressions, tmpLimit, nil)

	if fetchErr != nil {
		// Return an error only iterator if the initial fetch failed
		return func() (*document.Document, error) {
			if isDynamoAccessDeniedErr(fetchErr) {
				return nil, newErr(
					codes.PermissionDenied,
					"unable to query document values, have you requested access to this collection?",
					fetchErr,
				)
			}

			return nil, newErr(
				codes.Internal,
				"query error",
				fetchErr,
			)
		}
	}

	documents = res.Documents
	pagingToken = res.PagingToken

	return func() (*document.Document, error) {
		// check the iteration state
		if tmpLimit <= 0 && limit > 0 {
			// we've reached the limit of reading
			return nil, io.EOF
		} else if pagingToken != nil && len(documents) == 0 {
			// we've run out of documents and have more pages to read
			res, fetchErr = s.query(ctx, collection, expressions, tmpLimit, pagingToken)
			documents = res.Documents
			pagingToken = res.PagingToken
		} else if pagingToken == nil && len(documents) == 0 {
			// we're all out of documents and pages before hitting the limit
			return nil, io.EOF
		}

		// We received an error fetching the docs
		if fetchErr != nil {
			return nil, newErr(
				codes.Internal,
				"query error",
				fetchErr,
			)
		}

		if len(documents) == 0 {
			return nil, io.EOF
		}

		// pop the first element
		var doc document.Document
		doc, documents = documents[0], documents[1:]
		tmpLimit = tmpLimit - 1

		return &doc, nil
	}
}

// New - Create a new DynamoDB key value plugin implementation
func New(provider core.AwsProvider) (document.DocumentService, error) {
	awsRegion := utils.GetEnv("AWS_REGION", "us-east-1")

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

// NewWithClient - Mainly used for testing
func NewWithClient(provider core.AwsProvider, client *dynamodb.Client) (document.DocumentService, error) {
	return &DynamoDocService{
		provider: provider,
		client:   client,
	}, nil
}

// Private Functions ----------------------------------------------------------

func createKeyMap(key *document.Key) map[string]string {
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

func createItemMap(source map[string]interface{}, key *document.Key) map[string]interface{} {
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
	collection *document.Collection,
	expressions []document.QueryExpression,
	limit int,
	pagingToken map[string]string,
) (*document.QueryResult, error)

func (s *DynamoDocService) performQuery(
	ctx context.Context,
	collection *document.Collection,
	expressions []document.QueryExpression,
	limit int,
	pagingToken map[string]string,
) (*document.QueryResult, error) {
	if collection.Parent == nil {
		// Should never occur
		return nil, fmt.Errorf("cannot perform query without partion key defined")
	}

	// Sort expressions to help map where "A >= %1 AND A <= %2" to DynamoDB expression "A BETWEEN %1 AND %2"
	sort.Sort(document.ExpsSort(expressions))

	tableName, err := s.getTableName(ctx, *collection)
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
		limit64 := int32(limit)
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
	collection *document.Collection,
	expressions []document.QueryExpression,
	limit int,
	pagingToken map[string]string,
) (*document.QueryResult, error) {
	// Sort expressions to help map where "A >= %1 AND A <= %2" to DynamoDB expression "A BETWEEN %1 AND %2"
	sort.Sort(document.ExpsSort(expressions))

	tableName, err := s.getTableName(ctx, *collection)
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
		limit32 := int32(limit)
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

func marshalQueryResult(collection *document.Collection, items []map[string]types.AttributeValue, lastEvaluatedKey map[string]types.AttributeValue) (*document.QueryResult, error) {
	// Unmarshal Dynamo response items
	var pTkn map[string]string = nil
	var valueMaps []map[string]interface{}
	if err := attributevalue.UnmarshalListOfMaps(items, &valueMaps); err != nil {
		return nil, fmt.Errorf("error unmarshalling query response: %w", err)
	}

	docs := make([]document.Document, 0, len(valueMaps))

	// Strip keys & append results
	for _, m := range valueMaps {
		// Retrieve the original ID on the result
		var id string
		var c *document.Collection
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
			c = &document.Collection{
				Name: collection.Name,
				Parent: &document.Key{
					Collection: &document.Collection{
						Name: collection.Parent.Collection.Name,
					},
					Id: pk,
				},
			}
		}

		// Split out sort key value
		delete(m, AttribPk)
		delete(m, AttribSk)

		sdkDoc := document.Document{
			Key: &document.Key{
				Collection: c,
				Id:         id,
			},
			Content: m,
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

	return &document.QueryResult{
		Documents:   docs,
		PagingToken: pTkn,
	}, nil
}

func createFilterExpression(expressions []document.QueryExpression) string {
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

func isBetweenStart(index int, exps []document.QueryExpression) bool {
	if index < (len(exps) - 1) {
		if exps[index].Operand == exps[index+1].Operand &&
			exps[index].Operator == ">=" &&
			exps[index+1].Operator == "<=" {
			return true
		}
	}
	return false
}

func isBetweenEnd(index int, exps []document.QueryExpression) bool {
	if index > 0 && index < len(exps) {
		if exps[index-1].Operand == exps[index].Operand &&
			exps[index-1].Operator == ">=" &&
			exps[index].Operator == "<=" {
			return true
		}
	}
	return false
}

func (s *DynamoDocService) getTableName(ctx context.Context, collection document.Collection) (*string, error) {
	tables, err := s.provider.GetResources(ctx, core.AwsResource_Collection)
	if err != nil {
		return nil, fmt.Errorf("encountered an error retrieving the table list: %w", err)
	}

	coll := collection
	for coll.Parent != nil {
		coll = *coll.Parent.Collection
	}

	if table, ok := tables[coll.Name]; ok {
		tableName := strings.Split(table, "/")[1]

		// split the table arn to get the name
		return aws.String(tableName), nil
	}

	return nil, fmt.Errorf("collection %s does not exist", coll.Name)
}

func createDeleteQuery(table *string, key *document.Key, startKey map[string]types.AttributeValue) *dynamodb.QueryInput {
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
