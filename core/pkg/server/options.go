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

package server

import (
	"github.com/nitrictech/nitric/core/pkg/gateway"
	kvstorepb "github.com/nitrictech/nitric/core/pkg/proto/kvstore/v1"
	queuespb "github.com/nitrictech/nitric/core/pkg/proto/queues/v1"
	resourcespb "github.com/nitrictech/nitric/core/pkg/proto/resources/v1"
	secretspb "github.com/nitrictech/nitric/core/pkg/proto/secrets/v1"
	sqlpb "github.com/nitrictech/nitric/core/pkg/proto/sql/v1"
	storagepb "github.com/nitrictech/nitric/core/pkg/proto/storage/v1"
	topicspb "github.com/nitrictech/nitric/core/pkg/proto/topics/v1"
	websocketspb "github.com/nitrictech/nitric/core/pkg/proto/websockets/v1"
	"github.com/nitrictech/nitric/core/pkg/workers/apis"
	"github.com/nitrictech/nitric/core/pkg/workers/http"
	"github.com/nitrictech/nitric/core/pkg/workers/schedules"
	"github.com/nitrictech/nitric/core/pkg/workers/storage"
	"github.com/nitrictech/nitric/core/pkg/workers/topics"
	"github.com/nitrictech/nitric/core/pkg/workers/websockets"
)

type ServerOption func(opts *NitricServer)

func WithResourcesPlugin(resources resourcespb.ResourcesServer) ServerOption {
	return func(opts *NitricServer) {
		opts.ResourcesPlugin = resources
	}
}

func WithGatewayPlugin(gw gateway.GatewayService) ServerOption {
	return func(opts *NitricServer) {
		opts.GatewayPlugin = gw
	}
}

func WithKeyValuePlugin(kv kvstorepb.KvStoreServer) ServerOption {
	return func(opts *NitricServer) {
		opts.KeyValuePlugin = kv
	}
}

func WithTopicsPlugin(tp topicspb.TopicsServer) ServerOption {
	return func(opts *NitricServer) {
		opts.TopicsPlugin = tp
	}
}

func WithStoragePlugin(sp storagepb.StorageServer) ServerOption {
	return func(opts *NitricServer) {
		opts.StoragePlugin = sp
	}
}

func WithSecretManagerPlugin(sm secretspb.SecretManagerServer) ServerOption {
	return func(opts *NitricServer) {
		opts.SecretManagerPlugin = sm
	}
}

func WithWebsocketPlugin(ws websocketspb.WebsocketServer) ServerOption {
	return func(opts *NitricServer) {
		opts.WebsocketPlugin = ws
	}
}

func WithQueuesPlugin(qs queuespb.QueuesServer) ServerOption {
	return func(opts *NitricServer) {
		opts.QueuesPlugin = qs
	}
}

func WithSqlPlugin(sql sqlpb.SqlServer) ServerOption {
	return func(opts *NitricServer) {
		opts.SqlPlugin = sql
	}
}

func WithApiPlugin(api apis.ApiRequestHandler) ServerOption {
	return func(opts *NitricServer) {
		opts.ApiPlugin = api
	}
}

func WithHttpPlugin(http http.HttpRequestHandler) ServerOption {
	return func(opts *NitricServer) {
		opts.HttpPlugin = http
	}
}

func WithSchedulesPlugin(schedules schedules.ScheduleRequestHandler) ServerOption {
	return func(opts *NitricServer) {
		opts.SchedulesPlugin = schedules
	}
}

func WithTopicsListenerPlugin(topics topics.SubscriptionRequestHandler) ServerOption {
	return func(opts *NitricServer) {
		opts.TopicsListenerPlugin = topics
	}
}

func WithStorageListenerPlugin(storage storage.BucketRequestHandler) ServerOption {
	return func(opts *NitricServer) {
		opts.StorageListenerPlugin = storage
	}
}

func WithWebsocketListenerPlugin(websockets websockets.WebsocketRequestHandler) ServerOption {
	return func(opts *NitricServer) {
		opts.WebsocketListenerPlugin = websockets
	}
}

func WithServiceAddress(address string) ServerOption {
	return func(opts *NitricServer) {
		opts.ServiceAddress = address
	}
}

// WithMinWorkers - Set the minimum number of workers that need to be available.
// this option is ignored if the MIN_WORKERS environment variable is set
func WithMinWorkers(minWorkers int) ServerOption {
	return func(opts *NitricServer) {
		opts.MinWorkers = minWorkers
	}
}

func WithChildCommand(command []string) ServerOption {
	return func(opts *NitricServer) {
		opts.ChildCommand = command
	}
}

func WithPreCommands(commands [][]string) ServerOption {
	return func(opts *NitricServer) {
		opts.PreCommands = commands
	}
}

func WithChildTimeoutSeconds(timeout int) ServerOption {
	return func(opts *NitricServer) {
		opts.ChildTimeoutSeconds = timeout
	}
}
