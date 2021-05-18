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

package dynamodb_service_test

import (
	"fmt"
	"os"
	"sort"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	plugin "github.com/nitric-dev/membrane/plugins/kv/dynamodb"
	"github.com/nitric-dev/membrane/sdk"
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
	}, nil
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
	key := input.Key[plugin.KEY].S
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
	key := input.Key[plugin.KEY].S

	if _, ok := m.store[*key]; ok {
		m.store[*key] = nil
		return &dynamodb.DeleteItemOutput{}, nil
	}

	// TODO: match real error codes.
	return nil, awserr.New("TESTERR", "Document does not exist!", fmt.Errorf("No document found"))
}

func (m *mockDynamoDBClient) Query(input *dynamodb.QueryInput) (*dynamodb.QueryOutput, error) {
	queryOutput := &dynamodb.QueryOutput{}
	queryOutput.Items = make([]map[string]*dynamodb.AttributeValue, 0)

	i := int64(0)

	for _, k := range sortKeys(m.store[*input.TableName]) {
		row := m.store[*input.TableName][k]

		if input.Limit != nil && *input.Limit >= 0 && i == *input.Limit {
			break
		}

		key := *input.ExpressionAttributeValues[":key"].S
		if key == k {
			attValue, _ := dynamodbattribute.MarshalMap(plugin.NitricKVDocument{
				Value: row,
			})
			queryOutput.Items = append(queryOutput.Items, attValue)
		}

		i += 1
	}

	return queryOutput, nil
}

func (m *mockDynamoDBClient) Scan(input *dynamodb.ScanInput) (*dynamodb.ScanOutput, error) {
	scanOutput := &dynamodb.ScanOutput{}
	scanOutput.Items = []map[string]*dynamodb.AttributeValue{}

	i := int64(0)

	for _, k := range sortKeys(m.store[*input.TableName]) {
		row := m.store[*input.TableName][k]

		if input.Limit != nil && *input.Limit >= 0 && i == *input.Limit {
			break
		}

		attValue, _ := dynamodbattribute.MarshalMap(plugin.NitricKVDocument{
			Value: row,
		})
		scanOutput.Items = append(scanOutput.Items, attValue)

		i += 1
	}

	return scanOutput, nil
}

