package dynamodb_service_test

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	plugin "github.com/nitric-dev/membrane/plugins/aws/kv/dynamodb"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// Define a mock struct to be used in your unit tests of myFunc.
type mockDynamoDBClient struct {
	dynamodbiface.DynamoDBAPI
	store map[string]map[string]map[string]interface{}
}

func (m *mockDynamoDBClient) CreateTable(input *dynamodb.CreateTableInput) (*dynamodb.CreateTableOutput, error) {
	// mock response/functionality
	tableName := input.TableName
	if _, ok := m.store[*tableName]; ok {
		// Already exists, throw an error
		// TODO: Throw predefined AWS error...
		return nil, awserr.New("Some code", "Table already exists", fmt.Errorf("Could not create table"))
	} else {
		m.store[*tableName] = make(map[string]map[string]interface{})
	}

	// Currently the output is ignored in our usecase so leave this empty for now...
	return &dynamodb.CreateTableOutput{
		TableDescription: &dynamodb.TableDescription{
			TableStatus: aws.String("CREATING"),
		},
	} , nil
}

func (m *mockDynamoDBClient) DescribeTable(input *dynamodb.DescribeTableInput) (*dynamodb.DescribeTableOutput, error) {
	return &dynamodb.DescribeTableOutput{Table: &dynamodb.TableDescription{
		TableStatus: aws.String("ACTIVE"),
	}}, nil
}

func (m *mockDynamoDBClient) PutItem(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	tableName := input.TableName
	item := input.Item
	mapValue := plugin.NitricKVDocument{}

	dynamodbattribute.UnmarshalMap(item, &mapValue)

	// mock response/functionality
	if _, ok := m.store[*tableName]; ok {
		m.store[*tableName][mapValue.Key] = mapValue.Value
		// Already exists, throw an error
		// TODO: Throw predefined AWS error...
		return &dynamodb.PutItemOutput{}, nil
	} else {
		return nil, awserr.New(dynamodb.ErrCodeResourceNotFoundException, "Table does not exist", fmt.Errorf("No table found"))
	}

}

func (m *mockDynamoDBClient) GetItem(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
	// mock response/functionality
	key := input.Key["Key"].S
	tableName := input.TableName

	if item, ok := m.store[*tableName][*key]; ok {
		attValue, _ := dynamodbattribute.MarshalMap(plugin.NitricKVDocument{
			Key:   *key,
			Value: item,
		})
		return &dynamodb.GetItemOutput{Item: attValue}, nil
	}

	// TODO: match real error codes.
	return nil, awserr.New("TESTERR", "Document does not exist!", fmt.Errorf("No document found"))

}

func (m *mockDynamoDBClient) DeleteItem(input *dynamodb.DeleteItemInput) (*dynamodb.DeleteItemOutput, error) {
	// mock response/functionality
	key := input.Key["Key"].S

	if _, ok := m.store[*key]; ok {
		m.store[*key] = nil
		return &dynamodb.DeleteItemOutput{}, nil
	}

	// TODO: match real error codes.
	return nil, awserr.New("TESTERR", "Document does not exist!", fmt.Errorf("No document found"))
}

var _ = Describe("DynamoDb", func() {
	Context("Document Creation", func() {
		When("The dynamo-db client is operational", func() {
			// Setup Test
			item := map[string]interface{}{
				"Test": "Test",
			}
			mockSvc := &mockDynamoDBClient{
				store: make(map[string]map[string]map[string]interface{}),
			}

			// Inject the mock
			myMockClient, _ := plugin.NewWithClient(mockSvc)

			It("documentsClient.Create should store the document without error", func() {
				err := myMockClient.Put("Test", "Test", item)
				Expect(err).To(BeNil())

				storedItem, ok := mockSvc.store["Test"]["Test"]

				Expect(ok).To(BeTrue())
				Expect(storedItem).To(BeEquivalentTo(item))
			})
		})

		When("Creating a new document without a key", func() {
			// Setup Test
			item := map[string]interface{}{
				"Test": "Test",
			}
			mockSvc := &mockDynamoDBClient{
				store: make(map[string]map[string]map[string]interface{}),
			}

			// Inject the mock
			myMockClient, _ := plugin.NewWithClient(mockSvc)

			It("should return an error", func() {
				err := myMockClient.Put("Test", "", item)

				Expect(err.Error()).To(ContainSubstring("key auto-generation unimplemented, provide non-blank key"))
			})
		})
	})

	Context("Document Retrieval", func() {
		When("The document exists", func() {
			item := map[string]interface{}{
				"Test": "Test",
			}
			// Setup Test
			mockSvc := &mockDynamoDBClient{
				store: map[string]map[string]map[string]interface{}{
					"Test": {
						"Test": item,
					},
				},
			}
			// Inject the mock
			myMockClient, _ := plugin.NewWithClient(mockSvc)

			It("Should retrieve the stored document", func() {
				storedItem, _ := myMockClient.Get("Test", "Test")

				Expect(storedItem).To(BeEquivalentTo(item))
			})
		})

		When("The document does not exist", func() {
			// Setup Test
			mockSvc := &mockDynamoDBClient{
				store: map[string]map[string]map[string]interface{}{},
			}
			// Inject the mock
			myMockClient, _ := plugin.NewWithClient(mockSvc)

			It("Should fail when attempting to retrieve the document", func() {
				_, err := myMockClient.Get("Test", "Test")

				Expect(err.Error()).To(ContainSubstring("error getting value for key"))
			})
		})
	})

	Context("Document Deletion", func() {
		When("The Document Exists", func() {
			item := map[string]interface{}{
				"Test": "Test",
			}
			// Setup Test
			mockSvc := &mockDynamoDBClient{
				store: map[string]map[string]map[string]interface{}{
					"Test": {
						"Test": item,
					},
				},
			}
			// Inject the mock
			myMockClient, _ := plugin.NewWithClient(mockSvc)

			It("Should delete the stored document", func() {
				_ = myMockClient.Delete("Test", "Test")
				_, ok := mockSvc.store["Test"]["Test"]

				Expect(ok).To(BeFalse())
			})
		})

		When("The document does not exist", func() {
			// Setup Test
			mockSvc := &mockDynamoDBClient{
				store: map[string]map[string]map[string]interface{}{},
			}
			// Inject the mock
			myMockClient, _ := plugin.NewWithClient(mockSvc)

			It("Should delete the stored document", func() {
				err := myMockClient.Delete("Test", "Test")

				Expect(err.Error()).To(ContainSubstring("error deleting key"))
			})
		})
	})

	Context("Updating Documents", func() {
		When("The document does not exist", func() {
			item := map[string]interface{}{
				"Test": "Test",
			}
			mockSvc := &mockDynamoDBClient{
				store: make(map[string]map[string]map[string]interface{}),
			}

			// Inject the mock
			myMockClient, _ := plugin.NewWithClient(mockSvc)

			It("Should behave as Create", func() {
				err := myMockClient.Put("Test", "Test", item)

				Expect(err).To(BeNil())
				storedItem, ok := mockSvc.store["Test"]["Test"]
				Expect(ok).To(BeTrue())
				Expect(storedItem).To(BeEquivalentTo(item))
			})
		})
	})
})
