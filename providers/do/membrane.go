package main

import (
	"log"

	"github.com/nitric-dev/membrane/membrane"
	gateway "github.com/nitric-dev/membrane/plugins/gateway/app_platform"
)

func main() {

	gatewayPlugin, _ := gateway.New()

	m, err := membrane.New(&membrane.MembraneOptions{
		GatewayPlugin:  gatewayPlugin,
		// FIXME: Hardcode as true as we don't have plugins for other services for digital ocean yet...
		TolerateMissingServices: true,
	})

	if err != nil {
		log.Fatalf("There was an error initialising the membrane server: %v", err)
	}

	// Start the Membrane server
	m.Start()
}