var _ = Describe("DynamoDb", func() {
	defer GinkgoRecover()

	Context("Document Creation", func() {
		When("The dynamo-db client is operational", func() {
			// Setup Test
			key := map[string]interface{}{
				plugin.KEY: "Test",
			}
			item := map[string]interface{}{
				"Test": "Test",
			}
			mockSvc := &mockDynamoDBClient{
				store: make(map[string]map[string]map[string]interface{}),
			}

			// Inject the mock
			myMockClient, _ := plugin.NewWithClient(mockSvc)

			It("documentsClient.Create should store the document without error", func() {
				err := myMockClient.Put("Test", key, item)
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
				key := map[string]interface{}{
					plugin.KEY: "",
				}
				err := myMockClient.Put("Test", key, item)
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(ContainSubstring("provide non-empty key"))
			})
		})
	})

	Context("Document Retrieval", func() {
		When("The document exists", func() {
			key := map[string]interface{}{
				plugin.KEY: "Test",
			}
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
				storedItem, _ := myMockClient.Get("Test", key)

				Expect(storedItem).To(BeEquivalentTo(item))
			})
		})

		When("The document does not exist", func() {
			key := map[string]interface{}{
				plugin.KEY: "Test",
			}
			// Setup Test
			mockSvc := &mockDynamoDBClient{
				store: map[string]map[string]map[string]interface{}{},
			}
			// Inject the mock
			myMockClient, _ := plugin.NewWithClient(mockSvc)

			It("Should fail when attempting to retrieve the document", func() {
				_, err := myMockClient.Get("Test", key)

				Expect(err.Error()).To(ContainSubstring("error getting value for key"))
			})
		})
	})

	Context("Document Deletion", func() {
		When("The Document Exists", func() {
			key := map[string]interface{}{
				plugin.KEY: "Test",
			}
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
				_ = myMockClient.Delete("Test", key)
				_, ok := mockSvc.store["Test"]["Test"]

				Expect(ok).To(BeFalse())
			})
		})

		When("The document does not exist", func() {
			key := map[string]interface{}{
				plugin.KEY: "Test",
			}
			// Setup Test
			mockSvc := &mockDynamoDBClient{
				store: map[string]map[string]map[string]interface{}{},
			}
			// Inject the mock
			myMockClient, _ := plugin.NewWithClient(mockSvc)

			It("Should delete the stored document", func() {
				err := myMockClient.Delete("Test", key)

				Expect(err.Error()).To(ContainSubstring("error deleting key"))
			})
		})
	})

	Context("Updating Documents", func() {
		When("The document does not exist", func() {
			key := map[string]interface{}{
				plugin.KEY: "Test",
			}
			item := map[string]interface{}{
				"Test": "Test",
			}
			mockSvc := &mockDynamoDBClient{
				store: make(map[string]map[string]map[string]interface{}),
			}

			// Inject the mock
			myMockClient, _ := plugin.NewWithClient(mockSvc)

			It("Should behave as Create", func() {
				err := myMockClient.Put("Test", key, item)

				Expect(err).To(BeNil())
				storedItem, ok := mockSvc.store["Test"]["Test"]
				Expect(ok).To(BeTrue())
				Expect(storedItem).To(BeEquivalentTo(item))
			})
		})
	})

	Context("KV Query", func() {
		// Setup Test & Inject Mock into client
		mockSvc := &mockDynamoDBClient{
			store: map[string]map[string]map[string]interface{}{},
		}
		myMockClient, _ := plugin.NewWithClient(mockSvc)

		When("collection is blank", func() {
			It("Should return an error", func() {
				results, err := myMockClient.Query("", nil, 0)
				Expect(results).To(BeNil())
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(ContainSubstring("provide non-blank collection"))
			})
		})
		When("expressions is nill", func() {
			It("Should return an error", func() {
				results, err := myMockClient.Query("collection", nil, 0)
				Expect(results).To(BeNil())
				Expect(err.Error()).To(ContainSubstring("provide non-nil expressions"))
			})
		})

		// Setup Test & Inject Mock into client
		mockSvc = &mockDynamoDBClient{
			store: map[string]map[string]map[string]interface{}{
				"table": {
					"key1": {
						"email": "user1@server.com",
					},
					"key2": {
						"email": "user2@server.com",
					},
					"key3": {
						"email": "user3@server.com",
					},
				},
			},
		}
		myMockClient, _ = plugin.NewWithClient(mockSvc)

		When("scan with empty expression and no limit", func() {
			It("return all items", func() {
				exps := []sdk.QueryExpression{}
				results, err := myMockClient.Query("table", exps, 0)

				Expect(err).To(BeNil())
				Expect(results).NotTo(BeNil())
				Expect(len(results)).To(BeIdenticalTo(3))
				Expect(results[0]).To(BeEquivalentTo(map[string]interface{}{"email": "user1@server.com"}))
				Expect(results[1]).To(BeEquivalentTo(map[string]interface{}{"email": "user2@server.com"}))
				Expect(results[2]).To(BeEquivalentTo(map[string]interface{}{"email": "user3@server.com"}))
			})
		})

		When("scan with empty expression and 1 limit", func() {
			It("return first item", func() {
				exps := []sdk.QueryExpression{}
				results, err := myMockClient.Query("table", exps, 1)

				Expect(err).To(BeNil())
				Expect(results).NotTo(BeNil())
				Expect(len(results)).To(BeIdenticalTo(1))
				Expect(results[0]).To(BeEquivalentTo(map[string]interface{}{"email": "user1@server.com"}))
			})
		})

		When("query with one expression", func() {
			It("return second item", func() {
				exps := []sdk.QueryExpression{
					{Operand: "key", Operator: "==", Value: "key2"},
				}
				results, err := myMockClient.Query("table", exps, 0)

				Expect(err).To(BeNil())
				Expect(results).NotTo(BeNil())
				Expect(results[0]).To(BeEquivalentTo(map[string]interface{}{"email": "user2@server.com"}))
			})
		})
	})

	// TODO: develop integration testing mode: export INTEGRATION_TESTING=true
	integrationTesting := os.Getenv("INTEGRATION_TESTING")
	if integrationTesting == "" {
		return
	}

	// Perform Local DynamoDB integration testing with pre-loaded dataset
	Context("Local DynamoDB", func() {
		// create an aws session
		sess := session.Must(session.NewSession(&aws.Config{
			Region:   aws.String("x"),
			Endpoint: aws.String("http://127.0.0.1:8000"),
		}))

		localClient, _ := plugin.NewWithClient(dynamodb.New(sess))

		When("customer table query with no expressions", func() {
			It("return two items", func() {
				exps := []sdk.QueryExpression{}
				results, err := localClient.Query("customer", exps, 0)

				Expect(err).To(BeNil())
				Expect(results).NotTo(BeNil())
				Expect(len(results)).To(BeIdenticalTo(2))
				Expect(results[0]["email"]).To(BeEquivalentTo("jane.smith@server.com"))
				Expect(results[1]["email"]).To(BeEquivalentTo("paul.davis@server.com"))
			})
		})

		When("customer table query with no expressions, 1 limit", func() {
			It("return first item", func() {
				exps := []sdk.QueryExpression{}
				results, err := localClient.Query("customer", exps, 1)

				Expect(err).To(BeNil())
				Expect(results).NotTo(BeNil())
				Expect(len(results)).To(BeIdenticalTo(1))
				Expect(results[0]["email"]).To(BeEquivalentTo("jane.smith@server.com"))
			})
		})

		When("customer table query with one single key expression", func() {
			It("return second item", func() {
				exps := []sdk.QueryExpression{
					{Operand: "key", Operator: "==", Value: "paul.davis@server.com"},
				}
				results, err := localClient.Query("customer", exps, 0)

				Expect(err).To(BeNil())
				Expect(results).NotTo(BeNil())
				Expect(len(results)).To(BeIdenticalTo(1))
				Expect(results[0]["email"]).To(BeEquivalentTo("paul.davis@server.com"))
			})
		})

		When("application table query with customer key expression", func() {
			It("return customer items", func() {
				exps := []sdk.QueryExpression{
					{Operand: "pk", Operator: "==", Value: "Customer#1000"},
				}
				results, err := localClient.Query("application", exps, 0)

				Expect(err).To(BeNil())
				Expect(results).NotTo(BeNil())
				Expect(len(results)).To(BeIdenticalTo(4))
				Expect(results[0]["firstName"]).To(BeEquivalentTo("Jane"))
				Expect(results[1]["order"]).To(BeEquivalentTo("501"))
				Expect(results[2]["order"]).To(BeEquivalentTo("502"))
				Expect(results[3]["order"]).To(BeEquivalentTo("503"))
			})
		})

		When("application table query with pk customer and sk order expressions", func() {
			It("return customer orders items", func() {
				exps := []sdk.QueryExpression{
					{Operand: "pk", Operator: "==", Value: "Customer#1000"},
					{Operand: "sk", Operator: "startsWith", Value: "Order#"},
				}
				results, err := localClient.Query("application", exps, 0)

				Expect(err).To(BeNil())
				Expect(results).NotTo(BeNil())
				Expect(len(results)).To(BeIdenticalTo(3))
				Expect(results[0]["order"]).To(BeEquivalentTo("501"))
				Expect(results[1]["order"]).To(BeEquivalentTo("502"))
				Expect(results[2]["order"]).To(BeEquivalentTo("503"))
			})
		})

		When("application table query with pk customer and sk order expressions", func() {
			It("return customer order item", func() {
				exps := []sdk.QueryExpression{
					{Operand: "pk", Operator: "==", Value: "Customer#1000"},
					{Operand: "sk", Operator: ">", Value: "Order#502"},
				}
				results, err := localClient.Query("application", exps, 0)

				Expect(err).To(BeNil())
				Expect(results).NotTo(BeNil())
				Expect(len(results)).To(BeIdenticalTo(1))
				Expect(results[0]["order"]).To(BeEquivalentTo("503"))
			})
		})

		When("application table scan with no expressions", func() {
			It("return all items", func() {
				results, err := localClient.Query("application", []sdk.QueryExpression{}, 0)

				Expect(err).To(BeNil())
				Expect(results).NotTo(BeNil())
				Expect(len(results)).To(BeIdenticalTo(5))
				Expect(results[0]["firstName"]).To(BeEquivalentTo("Jane"))
				Expect(results[4]["sku"]).To(BeEquivalentTo("171-823-623"))
			})
		})
	})
})

func sortKeys(m map[string]map[string]interface{}) (keys []string) {

	keys = make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	return keys
}
