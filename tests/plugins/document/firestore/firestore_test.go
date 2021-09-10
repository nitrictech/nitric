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
	"os/exec"
	"strings"
	"syscall"

	firestore_service "github.com/nitric-dev/membrane/pkg/plugins/document/firestore"

	"cloud.google.com/go/firestore"
	test "github.com/nitric-dev/membrane/tests/plugins/document"
	. "github.com/onsi/ginkgo"
)

const shell = "/bin/sh"
const containerName = "firestore-nitric"
const port = "8080"

func startFirestoreContainer() {
	// Run dynamodb container
	args := []string{
		"docker",
		"run",
		"-d",
		"-p " + port + ":" + port,
		"--env \"FIRESTORE_PROJECT_ID=dummy-project-id\"",
		"--name " + containerName,
		"mtlynch/firestore-emulator-docker",
	}

	cmd := exec.Command("/bin/sh", "-c", strings.Join(args[:], " "))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error running Firestore Image %v : %v \n", cmd, err)
		panic(fmt.Sprintf("Error running Firestore Image %v : %v", cmd, err))
	}

	// Makes process killable
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
}

func stopFirestoreContainer() {
	// clean up
	stopArgs := []string{
		"docker",
		"container",
		"stop",
		containerName,
	}

	stopCmd := exec.Command(shell, "-c", strings.Join(stopArgs[:], " "))

	if err := stopCmd.Run(); err != nil {
		fmt.Printf("Error stopping Firestore container %v : %v \n", stopCmd, err)
		panic(fmt.Sprintf("Error stopping Firestore container %v : %v", stopCmd, err))
	}

	removeArgs := []string{
		"docker",
		"container",
		"rm",
		containerName,
	}

	removeCmd := exec.Command(shell, "-c", strings.Join(removeArgs[:], " "))

	if err := removeCmd.Run(); err != nil {
		fmt.Printf("Error removing Firestore container %v : %v \n", removeCmd, err)
		panic(fmt.Sprintf("Error removing Firestore container %v : %v", removeCmd, err))
	}
}

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
	startFirestoreContainer()

	ctx := context.Background()
	db := createFirestoreClient(ctx)

	AfterSuite(func() {
		stopFirestoreContainer()
	})

	docPlugin, err := firestore_service.NewWithClient(db, ctx)
	if err != nil {
		panic(err)
	}

	test.GetTests(docPlugin)
	test.SetTests(docPlugin)
	test.DeleteTests(docPlugin)
	test.QueryTests(docPlugin)
})
