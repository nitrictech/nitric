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
	"context"
	"fmt"
	"os"
	"time"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	mock_provider "github.com/nitrictech/nitric/core/mocks/provider"
	dynamodb_service "github.com/nitrictech/nitric/core/pkg/plugins/document/dynamodb"
	"github.com/nitrictech/nitric/core/pkg/providers/aws/core"
	"github.com/nitrictech/nitric/core/tests/plugins"
	test "github.com/nitrictech/nitric/core/tests/plugins/document"
)

const (
	containerName = "dynamodb-nitric"
	port          = "8000"
)

var _ = Describe("DynamoDb", func() {
	defer GinkgoRecover()

	os.Setenv("AWS_ACCESS_KEY_ID", "fakeMyKeyId")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "fakeSecretAccessKey")
	os.Setenv("AWS_REGION", "X")

	// Start Local DynamoDB
	// Run dynamodb container
	args := []string{
		"docker",
		"run",
		"-d",
		"-p " + port + ":" + port,
		"--name " + containerName,
		"amazon/dynamodb-local:latest",
	}
	plugins.StartContainer(containerName, args)

	// // Create DynamoDB client
	db := createDynamoClient()

	testConnection(db)

	BeforeEach(func() {
		// Table names suffixed with 7 alphanumeric chars to match pulumi deployment.
		createTable(db, "customers-1111111")
		createTable(db, "users-1111111")
		createTable(db, "items-1111111")
		createTable(db, "parentItems-1111111")
	})

	AfterEach(func() {
		deleteTable(db, "customers-1111111")
		deleteTable(db, "users-1111111")
		deleteTable(db, "items-1111111")
		deleteTable(db, "parentItems-1111111")
	})

	AfterSuite(func() {
		plugins.StopContainer(containerName)
	})

	ctrl := gomock.NewController(GinkgoT())
	provider := mock_provider.NewMockAwsProvider(ctrl)

	provider.EXPECT().GetResources(gomock.Any(), core.AwsResource_Collection).AnyTimes().Return(map[string]string{
		"customers":   "arn:${Partition}:dynamodb:${Region}:${Account}:table/customers-1111111",
		"users":       "arn:${Partition}:dynamodb:${Region}:${Account}:table/users-1111111",
		"items":       "arn:${Partition}:dynamodb:${Region}:${Account}:table/items-1111111",
		"parentItems": "arn:${Partition}:dynamodb:${Region}:${Account}:table/parentItems-1111111",
	}, nil)

	docPlugin, err := dynamodb_service.NewWithClient(provider, db)
	if err != nil {
		panic(err)
	}

	test.GetTests(docPlugin)
	test.SetTests(docPlugin)
	test.DeleteTests(docPlugin)
	test.QueryTests(docPlugin)
	test.QueryStreamTests(docPlugin)
})

func createDynamoClient() *dynamodb.Client {
	cfg, sessionError := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("x"),
		config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			return aws.Endpoint{URL: "http://localhost:8000"}, nil
		})),
	)
	if sessionError != nil {
		return nil
	}

	return dynamodb.NewFromConfig(cfg)
}

func testConnection(db *dynamodb.Client) {
	input := &dynamodb.ListTablesInput{}

	if _, err := db.ListTables(context.TODO(), input); err != nil {
		// Wait for Java DynamoDB process to get started
		time.Sleep(2 * time.Second)

		if _, err := db.ListTables(context.TODO(), input); err != nil {
			time.Sleep(4 * time.Second)

			if _, err := db.ListTables(context.TODO(), input); err != nil {
				fmt.Printf("DynamoDB connection error: %v \n", err)
				panic(err)
			}
		} else {
			return
		}
	}
}

func createTable(db *dynamodb.Client, tableName string) {
	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("_pk"),
				AttributeType: "S",
			},
			{
				AttributeName: aws.String("_sk"),
				AttributeType: "S",
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("_pk"),
				KeyType:       types.KeyTypeHash,
			},
			{
				AttributeName: aws.String("_sk"),
				KeyType:       types.KeyTypeRange,
			},
		},
		ProvisionedThroughput: &types.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
		TableName: aws.String(tableName),
		Tags: []types.Tag{
			{
				Key:   aws.String("x-nitric-name"),
				Value: aws.String(tableName),
			},
		},
	}
	_, err := db.CreateTable(context.TODO(), input)
	if err != nil {
		panic(fmt.Sprintf("Error calling CreateTable: %s", err))
	}
}

func deleteTable(db *dynamodb.Client, tableName string) {
	deleteInput := &dynamodb.DeleteTableInput{
		TableName: aws.String(tableName),
	}

	_, err := db.DeleteTable(context.TODO(), deleteInput)
	if err != nil {
		panic(fmt.Sprintf("Error calling DeleteTable: %s", err))
	}
}
