package main

import (
	http_service "github.com/nitric-dev/membrane/plugins/gateway/appservice"
	"github.com/nitric-dev/membrane/sdk"
	"log"

	"github.com/nitric-dev/membrane/membrane"
)

func main() {

	authPlugin := &sdk.UnimplementedAuthPlugin{}
	kvPlugin := &sdk.UnimplementedKeyValuePlugin{}
	eventingPlugin := &sdk.UnimplementedEventingPlugin{}
	gatewayPlugin, _ := http_service.New()
	storagePlugin := &sdk.UnimplementedStoragePlugin{}
	queuePlugin := &sdk.UnimplementedQueuePlugin{}

	m, err := membrane.New(&membrane.MembraneOptions{
		AuthPlugin:              authPlugin,
		KvPlugin:                kvPlugin,
		EventingPlugin:          eventingPlugin,
		GatewayPlugin:           gatewayPlugin,
		StoragePlugin:           storagePlugin,
		QueuePlugin:             queuePlugin,
	})

	if err != nil {
		log.Fatalf("There was an error initialising the m server: %v", err)
	}

	// Start the Membrane server
	m.Start()
}
