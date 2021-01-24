package main

import (
	"log"
	"strconv"

	"github.com/nitric-dev/membrane/membrane"
	documents "github.com/nitric-dev/membrane/plugins/aws/documents/dynamodb"
	eventing "github.com/nitric-dev/membrane/plugins/aws/eventing/sns"
	gateway "github.com/nitric-dev/membrane/plugins/aws/gateway/lambda"
	queue "github.com/nitric-dev/membrane/plugins/aws/queue/sqs"
	storage "github.com/nitric-dev/membrane/plugins/aws/storage/s3"
	"github.com/nitric-dev/membrane/utils"
)

func main() {
	serviceAddress := utils.GetEnv("SERVICE_ADDRESS", "127.0.0.1:50051")
	childAddress := utils.GetEnv("CHILD_ADDRESS", "127.0.0.1:8080")
	childCommand := utils.GetEnv("INVOKE", "")
	tolerateMissingServices := utils.GetEnv("TOLERATE_MISSING_SERVICES", "false")

	tolerateMissing, err := strconv.ParseBool(tolerateMissingServices)

	if err != nil {
		log.Fatalf("There was an error initialising the membrane server: %v", err)
	}

	eventingPlugin, _ := eventing.New()
	documentsPlugin, _ := documents.New()
	storagePlugin, _ := storage.New()
	gatewayPlugin, _ := gateway.New()
	queuePlugin, _ := queue.New()

	membrane, err := membrane.New(&membrane.MembraneOptions{
		ServiceAddress:          serviceAddress,
		ChildAddress:            childAddress,
		ChildCommand:            childCommand,
		EventingPlugin:          eventingPlugin,
		DocumentsPlugin:         documentsPlugin,
		StoragePlugin:           storagePlugin,
		QueuePlugin:             queuePlugin,
		GatewayPlugin:           gatewayPlugin,
		TolerateMissingServices: tolerateMissing,
	})

	if err != nil {
		log.Fatalf("There was an error initialising the membrane server: %v", err)
	}

	// Start the Membrane server
	membrane.Start()
}
