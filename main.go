package main

import (
	"github.com/nitric-dev/membrane-plugin-sdk/utils"
	"nitric.io/membrane/membrane"
)

func main() {
	serviceAddress := utils.GetEnv("SERVICE_ADDRESS", "127.0.0.1:50051")
	childAddress := utils.GetEnv("CHILD_ADDRESS", "127.0.0.1:8080")
	pluginDir := utils.GetEnv("PLUGIN_DIR", "./plugins")
	childCommand := utils.GetEnv("INVOKE", "")
	eventingPluginFile := utils.GetEnv("EVENTING_PLUGIN", "eventing.so")
	documentsPluginFile := utils.GetEnv("DOCUMENTS_PLUGIN", "documents.so")
	storagePluginFile := utils.GetEnv("STORAGE_PLUGIN", "storage.so")
	gatewayPluginFile := utils.GetEnv("GATEWAY_PLUGIN", "gateway.so")

	membrane, error := membrane.New(
		serviceAddress,
		childAddress,
		childCommand,
		pluginDir,
		eventingPluginFile,
		documentsPluginFile,
		storagePluginFile,
		gatewayPluginFile,
	)

	if error != nil {
		panic(error)
	}

	// Start the Membrane server
	membrane.Start()
}
