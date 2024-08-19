package runtime

import (
	"github.com/nitrictech/nitric/cloud/gcp/runtime/api"
	"github.com/nitrictech/nitric/cloud/gcp/runtime/gateway"
	"github.com/nitrictech/nitric/cloud/gcp/runtime/keyvalue"
	"github.com/nitrictech/nitric/cloud/gcp/runtime/queue"
	"github.com/nitrictech/nitric/cloud/gcp/runtime/resource"
	"github.com/nitrictech/nitric/cloud/gcp/runtime/secret"
	"github.com/nitrictech/nitric/cloud/gcp/runtime/storage"
	"github.com/nitrictech/nitric/cloud/gcp/runtime/topic"
	"github.com/nitrictech/nitric/core/pkg/server"
)

func NewGcpRuntimeServer(resourcesPlugin resource.GcpResourceResolver, opts ...server.ServerOption) (*server.NitricServer, error) {
	secretPlugin, _ := secret.New()
	keyValuePlugin, _ := keyvalue.New()
	topicsPlugin, _ := topic.New(resourcesPlugin)
	storagePlugin, _ := storage.New()

	queuesPlugin, _ := queue.New()

	gatewayPlugin, _ := gateway.New(resourcesPlugin)
	apiPlugin := api.NewGcpApiGatewayProvider(resourcesPlugin)

	defaultGcpOpts := []server.ServerOption{
		server.WithKeyValuePlugin(keyValuePlugin),
		server.WithSecretManagerPlugin(secretPlugin),
		server.WithGatewayPlugin(gatewayPlugin),
		server.WithStoragePlugin(storagePlugin),
		server.WithTopicsPlugin(topicsPlugin),
		server.WithQueuesPlugin(queuesPlugin),
		server.WithApiPlugin(apiPlugin),
	}

	// append overrides
	defaultGcpOpts = append(defaultGcpOpts, opts...)

	return server.New(defaultGcpOpts...)
}
