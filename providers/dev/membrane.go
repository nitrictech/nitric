package main

import (
	"log"
	"strconv"

	"github.com/nitric-dev/membrane/membrane"
	auth "github.com/nitric-dev/membrane/plugins/auth/dev"
	eventing "github.com/nitric-dev/membrane/plugins/eventing/dev"
	gateway "github.com/nitric-dev/membrane/plugins/gateway/dev"
	kv "github.com/nitric-dev/membrane/plugins/kv/dev"
	queue "github.com/nitric-dev/membrane/plugins/queue/dev"
	storage "github.com/nitric-dev/membrane/plugins/storage/dev"
	"github.com/nitric-dev/membrane/utils"
)

func main() {
	serviceAddress := utils.GetEnv("SERVICE_ADDRESS", "127.0.0.1:50051")
	childAddress := utils.GetEnv("CHILD_ADDRESS", "127.0.0.1:8080")
	childCommand := utils.GetEnv("INVOKE", "")
	tolerateMissingServices := utils.GetEnv("TOLERATE_MISSING_SERVICES", "false")
	membraneMode := utils.GetEnv("MEMBRANE_MODE", "FAAS")

	mode, err := membrane.ModeFromString(membraneMode)
	if err != nil {
		log.Fatalf("There was an error initialising the membrane server: %v", err)
	}

	tolerateMissing, err := strconv.ParseBool(tolerateMissingServices)
	if err != nil {
		log.Fatalf("There was an error initialising the membrane server: %v", err)
	}

	eventingPlugin, _ := eventing.New()
	kvPlugin, _ := kv.New()
	storagePlugin, _ := storage.New()
	gatewayPlugin, _ := gateway.New()
	queuePlugin, _ := queue.New()
	authPlugin, _ := auth.New()

	membraneServer, err := membrane.New(&membrane.MembraneOptions{
		ServiceAddress:          serviceAddress,
		ChildAddress:            childAddress,
		ChildCommand:            childCommand,
		EventingPlugin:          eventingPlugin,
		KvPlugin:                kvPlugin,
		StoragePlugin:           storagePlugin,
		QueuePlugin:             queuePlugin,
		AuthPlugin:              authPlugin,
		GatewayPlugin:           gatewayPlugin,
		TolerateMissingServices: tolerateMissing,
		Mode:                    mode,
	})

	if err != nil {
		log.Fatalf("There was an error initialising the membraneServer server: %v", err)
	}

	// Start the Membrane server
	membraneServer.Start()
}
