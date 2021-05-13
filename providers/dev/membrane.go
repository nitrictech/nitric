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
	"github.com/nitric-dev/membrane/membrane"
	auth "github.com/nitric-dev/membrane/plugins/auth/dev"
	eventing "github.com/nitric-dev/membrane/plugins/eventing/dev"
	gateway "github.com/nitric-dev/membrane/plugins/gateway/dev"
	kv "github.com/nitric-dev/membrane/plugins/kv/dev"
	queue "github.com/nitric-dev/membrane/plugins/queue/dev"
	storage "github.com/nitric-dev/membrane/plugins/storage/dev"
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

	eventingPlugin, _ := eventing.New()
	kvPlugin, _ := kv.New()
	storagePlugin, _ := storage.New()
	gatewayPlugin, _ := gateway.New()
	queuePlugin, _ := queue.New()
	authPlugin, _ := auth.New()

	m, err := membrane.New(&membrane.MembraneOptions{
		EventingPlugin:          eventingPlugin,
		KvPlugin:                kvPlugin,
		StoragePlugin:           storagePlugin,
		QueuePlugin:             queuePlugin,
		AuthPlugin:              authPlugin,
		GatewayPlugin:           gatewayPlugin,
	})

	if err != nil {
		log.Fatalf("There was an error initialising the membraneServer server: %v", err)
	}

	// Start the Membrane server
	println("starting server")
	go(m.Start)()

	println("wait for term signal")
	// Wait for a terminate interrupt
	<-term
	println("stopping server")
	m.Stop()
}
