package main

import (
	"log"

	"github.com/nitric-dev/membrane/membrane"
	gateway "github.com/nitric-dev/membrane/plugins/do/gateway/http"
	"github.com/nitric-dev/membrane/utils"
)

func main() {
	serviceAddress := utils.GetEnv("SERVICE_ADDRESS", "127.0.0.1:50051")
	childAddress := utils.GetEnv("CHILD_ADDRESS", "127.0.0.1:8080")
	childCommand := utils.GetEnv("INVOKE", "")
	// tolerateMissingServices := utils.GetEnv("TOLERATE_MISSING_SERVICES", "false")
	membraneMode := utils.GetEnv("MEMBRANE_MODE", "FAAS")

	mode, err := membrane.ModeFromString(membraneMode)
	if err != nil {
		log.Fatalf("There was an error initialising the membrane server: %v", err)
	}

	//tolerateMissing, err := strconv.ParseBool(tolerateMissingServices)
	//if err != nil {
	//	log.Fatalf("There was an error initialising the membrane server: %v", err)
	//}

	gatewayPlugin, _ := gateway.New()

	membrane, err := membrane.New(&membrane.MembraneOptions{
		ServiceAddress: serviceAddress,
		ChildAddress:   childAddress,
		ChildCommand:   childCommand,
		GatewayPlugin:  gatewayPlugin,
		// Hardcode as true as we don't have plugins
		// for other services for digital ocean yet...
		TolerateMissingServices: true,
		Mode:                    mode,
	})

	if err != nil {
		log.Fatalf("There was an error initialising the membrane server: %v", err)
	}

	// Start the Membrane server
	membrane.Start()
}
