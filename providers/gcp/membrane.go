package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/nitric-dev/membrane/membrane"
	auth "github.com/nitric-dev/membrane/plugins/auth/identityplatform"
	eventing "github.com/nitric-dev/membrane/plugins/eventing/pubsub"
	gateway "github.com/nitric-dev/membrane/plugins/gateway/cloudrun"
	kv "github.com/nitric-dev/membrane/plugins/kv/firestore"
	queue "github.com/nitric-dev/membrane/plugins/queue/pubsub"
	storage "github.com/nitric-dev/membrane/plugins/storage/storage"
	"github.com/nitric-dev/membrane/utils"
)

func main() {
	serviceAddress := utils.GetEnv("SERVICE_ADDRESS", "127.0.0.1:50051")
	childAddress := utils.GetEnv("CHILD_ADDRESS", "127.0.0.1:8080")
	childCommand := utils.GetEnv("INVOKE", "")
	tolerateMissingServices := utils.GetEnv("TOLERATE_MISSING_SERVICES", "false")

	tolerateMissing, err := strconv.ParseBool(tolerateMissingServices)

	mode, err := membrane.ModeFromString(utils.GetEnv("MEMBRANE_MODE", "FAAS"))
	if err != nil {
		log.Fatalf("There was an error initialising the membrane server: %v", err)
	}

	if err != nil {
		log.Fatalf("There was an error initialising the membrane server: %v", err)
	}

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

	membrane, err := membrane.New(&membrane.MembraneOptions{
		ServiceAddress:          serviceAddress,
		ChildAddress:            childAddress,
		ChildCommand:            childCommand,
		EventingPlugin:          eventingPlugin,
		KvPlugin:                kvPlugin,
		StoragePlugin:           storagePlugin,
		QueuePlugin:             queuePlugin,
		GatewayPlugin:           gatewayPlugin,
		AuthPlugin:              authPlugin,
		TolerateMissingServices: tolerateMissing,
		Mode:                    mode,
	})

	if err != nil {
		log.Fatalf("There was an error initialising the membrane server: %v", err)
	}

	// Start the Membrane server
	membrane.Start()
}
