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
	"os/exec"

	kv_plugin "github.com/nitric-dev/membrane/plugins/kv/dynamodb"
	data "github.com/nitric-dev/membrane/plugins/kv/test"
	"github.com/nitric-dev/membrane/sdk"
	"github.com/nitric-dev/membrane/utils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func startDynamoProcess() *exec.Cmd {
	// Start Local DynamoDB
	args := []string{
		"-Djava.library.path=/usr/local/dynamodb/DynamoDBLocal_lib",
		"-jar",
		"/usr/local/dynamodb/DynamoDBLocal.jar",
		"-inMemory",
	}
	cmd := exec.Command("/usr/bin/java", args[:]...)
	if err := cmd.Start(); err != nil {
		panic(fmt.Sprintf("Error starting Local DynamoDB %v : %v", cmd, err))
	}
	fmt.Printf("Started Local DynamoDB (PID %v) and loading data...\n", cmd.Process.Pid)

	return cmd
}

func stopDynamoProcess(cmd *exec.Cmd) {
	if err := cmd.Process.Kill(); err != nil {
		fmt.Printf("failed to kill DynamoDB %v : %v \n", cmd.Process.Pid, err)
	}
}

func createDynamoClient() *dynamodb.DynamoDB {
	sess := session.Must(session.NewSession(&aws.Config{
		Region:   aws.String("x"),
		Endpoint: aws.String("http://127.0.0.1:8000"),
	}))

	return dynamodb.New(sess)
}

