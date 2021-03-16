package dynamodb_service

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/dynamodb/expression"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/nitric-dev/membrane/plugins/sdk"
	"github.com/nitric-dev/membrane/utils"
)

type NitricDocument struct {
	Key   string
	Value map[string]interface{}
}

// AWS DynamoDB AWS nitric plugin
type DynamoDbDocumentService struct {
	sdk.UnimplementedDocumentsPlugin
	client dynamodbiface.DynamoDBAPI
}

func (s *DynamoDbDocumentService) createStandardDocumentTable(name string) error {
	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("Key"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("Key"),
				KeyType:       aws.String("HASH"),
			},
		},
		// TODO: This value is dependent on BillingMode, determine how to handle this. See: https://docs.aws.amazon.com/amazondynamodb/latest/APIReference/API_CreateTable.html#DDB-CreateTable-request-ProvisionedThroughput
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
		TableName: aws.String(name),
	}

	_, err := s.client.CreateTable(input)
	if err != nil {
		return fmt.Errorf("failed to create new dynamodb document table, with name %v. details: %v", name, err)
	}
	return nil
}

func (s *DynamoDbDocumentService) Create(collection string, key string, document map[string]interface{}) error {
	if key == "" {
		return fmt.Errorf("key auto-generation unimplemented, provide non-blank key")
	}

	// Construct DynamoDB attribute value object
	av, err := dynamodbattribute.MarshalMap(NitricDocument{
		Key:   key,
		Value: document,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal document")
	}

	condition := expression.AttributeNotExists(expression.Name("Key"))
	expr, err := expression.NewBuilder().WithCondition(condition).Build()

	if err != nil {
		return fmt.Errorf("failed to generate put request: %v", err)
	}

	// Store the NitricDocument attribute value to the specified table (collection)
	input := &dynamodb.PutItemInput{
		Item:                     av,
		TableName:                aws.String(collection),
		ExpressionAttributeNames: expr.Names(),
		ConditionExpression:      expr.Condition(),
	}

	var _, putError = s.client.PutItem(input)
	if putError != nil {
		if awsErr, ok := putError.(awserr.Error); ok {
			// Table not found,  try to create and put again
			if awsErr.Code() == dynamodb.ErrCodeResourceNotFoundException {
				createError := s.createStandardDocumentTable(collection)
				if createError != nil {
					return fmt.Errorf("table not found and failed to create: %v", createError)
				}
				_, putError = s.client.PutItem(input)
			}
			// Condition Expression failed
			if awsErr.Code() == dynamodb.ErrCodeConditionalCheckFailedException {
				return fmt.Errorf("document exists with key: %v", key)
			}
		}
	}

	if putError != nil {
		return fmt.Errorf("error creating new document: %v", putError)
	}

	return nil
}

func (s *DynamoDbDocumentService) Get(collection string, key string) (map[string]interface{}, error) {
	tableName := collection

	input := &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"Key": {
				S: aws.String(key),
			},
		},
	}

	result, getError := s.client.GetItem(input)
	if getError != nil {
		return nil, fmt.Errorf("error getting document: %v", getError)
	}

	if result.Item == nil {
		return nil, fmt.Errorf("Document not found")
	}

	document := NitricDocument{}
	unmarshalError := dynamodbattribute.UnmarshalMap(result.Item, &document)
	if unmarshalError != nil {
		return nil, fmt.Errorf("Failed to unmarshal document %v", unmarshalError)
	}

	return document.Value, nil
}

func (s *DynamoDbDocumentService) Update(collection string, key string, document map[string]interface{}) error {
	if key == "" {
		return fmt.Errorf("key auto-generation unimplemented, provide non-blank key")
	}

	// Construct DynamoDB attribute value object
	av, err := dynamodbattribute.MarshalMap(NitricDocument{
		Key:   key,
		Value: document,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal document")
	}

	condition := expression.AttributeExists(expression.Name("Key"))
	expr, err := expression.NewBuilder().WithCondition(condition).Build()

	if err != nil {
		return fmt.Errorf("failed to generate put request: %v", err)
	}

	// Store the NitricDocument attribute value to the specified table (collection)
	input := &dynamodb.PutItemInput{
		Item:                     av,
		TableName:                aws.String(collection),
		ExpressionAttributeNames: expr.Names(),
		ConditionExpression:      expr.Condition(),
	}

	var _, putError = s.client.PutItem(input)
	if putError != nil {
		if awsErr, ok := putError.(awserr.Error); ok {
			// Table not found,  try to create and put again
			if awsErr.Code() == dynamodb.ErrCodeResourceNotFoundException {
				createError := s.createStandardDocumentTable(collection)
				if createError != nil {
					return fmt.Errorf("table not found and failed to create: %v", createError)
				}
				_, putError = s.client.PutItem(input)
			}
			// Condition Expression failed
			if awsErr.Code() == dynamodb.ErrCodeConditionalCheckFailedException {
				return fmt.Errorf("no document found with key: %v", key)
			}
		}
	}

	if putError != nil {
		return fmt.Errorf("error creating new document: %v", putError)
	}

	return nil
}

func (s *DynamoDbDocumentService) Delete(collection string, key string) error {
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(collection),
		Key: map[string]*dynamodb.AttributeValue{
			"Key": {
				S: aws.String(key),
			},
		},
	}

	_, err := s.client.DeleteItem(input)
	if err != nil {
		return fmt.Errorf("error deleting document: %v", err)
	}

	return nil
}

// Create a New DynamoDB document plugin implementation
func New() (sdk.DocumentService, error) {
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

	return &DynamoDbDocumentService{
		client: dynamoClient,
	}, nil
}

// Mainly used for mock testing to inject a mock client into this plugin
func NewWithClient(client dynamodbiface.DynamoDBAPI) (sdk.DocumentService, error) {
	return &DynamoDbDocumentService{
		client: client,
	}, nil
}
