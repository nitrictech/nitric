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

package mongodb_service_test

import (
	"context"
	"fmt"

	mongodb_service "github.com/nitric-dev/membrane/pkg/plugins/document/mongodb"
	"github.com/nitric-dev/membrane/tests/plugins"
	test "github.com/nitric-dev/membrane/tests/plugins/document"
	. "github.com/onsi/ginkgo"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const containerName = "mongodb-nitric"

func createMongoClient(ctx context.Context) (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017").SetDirect(true)
	client, clientError := mongo.NewClient(clientOptions)

	if clientError != nil {
		return nil, fmt.Errorf("mongodb error creating client: %v", clientError)
	}

	connectError := client.Connect(ctx)

	if connectError != nil {
		return nil, fmt.Errorf("mongodb unable to initialize connection: %v", connectError)
	}

	pingError := client.Ping(ctx, nil)

	if pingError != nil {
		return nil, fmt.Errorf("mongodb unable to connect: %v", pingError)
	}

	return client, nil
}

var _ = Describe("MongoDB", func() {
	defer GinkgoRecover()

	// Start Mongo
	args := []string{
		"docker",
		"run",
		"-d",
		"-p 27017-27019:27017-27019",
		"--name " + containerName,
		"mongo:4.0",
	}
	plugins.StartContainer(containerName, args)

	AfterSuite(func() {
		plugins.StopContainer(containerName)
	})

	ctx := context.Background()

	client, err := createMongoClient(ctx)

	if err != nil {
		panic(err)
	}

	docPlugin := mongodb_service.NewWithClient(client, "testing", ctx)

	if err != nil {
		fmt.Printf("NewClient error: %v \n", err)
		panic(err)
	}

	test.GetTests(docPlugin)
	test.SetTests(docPlugin)
	test.DeleteTests(docPlugin)
	test.QueryTests(docPlugin)
	test.QueryStreamTests(docPlugin)
})
