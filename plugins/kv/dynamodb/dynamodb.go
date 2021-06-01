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

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/nitric-dev/membrane/plugins/kv"
	"github.com/nitric-dev/membrane/sdk"
	"github.com/nitric-dev/membrane/utils"
)

// NitricKVDocument - represents the structure of a Key Value record when stored in DynamoDB
type NitricKVDocument struct {
	Value map[string]interface{} `json:"value"`
}

// AWS DynamoDB AWS Nitric Key Value service
type DynamoDbKVService struct {
	sdk.UnimplementedKeyValuePlugin
	client dynamodbiface.DynamoDBAPI
}

func copy(source map[string]interface{}) map[string]interface{} {
	newMap := make(map[string]interface{})
	for key, value := range source {
		newMap[key] = value
	}
	return newMap
}

func marshalListOfMaps(items []map[string]*dynamodb.AttributeValue) ([]map[string]interface{}, error) {
	// Unmarshall Dynamo response items into Doc struct, the marshall into result map
	var valueDocs []NitricKVDocument
	err := dynamodbattribute.UnmarshalListOfMaps(items, &valueDocs)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling query response: %v", err)
	}

	results := []map[string]interface{}{}
	for _, m := range valueDocs {
		results = append(results, m.Value)
	}

	return results, nil
}

func (s *DynamoDbKVService) Put(collection string, key map[string]interface{}, value map[string]interface{}) error {
	err := kv.ValidateCollection(collection)
	if err != nil {
		return err
	}
	err = kv.ValidateKeyMap(collection, key)
	if err != nil {
		return err
	}

	if value == nil {
		return fmt.Errorf("provide non-nil value")
	}

	// Construct DynamoDB attribute value object
	itemMap := copy(key)

	// Project any collection filter attributes into Doc filter attributes
	filterAttributes, err := kv.Stack.CollectionFilterAttributes(collection)
	if filterAttributes != nil && err == nil {
		for _, name := range filterAttributes {
			if _, found := value[name]; found {
				itemMap[name] = fmt.Sprintf("%v", value[name])
			}
		}
	}

	// Add value map
	itemMap["value"] = value

	itemAttributeMap, err := dynamodbattribute.MarshalMap(itemMap)
	if err != nil {
		return fmt.Errorf("failed to marshal value")
	}

	input := &dynamodb.PutItemInput{
		Item:      itemAttributeMap,
		TableName: aws.String(collection),
	}

	_, err = s.client.PutItem(input)
	if err != nil {
		return err
	}

	return nil
}

func (s *DynamoDbKVService) Get(collection string, key map[string]interface{}) (map[string]interface{}, error) {
	err := kv.ValidateCollection(collection)
	if err != nil {
		return nil, err
	}
	err = kv.ValidateKeyMap(collection, key)
	if err != nil {
		return nil, err
	}

	attributeMap, err := dynamodbattribute.MarshalMap(key)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal key: %v", key)
	}

	input := &dynamodb.GetItemInput{
		Key:       attributeMap,
		TableName: aws.String(collection),
	}

	result, err := s.client.GetItem(input)
	if err != nil {
		return nil, fmt.Errorf("error getting %v value %s : %v", collection, key, err)
	}

	if result.Item == nil {
		return nil, fmt.Errorf("%v value %v not found", collection, key)
	}

	kvDocument := NitricKVDocument{}
	unmarshalError := dynamodbattribute.UnmarshalMap(result.Item, &kvDocument)
	if unmarshalError != nil {
		return nil, fmt.Errorf("failed to unmarshal key value document: %v", unmarshalError)
	}

	return kvDocument.Value, nil
}

func (s *DynamoDbKVService) Delete(collection string, key map[string]interface{}) error {
	err := kv.ValidateCollection(collection)
	if err != nil {
		return err
	}
	err = kv.ValidateKeyMap(collection, key)
	if err != nil {
		return err
	}

	attributeMap, err := dynamodbattribute.MarshalMap(key)
	if err != nil {
		return fmt.Errorf("failed to marshal key: %v", key)
	}

	input := &dynamodb.DeleteItemInput{
		Key:       attributeMap,
		TableName: aws.String(collection),
	}

	_, err = s.client.DeleteItem(input)
	if err != nil {
		return fmt.Errorf("error deleting %v item %v : %v", collection, key, err)
	}

	return nil
}

