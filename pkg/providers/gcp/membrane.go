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
	"github.com/nitric-dev/membrane/pkg/membrane"
	"github.com/nitric-dev/membrane/pkg/plugins/document/firestore"
	"github.com/nitric-dev/membrane/pkg/plugins/eventing/pubsub"
	"github.com/nitric-dev/membrane/pkg/plugins/gateway/cloudrun"
	"github.com/nitric-dev/membrane/pkg/plugins/queue/pubsub"
	"github.com/nitric-dev/membrane/pkg/plugins/storage/storage"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Setup signal interrupt handling for graceful shutdown
	term := make(chan os.Signal)
	signal.Notify(term, os.Interrupt, syscall.SIGTERM)
	signal.Notify(term, os.Interrupt, syscall.SIGINT)

	documentPlugin, err := firestore_service.New()
	if err != nil {
		fmt.Println("Failed to load document plugin:", err.Error())
	}
	eventingPlugin, err := pubsub_service.New()
	if err != nil {
		fmt.Println("Failed to load eventing plugin:", err.Error())
	}
	storagePlugin, err := storage_service.New()
	if err != nil {
		fmt.Println("Failed to load storage plugin:", err.Error())
	}
	gatewayPlugin, err := cloudrun_plugin.New()
	if err != nil {
		fmt.Println("Failed to load gateway plugin:", err.Error())
	}
	queuePlugin, err := pubsub_queue_service.New()
	if err != nil {
		fmt.Println("Failed to load queue plugin:", err.Error())
	}

	m, err := membrane.New(&membrane.MembraneOptions{
		DocumentPlugin: documentPlugin,
		EventingPlugin: eventingPlugin,
		GatewayPlugin:  gatewayPlugin,
		QueuePlugin:    queuePlugin,
		StoragePlugin:  storagePlugin,
	})

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
		fmt.Println(fmt.Sprintf("Membrane Error: %v, exiting", membraneError))
	case sigTerm := <-term:
		fmt.Println(fmt.Sprintf("Received %v, exiting", sigTerm))
	}

	m.Stop()
}
