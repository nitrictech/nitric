package runtime

import (
	"github.com/nitrictech/nitric/cloud/azure/runtime/api"
	az_gateway "github.com/nitrictech/nitric/cloud/azure/runtime/gateway"
	"github.com/nitrictech/nitric/cloud/azure/runtime/keyvalue"
	"github.com/nitrictech/nitric/cloud/azure/runtime/queue"
	"github.com/nitrictech/nitric/cloud/azure/runtime/resource"
	"github.com/nitrictech/nitric/cloud/azure/runtime/secret"
	sql_service "github.com/nitrictech/nitric/cloud/azure/runtime/sql"
	az_storage "github.com/nitrictech/nitric/cloud/azure/runtime/storage"
	"github.com/nitrictech/nitric/cloud/azure/runtime/topic"
	"github.com/nitrictech/nitric/core/pkg/server"
)

func NewAzureRuntimeServer(resourcesPlugin resource.AzResourceResolver, opts ...server.ServerOption) (*server.NitricServer, error) {
	secretPlugin, _ := secret.New()
	keyValuePlugin, _ := keyvalue.New()
	topicsPlugin, _ := topic.New(resourcesPlugin)
	storagePlugin, _ := az_storage.New()
	queuesPlugin, _ := queue.New()
	gatewayPlugin, _ := az_gateway.New(resourcesPlugin)
	apiPlugin := api.NewAzureApiGatewayProvider(resourcesPlugin)

	sqlPlugin, _ := sql_service.New()

	defaultAzureOpts := []server.ServerOption{
		server.WithKeyValuePlugin(keyValuePlugin),
		server.WithSecretManagerPlugin(secretPlugin),
		server.WithGatewayPlugin(gatewayPlugin),
		server.WithStoragePlugin(storagePlugin),
		server.WithTopicsPlugin(topicsPlugin),
		server.WithQueuesPlugin(queuesPlugin),
		server.WithApiPlugin(apiPlugin),
		server.WithSqlPlugin(sqlPlugin),
	}

	// append overrides
	defaultAzureOpts = append(defaultAzureOpts, opts...)

	return server.New(defaultAzureOpts...)
}
