package dynamodb_service

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/nitric-dev/membrane/plugins/sdk"
	"github.com/nitric-dev/membrane/utils"
)

// NitricKVDocument - represents the structure of a Key Value record when stored in DynamoDB
type NitricKVDocument struct {
	Key   string
	Value map[string]interface{}
}

// AWS DynamoDB AWS Nitric Key Value service
type DynamoDbKVService struct {
	sdk.UnimplementedKeyValuePlugin
	client dynamodbiface.DynamoDBAPI
}

func (s *DynamoDbKVService) createStandardKVTable(name string) error {
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
		return fmt.Errorf("failed to create new dynamodb key value table, with name %v. details: %v", name, err)
	}
	return nil
}

func (s *DynamoDbKVService) Put(collection string, key string, value map[string]interface{}) error {
	if key == "" {
		return fmt.Errorf("key auto-generation unimplemented, provide non-blank key")
	}

	// Construct DynamoDB attribute value object
	av, err := dynamodbattribute.MarshalMap(NitricKVDocument{
		Key:   key,
		Value: value,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal value")
	}

	if err != nil {
		return fmt.Errorf("failed to generate put request: %v", err)
	}

	// Store the NitricKVDocument attribute value to the specified table (collection)
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(collection),
	}

	var _, putError = s.client.PutItem(input)
	if putError != nil {
		if awsErr, ok := putError.(awserr.Error); ok {
			// Table not found,  try to create and put again
			if awsErr.Code() == dynamodb.ErrCodeResourceNotFoundException {
				createError := s.createStandardKVTable(collection)
				if createError != nil {
					return fmt.Errorf("table not found and failed to create: %v", createError)
				}
				_, putError = s.client.PutItem(input)
			}
		}
	}

	if putError != nil {
		return fmt.Errorf("error creating new value: %v", putError)
	}

	return nil
}

func (s *DynamoDbKVService) Get(collection string, key string) (map[string]interface{}, error) {
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
		return nil, fmt.Errorf("error getting value for key %s: %v", key, getError)
	}

	if result.Item == nil {
		return nil, fmt.Errorf("value not found")
	}

	kvDocument := NitricKVDocument{}
	unmarshalError := dynamodbattribute.UnmarshalMap(result.Item, &kvDocument)
	if unmarshalError != nil {
		return nil, fmt.Errorf("failed to unmarshal key value document: %v", unmarshalError)
	}

	return kvDocument.Value, nil
}

func (s *DynamoDbKVService) Delete(collection string, key string) error {
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
		return fmt.Errorf("error deleting key %s: %v", key, err)
	}

	return nil
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

// Mainly used for mock testing to inject a mock client into this plugin
func NewWithClient(client dynamodbiface.DynamoDBAPI) (sdk.KeyValueService, error) {
	return &DynamoDbKVService{
		client: client,
	}, nil
}
