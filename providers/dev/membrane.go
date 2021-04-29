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
)

func main() {
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
	m.Start()
}
