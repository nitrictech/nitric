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

package job

import (
	batchpb "github.com/nitrictech/nitric/core/pkg/proto/batch/v1"
	kvstorepb "github.com/nitrictech/nitric/core/pkg/proto/kvstore/v1"
	queuespb "github.com/nitrictech/nitric/core/pkg/proto/queues/v1"
	secretspb "github.com/nitrictech/nitric/core/pkg/proto/secrets/v1"
	sqlpb "github.com/nitrictech/nitric/core/pkg/proto/sql/v1"
	storagepb "github.com/nitrictech/nitric/core/pkg/proto/storage/v1"
	topicspb "github.com/nitrictech/nitric/core/pkg/proto/topics/v1"
	websocketspb "github.com/nitrictech/nitric/core/pkg/proto/websockets/v1"
)

type nitricJobServerOption = func(*NitricJobServer)

func WithTopicPlugin(srv topicspb.TopicsServer) nitricJobServerOption {
	return func(o *NitricJobServer) {
		o.topicServer = srv
	}
}

func WithStoragePlugin(srv storagepb.StorageServer) nitricJobServerOption {
	return func(o *NitricJobServer) {
		o.storageServer = srv
	}
}

func WithQueuePlugin(srv queuespb.QueuesServer) nitricJobServerOption {
	return func(o *NitricJobServer) {
		o.queueServer = srv
	}
}

func WithSecretsPlugin(srv secretspb.SecretManagerServer) nitricJobServerOption {
	return func(o *NitricJobServer) {
		o.secretServer = srv
	}
}

func WithSqlPlugin(srv sqlpb.SqlServer) nitricJobServerOption {
	return func(o *NitricJobServer) {
		o.sqlServer = srv
	}
}

func WithKvStorePlugin(srv kvstorepb.KvStoreServer) nitricJobServerOption {
	return func(o *NitricJobServer) {
		o.kvStoreServer = srv
	}
}

func WithWebsocketPlugin(srv websocketspb.WebsocketServer) nitricJobServerOption {
	return func(o *NitricJobServer) {
		o.websocketServer = srv
	}
}

func WithBatchPlugin(srv batchpb.BatchServer) nitricJobServerOption {
	return func(o *NitricJobServer) {
		o.batchServer = srv
	}
}
