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

package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/nitrictech/nitric/pkg/membrane"
	firestore_service "github.com/nitrictech/nitric/pkg/plugins/document/firestore"
	pubsub_service "github.com/nitrictech/nitric/pkg/plugins/events/pubsub"
	cloudrun_plugin "github.com/nitrictech/nitric/pkg/plugins/gateway/cloudrun"
	pubsub_queue_service "github.com/nitrictech/nitric/pkg/plugins/queue/pubsub"
	secret_manager_secret_service "github.com/nitrictech/nitric/pkg/plugins/secret/secret_manager"
	storage_service "github.com/nitrictech/nitric/pkg/plugins/storage/storage"
)

func main() {
	// Setup signal interrupt handling for graceful shutdown
	var err error
	term := make(chan os.Signal, 1)
	signal.Notify(term, os.Interrupt, syscall.SIGTERM)
	signal.Notify(term, os.Interrupt, syscall.SIGINT)

	membraneOpts := membrane.DefaultMembraneOptions()

	membraneOpts.SecretPlugin, err = secret_manager_secret_service.New()
	if err != nil {
		log.Default().Println("Failed to load secret plugin:", err.Error())
	}

	membraneOpts.DocumentPlugin, err = firestore_service.New()
	if err != nil {
		log.Default().Println("Failed to load document plugin:", err.Error())
	}

	membraneOpts.EventsPlugin, err = pubsub_service.New()
	if err != nil {
		log.Default().Println("Failed to load events plugin:", err.Error())
	}

	membraneOpts.StoragePlugin, err = storage_service.New()
	if err != nil {
		log.Default().Println("Failed to load storage plugin:", err.Error())
	}

	membraneOpts.GatewayPlugin, err = cloudrun_plugin.New()
	if err != nil {
		log.Default().Println("Failed to load gateway plugin:", err.Error())
	}

	membraneOpts.QueuePlugin, err = pubsub_queue_service.New()
	if err != nil {
		log.Default().Println("Failed to load queue plugin:", err.Error())
	}

	m, err := membrane.New(membraneOpts)

	if err != nil {
		log.Fatalf("There was an error initialising the membrane server: %v", err)
	}

	errChan := make(chan error)
	// Start the Membrane server
	go func(chan error) {
		errChan <- m.Start()
	}(errChan)

	select {
	case membraneError := <-errChan:
		log.Default().Println(fmt.Sprintf("Membrane Error: %v, exiting", membraneError))
	case sigTerm := <-term:
		log.Default().Println(fmt.Sprintf("Received %v, exiting", sigTerm))
	}

	m.Stop()
}