func createApplicationTable(db *dynamodb.DynamoDB) {
	tableName := "application"

	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("pk"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("sk"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("pk"),
				KeyType:       aws.String("HASH"),
			},
			{
				AttributeName: aws.String("sk"),
				KeyType:       aws.String("RANGE"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
		TableName: aws.String(tableName),
	}
	_, err := db.CreateTable(input)
	if err != nil {
		panic(fmt.Sprintf("Error calling CreateTable: %s", err))
	}
}

func createUsersTable(db *dynamodb.DynamoDB) {
	tableName := "users"

	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("key"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("key"),
				KeyType:       aws.String("HASH"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
		TableName: aws.String(tableName),
	}
	_, err := db.CreateTable(input)
	if err != nil {
		panic(fmt.Sprintf("Error calling CreateTable: %s", err))
	}
}

func deleteTable(db *dynamodb.DynamoDB, tableName string) {
	deleteInput := &dynamodb.DeleteTableInput{
		TableName: aws.String(tableName),
	}

	_, err := db.DeleteTable(deleteInput)
	if err != nil {
		panic(fmt.Sprintf("Error calling DeleteTable: %s", err))
	}
}

var _ = Describe("DynamoDb", func() {
	defer GinkgoRecover()

	os.Setenv("AWS_ACCESS_KEY_ID", "fakeMyKeyId")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "fakeSecretAccessKey")
	os.Setenv("AWS_REGION", "X")

	// Start Local DynamoDB
	dynaCmd := startDynamoProcess()

	// Create DyanmoDB client
	db := createDynamoClient()

	BeforeEach(func() {
		createApplicationTable(db)
		createUsersTable(db)
	})

	AfterEach(func() {
		deleteTable(db, "application")
		deleteTable(db, "users")
	})

	AfterSuite(func() {
		stopDynamoProcess(dynaCmd)
	})

	os.Setenv(utils.NITRIC_HOME, "../test/")
	os.Setenv(utils.NITRIC_YAML, "nitric.yaml")

	kvPlugin, err := kv_plugin.NewWithClient(db)
	if err != nil {
		panic(err)
	}

	Context("Put", func() {
		When("Blank collection", func() {
			It("Should return error", func() {
				err := kvPlugin.Put("", data.UserKey1, data.UserItem1)
				Expect(err).Should(HaveOccurred())
			})
		})
		When("Nil key", func() {
			It("Should return error", func() {
				err := kvPlugin.Put("users", nil, data.UserItem1)
				Expect(err).Should(HaveOccurred())
			})
		})
		When("Nil item map", func() {
			It("Should return error", func() {
				err := kvPlugin.Put("users", data.UserKey1, nil)
				Expect(err).Should(HaveOccurred())
			})
		})
		When("Valid New Put", func() {
			It("Should store new item successfully", func() {
				err := kvPlugin.Put("users", data.UserKey1, data.UserItem1)
				Expect(err).ShouldNot(HaveOccurred())

				doc, err := kvPlugin.Get("users", data.UserKey1)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(doc).ToNot(BeNil())
				Expect(doc["email"]).To(BeEquivalentTo(data.UserItem1["email"]))
			})
		})
		When("Valid Update Put", func() {
			It("Should store new item successfully", func() {
				err := kvPlugin.Put("users", data.UserKey1, data.UserItem1)
				Expect(err).ShouldNot(HaveOccurred())

				doc, err := kvPlugin.Get("users", data.UserKey1)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(doc).ToNot(BeNil())
				Expect(doc["email"]).To(BeEquivalentTo(data.UserItem1["email"]))

				err = kvPlugin.Put("users", data.UserKey1, data.UserItem2)
				Expect(err).ShouldNot(HaveOccurred())

				doc, err = kvPlugin.Get("users", data.UserKey1)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(doc).ToNot(BeNil())
				Expect(doc["email"]).To(BeEquivalentTo(data.UserItem2["email"]))
			})
		})
		When("Valid Compound Key Put", func() {
			It("Should store item successfully", func() {
				err := kvPlugin.Put("application", data.OrderKey1, data.OrderItem1)
				Expect(err).ShouldNot(HaveOccurred())

				doc, err := kvPlugin.Get("application", data.OrderKey1)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(doc).ToNot(BeNil())
				Expect(doc).To(BeEquivalentTo(data.OrderItem1))
			})
		})
	})

	Context("Get", func() {
		When("Blank collection", func() {
			It("Should return error", func() {
				_, err := kvPlugin.Get("", data.UserKey1)
				Expect(err).Should(HaveOccurred())
			})
		})
		When("Nil key", func() {
			It("Should return error", func() {
				_, err := kvPlugin.Get("users", nil)
				Expect(err).Should(HaveOccurred())
			})
		})
		When("Valid Get", func() {
			It("Should get item successfully", func() {
				kvPlugin.Put("users", data.UserKey1, data.UserItem1)

				doc, err := kvPlugin.Get("users", data.UserKey1)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(doc).ToNot(BeNil())
				Expect(doc["email"]).To(BeEquivalentTo(data.UserItem1["email"]))
			})
		})
		When("Valid Compound Key Get", func() {
			It("Should store item successfully", func() {
				kvPlugin.Put("application", data.OrderKey1, data.OrderItem1)

				doc, err := kvPlugin.Get("application", data.OrderKey1)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(doc).ToNot(BeNil())
				Expect(doc).To(BeEquivalentTo(data.OrderItem1))
			})
		})
	})

	Context("Delete", func() {
		When("Blank collection", func() {
			It("Should return error", func() {
				err := kvPlugin.Delete("", data.UserKey1)
				Expect(err).Should(HaveOccurred())
			})
		})
		When("Nil key", func() {
			It("Should return error", func() {
				err := kvPlugin.Delete("collection", nil)
				Expect(err).Should(HaveOccurred())
			})
		})
		When("Valid Delete", func() {
			It("Should delete item successfully", func() {
				kvPlugin.Put("users", data.UserKey1, data.UserItem1)

				err := kvPlugin.Delete("users", data.UserKey1)
				Expect(err).ShouldNot(HaveOccurred())

				doc, err := kvPlugin.Get("users", data.UserKey1)
				Expect(doc).To(BeNil())
				Expect(err).Should(HaveOccurred())
			})
		})
		When("Valid Compound Key Delete", func() {
			It("Should delete item successfully", func() {
				kvPlugin.Put("application", data.OrderKey1, data.OrderItem1)

				err := kvPlugin.Delete("application", data.OrderKey1)
				Expect(err).ShouldNot(HaveOccurred())

				doc, err := kvPlugin.Get("application", data.OrderKey1)
				Expect(doc).To(BeNil())
				Expect(err).Should(HaveOccurred())
			})
		})
	})

	Context("Query", func() {
		When("Blank collection argument", func() {
			It("Should return an error", func() {
				vals, err := kvPlugin.Query("", nil, 0)
				Expect(vals).To(BeNil())
				Expect(err).Should(HaveOccurred())
			})
		})
		When("Nil key argument", func() {
			It("Should return an error", func() {
				vals, err := kvPlugin.Query("users", nil, 0)
				Expect(vals).To(BeNil())
				Expect(err).Should(HaveOccurred())
			})
		})
		When("Empty database collection", func() {
			It("Should return empty list", func() {
				vals, err := kvPlugin.Query("users", []sdk.QueryExpression{}, 0)
				Expect(vals).ToNot(BeNil())
				Expect(err).ShouldNot(HaveOccurred())
				Expect(vals).To(HaveLen(0))
			})
		})
		When("Filter users collection", func() {
			It("Should return 1 item", func() {
				kvPlugin.Put("users", data.UserKey1, data.UserItem1)
				kvPlugin.Put("users", data.UserKey2, data.UserItem2)
				kvPlugin.Put("users", data.UserKey3, data.UserItem3)
				exps := []sdk.QueryExpression{
					{Operand: "country", Operator: "==", Value: "US"},
					{Operand: "age", Operator: ">", Value: "40"},
				}
				vals, err := kvPlugin.Query("users", exps, 0)
				Expect(vals).ToNot(BeNil())
				Expect(err).ShouldNot(HaveOccurred())
				Expect(vals).To(HaveLen(1))
				Expect(vals[0]["email"]).To(BeEquivalentTo(data.UserItem3["email"]))
			})
		})
		When("Empty query (Scan)", func() {
			It("Should return all items", func() {
				kvPlugin.Put("application", data.CustomerKey, data.CustomerItem)
				kvPlugin.Put("application", data.OrderKey1, data.OrderItem1)
				kvPlugin.Put("application", data.OrderKey2, data.OrderItem2)
				kvPlugin.Put("application", data.OrderKey3, data.OrderItem3)
				kvPlugin.Put("application", data.ProductKey, data.ProductItem)

				vals, err := kvPlugin.Query("application", []sdk.QueryExpression{}, 0)
				Expect(vals).ToNot(BeNil())
				Expect(err).ShouldNot(HaveOccurred())
				Expect(vals).To(HaveLen(5))
			})
		})
		When("Empty limit query (Scan)", func() {
			It("Should return specified items", func() {
				kvPlugin.Put("application", data.CustomerKey, data.CustomerItem)
				kvPlugin.Put("application", data.OrderKey1, data.OrderItem1)
				kvPlugin.Put("application", data.OrderKey2, data.OrderItem2)
				kvPlugin.Put("application", data.OrderKey3, data.OrderItem3)
				kvPlugin.Put("application", data.ProductKey, data.ProductItem)

				vals, err := kvPlugin.Query("application", []sdk.QueryExpression{}, 3)
				Expect(vals).ToNot(BeNil())
				Expect(err).ShouldNot(HaveOccurred())
				Expect(vals).To(HaveLen(3))

				// DynamoDB scan operations do not have any order, so results could be in any order
			})
		})
		When("PK and SK equality query", func() {
			It("Should return specified item", func() {
				kvPlugin.Put("application", data.CustomerKey, data.CustomerItem)
				kvPlugin.Put("application", data.OrderKey1, data.OrderItem1)
				kvPlugin.Put("application", data.OrderKey2, data.OrderItem2)
				kvPlugin.Put("application", data.OrderKey3, data.OrderItem3)

				exps := []sdk.QueryExpression{
					{Operand: "pk", Operator: "==", Value: "Customer#1000"},
					{Operand: "sk", Operator: "==", Value: "Customer#1000"},
				}
				vals, err := kvPlugin.Query("application", exps, 0)
				Expect(vals).ToNot(BeNil())
				Expect(err).ShouldNot(HaveOccurred())
				Expect(vals).To(HaveLen(1))
				Expect(vals[0]).To(BeEquivalentTo(data.CustomerItem))
			})
		})
		When("PK equality query", func() {
			It("Should return specified items", func() {
				kvPlugin.Put("application", data.CustomerKey, data.CustomerItem)
				kvPlugin.Put("application", data.OrderKey1, data.OrderItem1)
				kvPlugin.Put("application", data.OrderKey2, data.OrderItem2)
				kvPlugin.Put("application", data.OrderKey3, data.OrderItem3)

				exps := []sdk.QueryExpression{
					{Operand: "pk", Operator: "==", Value: "Customer#1000"},
				}
				vals, err := kvPlugin.Query("application", exps, 0)
				Expect(vals).ToNot(BeNil())
				Expect(err).ShouldNot(HaveOccurred())
				Expect(vals).To(HaveLen(4))
				Expect(vals[0]).To(BeEquivalentTo(data.CustomerItem))
				Expect(vals[1]).To(BeEquivalentTo(data.OrderItem1))
				Expect(vals[2]).To(BeEquivalentTo(data.OrderItem2))
				Expect(vals[3]).To(BeEquivalentTo(data.OrderItem3))
			})
		})
		When("PK equality limit query", func() {
			It("Should return specified items", func() {
				kvPlugin.Put("application", data.CustomerKey, data.CustomerItem)
				kvPlugin.Put("application", data.OrderKey1, data.OrderItem1)
				kvPlugin.Put("application", data.OrderKey2, data.OrderItem2)
				kvPlugin.Put("application", data.OrderKey3, data.OrderItem3)

				exps := []sdk.QueryExpression{
					{Operand: "pk", Operator: "==", Value: "Customer#1000"},
				}
				vals, err := kvPlugin.Query("application", exps, 3)
				Expect(vals).ToNot(BeNil())
				Expect(err).ShouldNot(HaveOccurred())
				Expect(vals).To(HaveLen(3))
				Expect(vals[0]).To(BeEquivalentTo(data.CustomerItem))
				Expect(vals[1]).To(BeEquivalentTo(data.OrderItem1))
				Expect(vals[2]).To(BeEquivalentTo(data.OrderItem2))
			})
		})
		When("PK equality and SK startsWith", func() {
			It("Should return specified items", func() {
				kvPlugin.Put("application", data.CustomerKey, data.CustomerItem)
				kvPlugin.Put("application", data.OrderKey1, data.OrderItem1)
				kvPlugin.Put("application", data.OrderKey2, data.OrderItem2)
				kvPlugin.Put("application", data.OrderKey3, data.OrderItem3)
				kvPlugin.Put("application", data.ProductKey, data.ProductItem)

				exps := []sdk.QueryExpression{
					{Operand: "pk", Operator: "==", Value: "Customer#1000"},
					{Operand: "sk", Operator: "startsWith", Value: "Order#"},
				}
				vals, err := kvPlugin.Query("application", exps, 0)
				Expect(vals).ToNot(BeNil())
				Expect(err).ShouldNot(HaveOccurred())
				Expect(vals).To(HaveLen(3))
				Expect(vals[0]).To(BeEquivalentTo(data.OrderItem1))
				Expect(vals[1]).To(BeEquivalentTo(data.OrderItem2))
				Expect(vals[2]).To(BeEquivalentTo(data.OrderItem3))
			})
		})
		When("PK equality and SK >", func() {
			It("Should return specified items", func() {
				kvPlugin.Put("application", data.CustomerKey, data.CustomerItem)
				kvPlugin.Put("application", data.OrderKey1, data.OrderItem1)
				kvPlugin.Put("application", data.OrderKey2, data.OrderItem2)
				kvPlugin.Put("application", data.OrderKey3, data.OrderItem3)
				kvPlugin.Put("application", data.ProductKey, data.ProductItem)

				exps := []sdk.QueryExpression{
					{Operand: "pk", Operator: "==", Value: "Customer#1000"},
					{Operand: "sk", Operator: ">", Value: "Order#501"},
				}
				vals, err := kvPlugin.Query("application", exps, 0)
				Expect(vals).ToNot(BeNil())
				Expect(err).ShouldNot(HaveOccurred())
				Expect(vals).To(HaveLen(2))
				Expect(vals[0]).To(BeEquivalentTo(data.OrderItem2))
				Expect(vals[1]).To(BeEquivalentTo(data.OrderItem3))
			})
		})
		When("PK equality and SK >=", func() {
			It("Should return specified items", func() {
				kvPlugin.Put("application", data.CustomerKey, data.CustomerItem)
				kvPlugin.Put("application", data.OrderKey1, data.OrderItem1)
				kvPlugin.Put("application", data.OrderKey2, data.OrderItem2)
				kvPlugin.Put("application", data.OrderKey3, data.OrderItem3)
				kvPlugin.Put("application", data.ProductKey, data.ProductItem)

				exps := []sdk.QueryExpression{
					{Operand: "pk", Operator: "==", Value: "Customer#1000"},
					{Operand: "sk", Operator: ">=", Value: "Order#501"},
				}
				vals, err := kvPlugin.Query("application", exps, 0)
				Expect(vals).ToNot(BeNil())
				Expect(err).ShouldNot(HaveOccurred())
				Expect(vals).To(HaveLen(3))
				Expect(vals[0]).To(BeEquivalentTo(data.OrderItem1))
				Expect(vals[1]).To(BeEquivalentTo(data.OrderItem2))
				Expect(vals[2]).To(BeEquivalentTo(data.OrderItem3))
			})
		})
		When("PK equality and SK <", func() {
			It("Should return specified items", func() {
				kvPlugin.Put("application", data.CustomerKey, data.CustomerItem)
				kvPlugin.Put("application", data.OrderKey1, data.OrderItem1)
				kvPlugin.Put("application", data.OrderKey2, data.OrderItem2)
				kvPlugin.Put("application", data.OrderKey3, data.OrderItem3)
				kvPlugin.Put("application", data.ProductKey, data.ProductItem)

				exps := []sdk.QueryExpression{
					{Operand: "pk", Operator: "==", Value: "Customer#1000"},
					{Operand: "sk", Operator: "<", Value: "Order#501"},
				}
				vals, err := kvPlugin.Query("application", exps, 0)
				Expect(vals).ToNot(BeNil())
				Expect(err).ShouldNot(HaveOccurred())
				Expect(vals).To(HaveLen(1))
				Expect(vals[0]).To(BeEquivalentTo(data.CustomerItem))
			})
		})
		When("PK equality and SK <=", func() {
			It("Should return specified items", func() {
				kvPlugin.Put("application", data.CustomerKey, data.CustomerItem)
				kvPlugin.Put("application", data.OrderKey1, data.OrderItem1)
				kvPlugin.Put("application", data.OrderKey2, data.OrderItem2)
				kvPlugin.Put("application", data.OrderKey3, data.OrderItem3)
				kvPlugin.Put("application", data.ProductKey, data.ProductItem)

				exps := []sdk.QueryExpression{
					{Operand: "pk", Operator: "==", Value: "Customer#1000"},
					{Operand: "sk", Operator: "<=", Value: "Order#501"},
				}
				vals, err := kvPlugin.Query("application", exps, 0)
				Expect(vals).ToNot(BeNil())
				Expect(err).ShouldNot(HaveOccurred())
				Expect(vals).To(HaveLen(2))
				Expect(vals[0]).To(BeEquivalentTo(data.CustomerItem))
				Expect(vals[1]).To(BeEquivalentTo(data.OrderItem1))
			})
		})
		// Firestore: cant support multiple property inequality operators
		// When("PK equality and SK startsWith and filter", func() {
		// 	It("Should return specified items", func() {
		// 		kvPlugin.Put("application", data.CustomerKey, data.CustomerItem)
		// 		kvPlugin.Put("application", data.OrderKey1, data.OrderItem1)
		// 		kvPlugin.Put("application", data.OrderKey2, data.OrderItem2)
		// 		kvPlugin.Put("application", data.OrderKey3, data.OrderItem3)
		// 		kvPlugin.Put("application", data.ProductKey, data.ProductItem)

		// 		exps := []sdk.QueryExpression{
		// 			{Operand: "pk", Operator: "==", Value: "Customer#1000"},
		// 			{Operand: "sk", Operator: "startsWith", Value: "Order#"},
		// 			{Operand: "number", Operator: ">", Value: "1"},
		// 			{Operand: "price", Operator: "<", Value: "20"},
		// 		}
		// 		vals, err := kvPlugin.Query("application", exps, 0)
		// 		Expect(vals).ToNot(BeNil())
		// 		Expect(err).ShouldNot(HaveOccurred())
		// 		Expect(vals).To(HaveLen(1))
		// 		Expect(vals[0]).To(BeEquivalentTo(data.OrderItem2))
		// 	})
		// })
		// When("PK equality and SK startsWith and between filter", func() {
		// 	It("Should return specified items", func() {
		// 		kvPlugin.Put("application", data.CustomerKey, data.CustomerItem)
		// 		kvPlugin.Put("application", data.OrderKey1, data.OrderItem1)
		// 		kvPlugin.Put("application", data.OrderKey2, data.OrderItem2)
		// 		kvPlugin.Put("application", data.OrderKey3, data.OrderItem3)
		// 		kvPlugin.Put("application", data.ProductKey, data.ProductItem)

		// 		exps := []sdk.QueryExpression{
		// 			{Operand: "pk", Operator: "==", Value: "Customer#1000"},
		// 			{Operand: "sk", Operator: "startsWith", Value: "Order#"},
		// 			{Operand: "number", Operator: ">=", Value: "0"},
		// 			{Operand: "number", Operator: "<=", Value: "1"},
		// 		}
		// 		vals, err := kvPlugin.Query("application", exps, 0)
		// 		Expect(vals).ToNot(BeNil())
		// 		Expect(err).ShouldNot(HaveOccurred())
		// 		Expect(vals).To(HaveLen(1))
		// 		Expect(vals[0]).To(BeEquivalentTo(data.OrderItem1))
		// 	})
		// })
		// When("PK equality and SK startsWith and between filters with reversed order", func() {
		// 	It("Should return specified items", func() {
		// 		kvPlugin.Put("application", data.CustomerKey, data.CustomerItem)
		// 		kvPlugin.Put("application", data.OrderKey1, data.OrderItem1)
		// 		kvPlugin.Put("application", data.OrderKey2, data.OrderItem2)
		// 		kvPlugin.Put("application", data.OrderKey3, data.OrderItem3)
		// 		kvPlugin.Put("application", data.ProductKey, data.ProductItem)

		// 		exps := []sdk.QueryExpression{
		// 			{Operand: "pk", Operator: "==", Value: "Customer#1000"},
		// 			{Operand: "sk", Operator: "startsWith", Value: "Order#"},
		// 			{Operand: "number", Operator: "<=", Value: "1"},
		// 			{Operand: "number", Operator: ">=", Value: "0"},
		// 		}
		// 		vals, err := kvPlugin.Query("application", exps, 0)
		// 		Expect(vals).ToNot(BeNil())
		// 		Expect(err).ShouldNot(HaveOccurred())
		// 		Expect(vals).To(HaveLen(1))
		// 		Expect(vals[0]).To(BeEquivalentTo(data.OrderItem1))
		// 	})
		// })
	})

})
