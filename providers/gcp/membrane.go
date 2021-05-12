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

	"github.com/nitric-dev/membrane/membrane"
	auth "github.com/nitric-dev/membrane/plugins/auth/identityplatform"
	eventing "github.com/nitric-dev/membrane/plugins/eventing/pubsub"
	gateway "github.com/nitric-dev/membrane/plugins/gateway/cloudrun"
	kv "github.com/nitric-dev/membrane/plugins/kv/firestore"
	queue "github.com/nitric-dev/membrane/plugins/queue/pubsub"
	storage "github.com/nitric-dev/membrane/plugins/storage/storage"
)

func main() {
	// Setup signal interrupt handling for graceful shutdown
	term := make(chan os.Signal)
	signal.Notify(term, os.Interrupt, syscall.SIGTERM)
	signal.Notify(term, os.Interrupt, syscall.SIGINT)

	eventingPlugin, err := eventing.New()
	if err != nil {
		fmt.Println("Failed to load eventing plugin:", err.Error())
	}
	kvPlugin, err := kv.New()
	if err != nil {
		fmt.Println("Failed to load kv plugin:", err.Error())
	}
	storagePlugin, err := storage.New()
	if err != nil {
		fmt.Println("Failed to load storage plugin:", err.Error())
	}
	gatewayPlugin, err := gateway.New()
	if err != nil {
		fmt.Println("Failed to load gateway plugin:", err.Error())
	}
	queuePlugin, err := queue.New()
	if err != nil {
		fmt.Println("Failed to load queue plugin:", err.Error())
	}
	authPlugin, err := auth.New()
	if err != nil {
		fmt.Println("Failed to load auth plugin:", err.Error())
	}

	m, err := membrane.New(&membrane.MembraneOptions{
		EventingPlugin: eventingPlugin,
		KvPlugin:       kvPlugin,
		StoragePlugin:  storagePlugin,
		QueuePlugin:    queuePlugin,
		GatewayPlugin:  gatewayPlugin,
		AuthPlugin:     authPlugin,
	})

	if err != nil {
		log.Fatalf("There was an error initialising the membrane server: %v", err)
	}
	// Start the Membrane server
	println("starting server")
	go (m.Start)()

	println("wait for term signal")
	// Wait for a terminate interrupt
	<-term
	println("stopping server")
	m.Stop()
}
