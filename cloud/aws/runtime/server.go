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
	"github.com/nitrictech/nitric/core/pkg/server"
)

func NewAwsRuntimeServer(resolver resource.AwsResourceResolver, opts ...server.ServerOption) (*server.NitricServer, error) {
	secretPlugin, _ := secret.New(resolver)
	keyValuePlugin, _ := keyvalue.New(resolver)
	topicsPlugin, _ := topic.New(resolver)
	storagePlugin, _ := aws_storage.New(resolver)

	websocketPlugin, _ := websocket.NewAwsApiGatewayWebsocket(resolver)
	queuesPlugin, _ := queue.New(resolver)

	gatewayPlugin := aws_gateway.New(resolver)
	apiPlugin := api.NewAwsApiGatewayProvider(resolver)

	sqlPlugin := sql_service.NewRdsSqlService()

	defaultAwsOpts := []server.ServerOption{
		server.WithKeyValuePlugin(keyValuePlugin),
		server.WithSecretManagerPlugin(secretPlugin),
		server.WithGatewayPlugin(gatewayPlugin),
		server.WithStoragePlugin(storagePlugin),
		server.WithWebsocketPlugin(websocketPlugin),
		server.WithTopicsPlugin(topicsPlugin),
		server.WithQueuesPlugin(queuesPlugin),
		server.WithApiPlugin(apiPlugin),
		server.WithSqlPlugin(sqlPlugin),
	}

	// append overrides
	defaultAwsOpts = append(defaultAwsOpts, opts...)

	return server.New(defaultAwsOpts...)
}
