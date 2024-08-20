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
