package runtime

import (
	"github.com/nitrictech/nitric/cloud/aws/runtime/api"
	aws_gateway "github.com/nitrictech/nitric/cloud/aws/runtime/gateway"
	"github.com/nitrictech/nitric/cloud/aws/runtime/keyvalue"
	"github.com/nitrictech/nitric/cloud/aws/runtime/queue"
	"github.com/nitrictech/nitric/cloud/aws/runtime/resource"
	"github.com/nitrictech/nitric/cloud/aws/runtime/secret"
	sql_service "github.com/nitrictech/nitric/cloud/aws/runtime/sql"
	aws_storage "github.com/nitrictech/nitric/cloud/aws/runtime/storage"
	"github.com/nitrictech/nitric/cloud/aws/runtime/topic"
	"github.com/nitrictech/nitric/cloud/aws/runtime/websocket"
	"github.com/nitrictech/nitric/core/pkg/membrane"
)

func NewAwsRuntimeServer(resolver resource.AwsResourceResolver, opts ...membrane.RuntimeServerOption) (*membrane.Membrane, error) {
	secretPlugin, _ := secret.New(resolver)
	keyValuePlugin, _ := keyvalue.New(resolver)
	topicsPlugin, _ := topic.New(resolver)
	storagePlugin, _ := aws_storage.New(resolver)

	websocketPlugin, _ := websocket.NewAwsApiGatewayWebsocket(resolver)
	queuesPlugin, _ := queue.New(resolver)

	gatewayPlugin := aws_gateway.New(resolver)
	apiPlugin := api.NewAwsApiGatewayProvider(resolver)

	sqlPlugin := sql_service.NewRdsSqlService()

	defaultAwsMembraneOpts := []membrane.RuntimeServerOption{
		membrane.WithKeyValuePlugin(keyValuePlugin),
		membrane.WithSecretManagerPlugin(secretPlugin),
		membrane.WithGatewayPlugin(gatewayPlugin),
		membrane.WithStoragePlugin(storagePlugin),
		membrane.WithWebsocketPlugin(websocketPlugin),
		membrane.WithTopicsPlugin(topicsPlugin),
		membrane.WithQueuesPlugin(queuesPlugin),
		membrane.WithApiPlugin(apiPlugin),
		membrane.WithSqlPlugin(sqlPlugin),
	}

	// append overrides
	defaultAwsMembraneOpts = append(defaultAwsMembraneOpts, opts...)

	return membrane.New(defaultAwsMembraneOpts...)
}
