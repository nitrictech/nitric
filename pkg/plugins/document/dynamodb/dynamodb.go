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

package dynamodb_service

import (
	"fmt"
	"sort"
	"strings"

	"github.com/nitric-dev/membrane/pkg/plugins/document"
	"github.com/nitric-dev/membrane/pkg/utils"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

const ATTRIB_PK = "_pk"
const ATTRIB_SK = "_sk"

// DynamoDocService - AWS DynamoDB AWS Nitric Document service
type DynamoDocService struct {
	document.UnimplementedDocumentPlugin
	client dynamodbiface.DynamoDBAPI
}

func (s *DynamoDocService) Get(key *document.Key) (*document.Document, error) {
	err := document.ValidateKey(key)
	if err != nil {
		return nil, err
	}

	keyMap := createKeyMap(key)
	attributeMap, err := dynamodbattribute.MarshalMap(keyMap)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal key: %v", key)
	}

	input := &dynamodb.GetItemInput{
		Key:       attributeMap,
		TableName: getTableName(*key.Collection),
	}

	result, err := s.client.GetItem(input)
	if err != nil {
		return nil, fmt.Errorf("error getting %v : %v", key, err)
	}

	if result.Item == nil {
		return nil, fmt.Errorf("%v value not found", key)
	}

	var itemMap map[string]interface{}
	err = dynamodbattribute.UnmarshalMap(result.Item, &itemMap)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling item: %v", err)
	}

	delete(itemMap, ATTRIB_PK)
	delete(itemMap, ATTRIB_SK)

	return &document.Document{
		Key:     key,
		Content: itemMap,
	}, nil
}

func (s *DynamoDocService) Set(key *document.Key, value map[string]interface{}) error {
	err := document.ValidateKey(key)
	if err != nil {
		return err
	}

	if value == nil {
		return fmt.Errorf("provide non-nil value")
	}

	// Construct DynamoDB attribute value object
	itemMap := createItemMap(value, key)
	itemAttributeMap, err := dynamodbattribute.MarshalMap(itemMap)
	if err != nil {
		return fmt.Errorf("failed to marshal value")
	}

	input := &dynamodb.PutItemInput{
		Item:      itemAttributeMap,
		TableName: getTableName(*key.Collection),
	}

	_, err = s.client.PutItem(input)
	if err != nil {
		return err
	}

	return nil
}

func (s *DynamoDocService) Delete(key *document.Key) error {
	err := document.ValidateKey(key)
	if err != nil {
		return err
	}

	keyMap := createKeyMap(key)
	attributeMap, err := dynamodbattribute.MarshalMap(keyMap)
	if err != nil {
		return fmt.Errorf("failed to marshal keys: %v", key)
	}

	input := &dynamodb.DeleteItemInput{
		Key:       attributeMap,
		TableName: getTableName(*key.Collection),
	}

	_, err = s.client.DeleteItem(input)
	if err != nil {
		return fmt.Errorf("error deleting %v item %v : %v", key.Collection, key.Id, err)
	}

	// TODO: delete sub collection records

	return nil
}

func (s *DynamoDocService) Query(collection *document.Collection, expressions []document.QueryExpression, limit int, pagingToken map[string]string) (*document.QueryResult, error) {
	err := document.ValidateQueryCollection(collection)
	if err != nil {
		return nil, err
	}

	err = document.ValidateExpressions(expressions)
	if err != nil {
		return nil, err
	}

	queryResult := &document.QueryResult{
		Documents: make([]document.Document, 0),
	}

	// If partion key defined then perform a query
	if collection.Parent != nil && collection.Parent.Id != "" {
		err := s.performQuery(collection, expressions, limit, pagingToken, queryResult)
		if err != nil {
			return nil, err
		}

		remainingLimit := limit - len(queryResult.Documents)

		// If more results available, perform additional queries
		for remainingLimit > 0 &&
			(queryResult.PagingToken != nil && len(queryResult.PagingToken) > 0) {

			err := s.performQuery(collection, expressions, remainingLimit, queryResult.PagingToken, queryResult)
			if err != nil {
				return nil, err
			}

			remainingLimit = limit - len(queryResult.Documents)
		}

	} else {
		err := s.performScan(collection, expressions, limit, pagingToken, queryResult)
		if err != nil {
			return nil, err
		}

		remainingLimit := limit - len(queryResult.Documents)

		// If more results available, perform additional scans
		for remainingLimit > 0 &&
			(queryResult.PagingToken != nil && len(queryResult.PagingToken) > 0) {

			err := s.performScan(collection, expressions, remainingLimit, queryResult.PagingToken, queryResult)
			if err != nil {
				return nil, err
			}

			remainingLimit = limit - len(queryResult.Documents)
		}
	}

	return queryResult, nil
}