func (s *DynamoDbKVService) Query(collection string, expressions []sdk.QueryExpression, limit int) ([]map[string]interface{}, error) {
	err := kv.ValidateCollection(collection)
	if err != nil {
		return nil, err
	}
	err = kv.ValidateExpressions(collection, expressions)
	if err != nil {
		return nil, err
	}

	// Sort expressions to help map where "A >= %1 AND A <= %2" to DynamoDB expression "A BETWEEN %1 AND %2"
	sort.Sort(kv.ExpsSort(expressions))

	// If expressions perform a query
	if isQueryOperation(collection, expressions) {

		input := &dynamodb.QueryInput{
			TableName:            aws.String(collection),
			ProjectionExpression: aws.String("#value"),
		}

		// Configure KeyConditionExpression
		keyExp := createKeyExpression(collection, expressions)
		input.KeyConditionExpression = aws.String(keyExp)

		// Configure FilterExpression
		filterExp := createFilterExpression(collection, expressions)
		if filterExp != "" {
			input.FilterExpression = aws.String(filterExp)
		}

		// Configure ExpressionAttributeNames
		input.ExpressionAttributeNames = make(map[string]*string)
		for _, exp := range expressions {
			input.ExpressionAttributeNames["#"+exp.Operand] = aws.String(exp.Operand)
		}
		input.ExpressionAttributeNames["#value"] = aws.String("value")

		// Configure ExpressionAttributeValues
		input.ExpressionAttributeValues = make(map[string]*dynamodb.AttributeValue)
		for i, exp := range expressions {
			expKey := fmt.Sprintf(":%v%v", exp.Operand, i)
			input.ExpressionAttributeValues[expKey] = &dynamodb.AttributeValue{
				S: aws.String(exp.Value),
			}
		}

		// Configure fetch Limit
		if limit > 0 {
			limit64 := int64(limit)
			input.Limit = &(limit64)
		}

		// Perform query
		resp, err := s.client.Query(input)
		if err != nil {
			return nil, fmt.Errorf("error performing query %v: %v", input, err)
		}

		return marshalListOfMaps(resp.Items)

	} else {
		input := &dynamodb.ScanInput{
			TableName: aws.String(collection),
			ExpressionAttributeNames: map[string]*string{
				"#value": aws.String("value"),
			},
			ProjectionExpression: aws.String("#value"),
		}

		filterExp := createFilterExpression(collection, expressions)
		if filterExp != "" {
			// Configure FilterExpression
			input.FilterExpression = aws.String(filterExp)

			// Configure ExpressionAttributeNames
			input.ExpressionAttributeNames = make(map[string]*string)
			for _, exp := range expressions {
				input.ExpressionAttributeNames["#"+exp.Operand] = aws.String(exp.Operand)
			}
			input.ExpressionAttributeNames["#value"] = aws.String("value")

			// Configure ExpressionAttributeValues
			input.ExpressionAttributeValues = make(map[string]*dynamodb.AttributeValue)
			for i, exp := range expressions {
				expKey := fmt.Sprintf(":%v%v", exp.Operand, i)
				input.ExpressionAttributeValues[expKey] = &dynamodb.AttributeValue{
					S: aws.String(exp.Value),
				}
			}
		}

		// Configure fetch Limit
		if limit > 0 {
			limit64 := int64(limit)
			input.Limit = &(limit64)
		}

		resp, err := s.client.Scan(input)
		if err != nil {
			return nil, fmt.Errorf("error performing scan %v: %v", input, err)
		}

		return marshalListOfMaps(resp.Items)
	}
}

// Create a New DynamoDB key value plugin implementation
func New() (sdk.KeyValueService, error) {
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

	return &DynamoDbKVService{
		client: dynamoClient,
	}, nil
}

// Mainly used for testing
func NewWithClient(client *dynamodb.DynamoDB) (sdk.KeyValueService, error) {
	return &DynamoDbKVService{
		client: client,
	}, nil
}

// Return true if should perform a query operation against keys or false
// if should perform a Dynamodb scan operation
func isQueryOperation(collection string, exps []sdk.QueryExpression) bool {
	if len(exps) == 0 {
		return false
	}

	indexes, _ := kv.Stack.CollectionIndexes(collection)
	for _, exp := range exps {
		if utils.IndexOf(indexes, exp.Operand) != -1 {
			return true
		}
	}

	return false
}

func createKeyExpression(collection string, expressions []sdk.QueryExpression) string {
	indexAttributes, _ := kv.Stack.CollectionIndexes(collection)

	keyExp := ""
	for i, exp := range expressions {
		if utils.IndexOf(indexAttributes, exp.Operand) != -1 {
			if keyExp != "" {
				keyExp += " AND "
			}
			if exp.Operator == "startsWith" {
				keyExp += "begins_with(#" + exp.Operand + ", :" + fmt.Sprintf("%v%v", exp.Operand, i) + ")"

			} else if exp.Operator == "==" {
				keyExp += "#" + exp.Operand + " = :" + fmt.Sprintf("%v%v", exp.Operand, i)

			} else {
				keyExp += "#" + exp.Operand + " " + exp.Operator + " :" + fmt.Sprintf("%v%v", exp.Operand, i)
			}
		}
	}

	return keyExp
}

func createFilterExpression(collection string, expressions []sdk.QueryExpression) string {
	filterAttributes, _ := kv.Stack.CollectionFilterAttributes(collection)

	keyExp := ""
	for i, exp := range expressions {
		if utils.IndexOf(filterAttributes, exp.Operand) != -1 {
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
	}

	return keyExp
}

func isBetweenStart(index int, exps []sdk.QueryExpression) bool {
	if index < (len(exps) - 1) {
		if exps[index].Operand == exps[index+1].Operand &&
			exps[index].Operator == ">=" &&
			exps[index+1].Operator == "<=" {
			return true
		}
	}
	return false
}

func isBetweenEnd(index int, exps []sdk.QueryExpression) bool {
	if index > 0 && index < len(exps) {
		if exps[index-1].Operand == exps[index].Operand &&
			exps[index-1].Operator == ">=" &&
			exps[index].Operator == "<=" {
			return true
		}
	}
	return false
}
