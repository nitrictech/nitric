package runtime

import (
	"github.com/nitrictech/nitric/cloud/azure/runtime/api"
	aws_gateway "github.com/nitrictech/nitric/cloud/azure/runtime/gateway"
	"github.com/nitrictech/nitric/cloud/azure/runtime/keyvalue"
	"github.com/nitrictech/nitric/cloud/azure/runtime/queue"
	"github.com/nitrictech/nitric/cloud/azure/runtime/resource"
	"github.com/nitrictech/nitric/cloud/azure/runtime/secret"
	sql_service "github.com/nitrictech/nitric/cloud/azure/runtime/sql"
	aws_storage "github.com/nitrictech/nitric/cloud/azure/runtime/storage"
	"github.com/nitrictech/nitric/cloud/azure/runtime/topic"
	"github.com/nitrictech/nitric/core/pkg/membrane"
)

func NewAzureRuntimeServer(resourcesPlugin resource.AzResourceResolver, opts ...membrane.RuntimeServerOption) (*membrane.Membrane, error) {
	secretPlugin, _ := secret.New()
	keyValuePlugin, _ := keyvalue.New()
	topicsPlugin, _ := topic.New(resourcesPlugin)
	storagePlugin, _ := aws_storage.New()
	queuesPlugin, _ := queue.New()
	gatewayPlugin, _ := aws_gateway.New(resourcesPlugin)
	apiPlugin := api.NewAzureApiGatewayProvider(resourcesPlugin)

	sqlPlugin, _ := sql_service.New()

	defaultAwsMembraneOpts := []membrane.RuntimeServerOption{
		membrane.WithKeyValuePlugin(keyValuePlugin),
		membrane.WithSecretManagerPlugin(secretPlugin),
		membrane.WithGatewayPlugin(gatewayPlugin),
		membrane.WithStoragePlugin(storagePlugin),
		membrane.WithTopicsPlugin(topicsPlugin),
		membrane.WithQueuesPlugin(queuesPlugin),
		membrane.WithApiPlugin(apiPlugin),
		membrane.WithSqlPlugin(sqlPlugin),
	}

	// append overrides
	defaultAwsMembraneOpts = append(defaultAwsMembraneOpts, opts...)

	return membrane.New(defaultAwsMembraneOpts...)
}
