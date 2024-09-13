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
	"github.com/nitrictech/nitric/cloud/gcp/runtime/api"
	"github.com/nitrictech/nitric/cloud/gcp/runtime/gateway"
	"github.com/nitrictech/nitric/cloud/gcp/runtime/keyvalue"
	"github.com/nitrictech/nitric/cloud/gcp/runtime/queue"
	"github.com/nitrictech/nitric/cloud/gcp/runtime/resource"
	"github.com/nitrictech/nitric/cloud/gcp/runtime/secret"
	sql_service "github.com/nitrictech/nitric/cloud/gcp/runtime/sql"
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

	sqlPlugin := sql_service.New()

	defaultGcpOpts := []server.ServerOption{
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
	defaultGcpOpts = append(defaultGcpOpts, opts...)

	return server.New(defaultGcpOpts...)
}
