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
	"github.com/nitric-dev/membrane/pkg/plugins/document/boltdb"
	"github.com/nitric-dev/membrane/pkg/plugins/eventing/dev"
	"github.com/nitric-dev/membrane/pkg/plugins/gateway/dev"
	"github.com/nitric-dev/membrane/pkg/plugins/queue/dev"
	"github.com/nitric-dev/membrane/pkg/plugins/storage/boltdb"
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

	documentPlugin, _ := boltdb_service.New()
	eventingPlugin, _ := eventing_service.New()
	gatewayPlugin, _ := gateway_plugin.New()
	queuePlugin, _ := queue_service.New()
	storagePlugin, _ := boltdb_storage_service.New()

	m, err := membrane.New(&membrane.MembraneOptions{
		DocumentPlugin: documentPlugin,
		EventingPlugin: eventingPlugin,
		GatewayPlugin:  gatewayPlugin,
		QueuePlugin:    queuePlugin,
		StoragePlugin:  storagePlugin,
	})

	if err != nil {
		log.Fatalf("There was an error initialising the membraneServer server: %v", err)
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