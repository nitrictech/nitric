package main

import (
	"log"
	"strconv"

	"github.com/nitric-dev/membrane-plugin-sdk/utils"
	"nitric.io/membrane/membrane"
)

func main() {
	serviceAddress := utils.GetEnv("SERVICE_ADDRESS", "127.0.0.1:50051")
	childAddress := utils.GetEnv("CHILD_ADDRESS", "127.0.0.1:8080")
	pluginDir := utils.GetEnv("PLUGIN_DIR", "./plugins")
	childCommand := utils.GetEnv("INVOKE", "")
	tolerateMissingServices := utils.GetEnv("TOLERATE_MISSING_SERVICES", "false")
	eventingPluginFile := utils.GetEnv("EVENTING_PLUGIN", "eventing.so")
	documentsPluginFile := utils.GetEnv("DOCUMENTS_PLUGIN", "documents.so")
	storagePluginFile := utils.GetEnv("STORAGE_PLUGIN", "storage.so")
	gatewayPluginFile := utils.GetEnv("GATEWAY_PLUGIN", "gateway.so")

	tolerateMissing := strconv.ParseBool(tolerateMissingServices)

	membrane, err := membrane.New(&membrane.MembraneOptions{
		ServiceAddress:      serviceAddress,
		ChildAddress:        childAddress,
		ChildCommand:        childCommand,
		PluginDir:           pluginDir,
		EventingPluginFile:  eventingPluginFile,
		DocumentsPluginFile: documentsPluginFile,
		StoragePluginFile:   storagePluginFile,
		GatewayPluginFile:   gatewayPluginFile,
	})

	if err != nil {
		log.Fatalf("There was an error initialising the membrane server: %v", err)
	}

	// Start the Membrane server
	membrane.Start()
}
