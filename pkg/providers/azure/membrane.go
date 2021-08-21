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

	azure_storage_queue_service "github.com/nitric-dev/membrane/pkg/plugins/queue/azure_storage"

	"github.com/nitric-dev/membrane/pkg/membrane"
	"github.com/nitric-dev/membrane/pkg/plugins/document"
	"github.com/nitric-dev/membrane/pkg/plugins/events"
	http_service "github.com/nitric-dev/membrane/pkg/plugins/gateway/appservice"
	"github.com/nitric-dev/membrane/pkg/plugins/queue"
	azblob_service "github.com/nitric-dev/membrane/pkg/plugins/storage/azblob"
	"github.com/nitric-dev/membrane/pkg/plugins/storage"
)

func main() {
	term := make(chan os.Signal)
	signal.Notify(term, os.Interrupt, syscall.SIGTERM)
	signal.Notify(term, os.Interrupt, syscall.SIGINT)

	documentPlugin := &document.UnimplementedDocumentPlugin{}
	eventsPlugin := &events.UnimplementedeventsPlugin{}
	gatewayPlugin, _ := http_service.New()
	queuePlugin, _ := azure_storage_queue_service.New()
	storagePlugin, _ := azblob_service.New()

	m, err := membrane.New(&membrane.MembraneOptions{
		DocumentPlugin: documentPlugin,
		EventsPlugin:   eventsPlugin,
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
		fmt.Printf("Membrane Error: %v, exiting\n", membraneError)
	case sigTerm := <-term:
		fmt.Printf("Received %v, exiting\n", sigTerm)
	}

	m.Stop()
}
