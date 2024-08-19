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
	"github.com/nitrictech/nitric/core/pkg/membrane"
)

func NewGcpRuntimeServer(resourcesPlugin resource.GcpResourceResolver, opts ...membrane.RuntimeServerOption) (*membrane.Membrane, error) {
	secretPlugin, _ := secret.New()
	keyValuePlugin, _ := keyvalue.New()
	topicsPlugin, _ := topic.New(resourcesPlugin)
	storagePlugin, _ := storage.New()

	queuesPlugin, _ := queue.New()

	gatewayPlugin, _ := gateway.New(resourcesPlugin)
	apiPlugin := api.NewGcpApiGatewayProvider(resourcesPlugin)

	defaultGcpMembraneOpts := []membrane.RuntimeServerOption{
		membrane.WithKeyValuePlugin(keyValuePlugin),
		membrane.WithSecretManagerPlugin(secretPlugin),
		membrane.WithGatewayPlugin(gatewayPlugin),
		membrane.WithStoragePlugin(storagePlugin),
		membrane.WithTopicsPlugin(topicsPlugin),
		membrane.WithQueuesPlugin(queuesPlugin),
		membrane.WithApiPlugin(apiPlugin),
	}

	// append overrides
	defaultGcpMembraneOpts = append(defaultGcpMembraneOpts, opts...)

	return membrane.New(defaultGcpMembraneOpts...)
}