// New - Create a new DynamoDB key value plugin implementation
func New() (document.DocumentService, error) {
	awsRegion := utils.GetEnv("AWS_REGION", "us-east-1")

	// Create a new AWS session
	sess, sessionError := session.NewSession(&aws.Config{
		// FIXME: Use env config
		Region: aws.String(awsRegion),
	})

	if sessionError != nil {
		return nil, fmt.Errorf("error creating new AWS session %v", sessionError)
	}

	dynamoClient := dynamodb.New(sess)

	return &DynamoDocService{
		client: dynamoClient,
	}, nil
}

// NewWithClient - Mainly used for testing
func NewWithClient(client *dynamodb.DynamoDB) (document.DocumentService, error) {
	return &DynamoDocService{
		client: client,
	}, nil
}

// Private Functions ----------------------------------------------------------

func createKeyMap(key *document.Key) map[string]string {
	keyMap := make(map[string]string)

	parentKey := key.Collection.Parent

	if parentKey == nil {
		keyMap[ATTRIB_PK] = key.Id
		keyMap[ATTRIB_SK] = key.Collection.Name + "#"

	} else {
		keyMap[ATTRIB_PK] = parentKey.Id
		keyMap[ATTRIB_SK] = key.Collection.Name + "#" + key.Id
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
	newMap[ATTRIB_PK] = keyMap[ATTRIB_PK]
	newMap[ATTRIB_SK] = keyMap[ATTRIB_SK]

	return newMap
}

func (s *DynamoDocService) performQuery(
	collection *document.Collection,
	expressions []document.QueryExpression,
	limit int,
	pagingToken map[string]string,
	queryResult *document.QueryResult) error {

	if collection.Parent == nil {
		// Should never occur
		return fmt.Errorf("cannot perform query without partion key defined")
	}

	// Sort expressions to help map where "A >= %1 AND A <= %2" to DynamoDB expression "A BETWEEN %1 AND %2"
	sort.Sort(document.ExpsSort(expressions))

	input := &dynamodb.QueryInput{
		TableName: getTableName(*collection),
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
	input.ExpressionAttributeNames = make(map[string]*string)
	input.ExpressionAttributeNames["#pk"] = aws.String("_pk")
	input.ExpressionAttributeNames["#sk"] = aws.String("_sk")
	for _, exp := range expressions {
		input.ExpressionAttributeNames["#"+exp.Operand] = aws.String(exp.Operand)
	}

	// Configure ExpressionAttributeValues
	input.ExpressionAttributeValues = make(map[string]*dynamodb.AttributeValue)
	input.ExpressionAttributeValues[":pk"] = &dynamodb.AttributeValue{
		S: aws.String(collection.Parent.Id),
	}
	input.ExpressionAttributeValues[":sk"] = &dynamodb.AttributeValue{
		S: aws.String(collection.Name + "#"),
	}
	for i, exp := range expressions {
		expKey := fmt.Sprintf(":%v%v", exp.Operand, i)
		valAttrib, err := dynamodbattribute.Marshal(exp.Value)
		if err != nil {
			return fmt.Errorf("error marshalling %v: %v", exp.Operand, exp.Value)
		}
		input.ExpressionAttributeValues[expKey] = valAttrib
	}

	// Configure fetch Limit
	if limit > 0 {
		limit64 := int64(limit)
		input.Limit = &(limit64)

		if len(pagingToken) > 0 {
			startKey, err := dynamodbattribute.MarshalMap(pagingToken)
			if err != nil {
				return fmt.Errorf("error performing query %v: %v", input, err)
			}
			input.SetExclusiveStartKey(startKey)
		}
	}

	// Perform query
	resp, err := s.client.Query(input)

	if err != nil {
		return fmt.Errorf("error performing query %v: %v", input, err)
	}

	return marshalQueryResult(collection, resp.Items, resp.LastEvaluatedKey, limit, queryResult)
}

func (s *DynamoDocService) performScan(
	collection *document.Collection,
	expressions []document.QueryExpression,
	limit int,
	pagingToken map[string]string,
	queryResult *document.QueryResult) error {

	// Sort expressions to help map where "A >= %1 AND A <= %2" to DynamoDB expression "A BETWEEN %1 AND %2"
	sort.Sort(document.ExpsSort(expressions))

	input := &dynamodb.ScanInput{
		TableName: getTableName(*collection),
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
	input.ExpressionAttributeNames = make(map[string]*string)
	input.ExpressionAttributeNames["#sk"] = aws.String("_sk")

	for _, exp := range expressions {
		input.ExpressionAttributeNames["#"+exp.Operand] = aws.String(exp.Operand)
	}

	// Configure ExpressionAttributeValues
	input.ExpressionAttributeValues = make(map[string]*dynamodb.AttributeValue)
	keyAttrib := &dynamodb.AttributeValue{S: aws.String(collection.Name + "#")}

	input.ExpressionAttributeValues[":sk"] = keyAttrib
	for i, exp := range expressions {
		expKey := fmt.Sprintf(":%v%v", exp.Operand, i)
		valAttrib, err := dynamodbattribute.Marshal(exp.Value)
		if err != nil {
			return fmt.Errorf("error marshalling %v: %v", exp.Operand, exp.Value)
		}
		input.ExpressionAttributeValues[expKey] = valAttrib
	}

	// Configure fetch Limit
	if limit > 0 {
		// Account for parent record in fetch limit
		limit64 := int64(limit)
		input.Limit = &(limit64)

		if len(pagingToken) > 0 {
			startKey, err := dynamodbattribute.MarshalMap(pagingToken)
			if err != nil {
				return fmt.Errorf("error performing scan %v: %v", input, err)
			}
			input.SetExclusiveStartKey(startKey)
		}
	}

	resp, err := s.client.Scan(input)

	if err != nil {
		return fmt.Errorf("error performing scan %v: %v", input, err)
	}

	return marshalQueryResult(collection, resp.Items, resp.LastEvaluatedKey, limit, queryResult)
}

func marshalQueryResult(
	collection *sdk.Collection,
	items []map[string]*dynamodb.AttributeValue,
	lastEvaluatedKey map[string]*dynamodb.AttributeValue,
	limit int,
	queryResult *document.QueryResult) error {

	// Unmarshal Dynamo response items
	var valueMaps []map[string]interface{}
	err := dynamodbattribute.UnmarshalListOfMaps(items, &valueMaps)
	if err != nil {
		return fmt.Errorf("error unmarshalling query response: %v", err)
	}

	// Strip keys & append results
	for _, m := range valueMaps {
		// Retrieve the original ID on the result
		var id string
		var c *sdk.Collection
		if collection.Parent == nil {
			// We know this is a root document so its key will be located in PK
			pk, _ := m[ATTRIB_PK].(string)
			id = pk
			c = collection
		} else {
			// We know this is a child document so its key will be located in the SK
			pk, _ := m[ATTRIB_PK].(string)
			sk, _ := m[ATTRIB_SK].(string)
			idStr := strings.Split(sk, "#")
			id = idStr[len(idStr)-1]
			c = &sdk.Collection{
				Name: collection.Name,
				Parent: &sdk.Key{
					Collection: &sdk.Collection{
						Name: collection.Parent.Collection.Name,
					},
					Id: pk,
				},
			}
		}

		// Split out sort key value
		delete(m, ATTRIB_PK)
		delete(m, ATTRIB_SK)

		sdkDoc := document.Document{
			Key: &document.Key{
				Collection: c,
				Id:         id,
			},
			Content: m,
		}
		queryResult.Documents = append(queryResult.Documents, sdkDoc)
	}

	// Unmarshal lastEvalutedKey
	var resultPagingToken map[string]string
	if len(lastEvaluatedKey) > 0 {
		err = dynamodbattribute.UnmarshalMap(lastEvaluatedKey, &resultPagingToken)
		if err != nil {
			return fmt.Errorf("error unmarshalling query lastEvaluatedKey: %v", err)
		}
		queryResult.PagingToken = resultPagingToken
	}

	return nil
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

func getTableName(collection document.Collection) *string {
	coll := collection
	for coll.Parent != nil {
		coll = *coll.Parent.Collection
	}

	return aws.String(coll.Name)
}
