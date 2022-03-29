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

package firestore_service_test

import (
	"context"
	"fmt"
	"os"

	firestore_service "github.com/nitrictech/nitric/pkg/plugins/document/firestore"

	"cloud.google.com/go/firestore"
	. "github.com/onsi/ginkgo"

	"github.com/nitrictech/nitric/tests/plugins"
	test "github.com/nitrictech/nitric/tests/plugins/document"
)

const containerName = "firestore-nitric"
const port = "8080"

func createFirestoreClient(ctx context.Context) *firestore.Client {
	client, err := firestore.NewClient(ctx, "test")
	if err != nil {
		fmt.Printf("NewClient error: %v \n", err)
		panic(err)
	}

	return client
}

var _ = Describe("Firestore", func() {
	defer GinkgoRecover()

	// Start Local DynamoDB
	os.Setenv("FIRESTORE_EMULATOR_HOST", "localhost:"+port)

	// Start Firestore Emulator
	args := []string{
		"docker",
		"run",
		"-d",
		"-p " + port + ":" + port,
		"--env \"FIRESTORE_PROJECT_ID=dummy-project-id\"",
		"--name " + containerName,
		"mtlynch/firestore-emulator-docker",
	}
	plugins.StartContainer(containerName, args)

	ctx := context.Background()
	db := createFirestoreClient(ctx)

	AfterSuite(func() {
		plugins.StopContainer(containerName)
	})

	docPlugin, err := firestore_service.NewWithClient(db, ctx)
	if err != nil {
		panic(err)
	}

	test.GetTests(docPlugin)
	test.SetTests(docPlugin)
	test.DeleteTests(docPlugin)
	test.QueryTests(docPlugin)
	test.QueryStreamTests(docPlugin)
})
