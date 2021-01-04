package main

import (
	"fmt"
	"log"
	"plugin"
	"strconv"

	"github.com/nitric-dev/membrane/membrane"
	"github.com/nitric-dev/membrane/plugins/sdk"
	"github.com/nitric-dev/membrane/utils"
)

// Pluggable version of the Nitric membrane
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

	tolerateMissing, err := strconv.ParseBool(tolerateMissingServices)

	if err != nil {
		log.Fatalf("There was an error initialising the membrane server: %v", err)
	}

	var eventingPlugin sdk.EventingPlugin = nil
	var storagePlugin sdk.StoragePlugin = nil
	var documentsPlugin sdk.DocumentsPlugin = nil
	var gatewayPlugin sdk.GatewayPlugin = nil

	// Load the Eventing Plugin
	if plug, err := plugin.Open(fmt.Sprintf("%s/%s", pluginDir, eventingPluginFile)); err == nil {
		if symbol, err := plug.Lookup("New"); err == nil {
			if newFunc, ok := symbol.(func() (sdk.EventingPlugin, error)); ok {
				if plugin, err := newFunc(); err == nil {
					eventingPlugin = plugin
				}
			}
		}
	}

	// Load the Storage Plugin
	if plug, err := plugin.Open(fmt.Sprintf("%s/%s", pluginDir, storagePluginFile)); err == nil {
		if symbol, err := plug.Lookup("New"); err == nil {
			if newFunc, ok := symbol.(func() (sdk.StoragePlugin, error)); ok {
				if plugin, err := newFunc(); err == nil {
					storagePlugin = plugin
				}
			}
		}
	}

	// Load the Documents Plugin
	if plug, err := plugin.Open(fmt.Sprintf("%s/%s", pluginDir, documentsPluginFile)); err == nil {
		if symbol, err := plug.Lookup("New"); err == nil {
			if newFunc, ok := symbol.(func() (sdk.DocumentsPlugin, error)); ok {
				if plugin, err := newFunc(); err == nil {
					documentsPlugin = plugin
				}
			}
		}
	}

	// Load the Gateway Plugin
	if plug, err := plugin.Open(fmt.Sprintf("%s/%s", pluginDir, gatewayPluginFile)); err == nil {
		if symbol, err := plug.Lookup("New"); err == nil {
			if newFunc, ok := symbol.(func() (sdk.GatewayPlugin, error)); ok {
				if plugin, err := newFunc(); err == nil {
					gatewayPlugin = plugin
				}
			}
		}
	}

	membrane, err := membrane.New(&membrane.MembraneOptions{
		ServiceAddress:          serviceAddress,
		ChildAddress:            childAddress,
		ChildCommand:            childCommand,
		EventingPlugin:          eventingPlugin,
		DocumentsPlugin:         documentsPlugin,
		StoragePlugin:           storagePlugin,
		GatewayPlugin:           gatewayPlugin,
		TolerateMissingServices: tolerateMissing,
	})

	if err != nil {
		log.Fatalf("There was an error initialising the membrane server: %v", err)
	}

	// Start the Membrane server
	membrane.Start()
}
