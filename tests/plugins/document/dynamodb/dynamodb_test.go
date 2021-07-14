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

	ds_plugin "github.com/nitric-dev/membrane/plugins/document/dynamodb"
	test "github.com/nitric-dev/membrane/tests/plugins/document"
	. "github.com/onsi/ginkgo"

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
		createTable(db, "customers")
		createTable(db, "users")
		createTable(db, "items")
		createTable(db, "parentItems")
	})

	AfterEach(func() {
		deleteTable(db, "customers")
		deleteTable(db, "users")
		deleteTable(db, "items")
		deleteTable(db, "parentItems")
	})

	AfterSuite(func() {
		stopDynamoProcess(dynaCmd)
	})

	docPlugin, err := ds_plugin.NewWithClient(db)
	if err != nil {
		panic(err)
	}

	test.GetTests(docPlugin)
	test.SetTests(docPlugin)
	test.DeleteTests(docPlugin)
	test.QueryTests(docPlugin)
})
