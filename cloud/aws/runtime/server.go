// Copyright 2021 Nitric Technologies Pty Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
