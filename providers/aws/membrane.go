package main

import (
	"log"
	"strconv"

	"github.com/nitric-dev/membrane/sdk"

	"github.com/nitric-dev/membrane/membrane"
	auth "github.com/nitric-dev/membrane/plugins/auth/cognito"
	eventing "github.com/nitric-dev/membrane/plugins/eventing/sns"
	httpGateway "github.com/nitric-dev/membrane/plugins/gateway/ecs"
	lambdaGateway "github.com/nitric-dev/membrane/plugins/gateway/lambda"
	documents "github.com/nitric-dev/membrane/plugins/kv/dynamodb"
	queue "github.com/nitric-dev/membrane/plugins/queue/sqs"
	storage "github.com/nitric-dev/membrane/plugins/storage/s3"
	"github.com/nitric-dev/membrane/utils"
)

func main() {
	serviceAddress := utils.GetEnv("SERVICE_ADDRESS", "127.0.0.1:50051")
	childAddress := utils.GetEnv("CHILD_ADDRESS", "127.0.0.1:8080")
	childCommand := utils.GetEnv("INVOKE", "")
	tolerateMissingServices := utils.GetEnv("TOLERATE_MISSING_SERVICES", "false")
	gatewayEnv := utils.GetEnv("GATEWAY_ENVIRONMENT", "lambda")

	mode, err := membrane.ModeFromString(utils.GetEnv("MEMBRANE_MODE", "FAAS"))
	if err != nil {
		log.Fatalf("There was an error initialising the membrane server: %v", err)
	}

	tolerateMissing, err := strconv.ParseBool(tolerateMissingServices)

	if err != nil {
		log.Fatalf("There was an error initialising the m server: %v", err)
	}

	eventingPlugin, _ := eventing.New()
	keyValuePlugin, _ := documents.New()
	storagePlugin, _ := storage.New()
	queuePlugin, _ := queue.New()
	authPlugin, _ := auth.New()

	// Load the appropriate gateway based on the environment.
	var gatewayPlugin sdk.GatewayService
	switch gatewayEnv {
	case "lambda":
		gatewayPlugin, _ = lambdaGateway.New()
	default:
		gatewayPlugin, _ = httpGateway.New()
	}

	m, err := membrane.New(&membrane.MembraneOptions{
		ServiceAddress:          serviceAddress,
		ChildAddress:            childAddress,
		ChildCommand:            childCommand,
		EventingPlugin:          eventingPlugin,
		KvPlugin:                keyValuePlugin,
		StoragePlugin:           storagePlugin,
		QueuePlugin:             queuePlugin,
		GatewayPlugin:           gatewayPlugin,
		AuthPlugin:              authPlugin,
		TolerateMissingServices: tolerateMissing,
		Mode:                    mode,
	})

	if err != nil {
		log.Fatalf("There was an error initialising the m server: %v", err)
	}

	// Start the Membrane server
	m.Start()
}