package main

import (
	"fmt"
	"github.com/nitric-dev/membrane/membrane"
	auth "github.com/nitric-dev/membrane/plugins/auth/identityplatform"
	eventing "github.com/nitric-dev/membrane/plugins/eventing/pubsub"
	gateway "github.com/nitric-dev/membrane/plugins/gateway/cloudrun"
	kv "github.com/nitric-dev/membrane/plugins/kv/firestore"
	queue "github.com/nitric-dev/membrane/plugins/queue/pubsub"
	storage "github.com/nitric-dev/membrane/plugins/storage/storage"
	"log"
)

func main() {
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
		EventingPlugin:          eventingPlugin,
		KvPlugin:                kvPlugin,
		StoragePlugin:           storagePlugin,
		QueuePlugin:             queuePlugin,
		GatewayPlugin:           gatewayPlugin,
		AuthPlugin:              authPlugin,
	})

	if err != nil {
		log.Fatalf("There was an error initialising the membrane server: %v", err)
	}

	// Start the Membrane server
	m.Start()
}
