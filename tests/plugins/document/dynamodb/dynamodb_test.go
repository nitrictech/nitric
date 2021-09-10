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

	dynamodb_service "github.com/nitric-dev/membrane/pkg/plugins/document/dynamodb"

	"github.com/nitric-dev/membrane/tests/plugins"
	test "github.com/nitric-dev/membrane/tests/plugins/document"
	. "github.com/onsi/ginkgo"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

const shell = "/bin/sh"
const containerName = "dynamodb-nitric"
const port = "8000"

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

	// Create DynamoDB client
	db := createDynamoClient()

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

	docPlugin, err := dynamodb_service.NewWithClient(db)
	if err != nil {
		panic(err)
	}

	test.GetTests(docPlugin)
	test.SetTests(docPlugin)
	test.DeleteTests(docPlugin)
	test.QueryTests(docPlugin)
})

func createDynamoClient() *dynamodb.DynamoDB {
	sess := session.Must(session.NewSession(&aws.Config{
		Region:   aws.String("x"),
		Endpoint: aws.String("http://localhost:" + port),
	}))

	return dynamodb.New(sess)
}

func createTable(db *dynamodb.DynamoDB, tableName string) {
	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("_pk"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("_sk"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("_pk"),
				KeyType:       aws.String("HASH"),
			},
			{
				AttributeName: aws.String("_sk"),
				KeyType:       aws.String("RANGE"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
		TableName: aws.String(tableName),
		Tags: []*dynamodb.Tag{
			{
				Key:   aws.String("x-nitric-name"),
				Value: aws.String(tableName),
			},
		},
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
