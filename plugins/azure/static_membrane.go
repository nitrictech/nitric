package main

import (
	"github.com/nitric-dev/membrane/plugins/sdk"
	"log"
	"strconv"

	"github.com/nitric-dev/membrane/membrane"
	"github.com/nitric-dev/membrane/utils"
)

func main() {
	serviceAddress := utils.GetEnv("SERVICE_ADDRESS", "127.0.0.1:50051")
	childAddress := utils.GetEnv("CHILD_ADDRESS", "127.0.0.1:8080")
	childCommand := utils.GetEnv("INVOKE", "")
	tolerateMissingServices := utils.GetEnv("TOLERATE_MISSING_SERVICES", "false")

	tolerateMissing, err := strconv.ParseBool(tolerateMissingServices)

	if err != nil {
		log.Fatalf("There was an error initialising the m server: %v", err)
	}

	//eventingPlugin, _ := eventing.New()
	//documentsPlugin, _ := documents.New()
	//storagePlugin, _ := storage.New()
	//queuePlugin, _ := queue.New()
	//authPlugin, _ := auth.New()
	//gatewayPlugin, _ := httpGateway.New()
	authPlugin := &sdk.UnimplementedAuthPlugin{}
	documentsPlugin := &sdk.UnimplementedDocumentsPlugin{}
	eventingPlugin := &sdk.UnimplementedEventingPlugin{}
	gatewayPlugin := &sdk.UnimplementedGatewayPlugin{}
	storagePlugin := &sdk.UnimplementedStoragePlugin{}
	queuePlugin := &sdk.UnimplementedQueuePlugin{}



	m, err := membrane.New(&membrane.MembraneOptions{
		ServiceAddress:          serviceAddress,
		ChildAddress:            childAddress,
		ChildCommand:            childCommand,
		TolerateMissingServices: tolerateMissing,
		AuthPlugin:              authPlugin,
		DocumentsPlugin:         documentsPlugin,
		EventingPlugin:          eventingPlugin,
		GatewayPlugin:           gatewayPlugin,
		StoragePlugin:           storagePlugin,
		QueuePlugin:             queuePlugin,
	})

	if err != nil {
		log.Fatalf("There was an error initialising the m server: %v", err)
	}

	// Start the Membrane server
	_ = m.Start()
}
