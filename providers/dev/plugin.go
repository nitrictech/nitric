// Copyright 2021 Nitric Pty Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	eventing_plugin "github.com/nitric-dev/membrane/plugins/eventing/dev"
	gateway_plugin "github.com/nitric-dev/membrane/plugins/gateway/dev"
	kv_plugin "github.com/nitric-dev/membrane/plugins/kv/boltdb"
	queue_plugin "github.com/nitric-dev/membrane/plugins/queue/dev"
	storage_plugin "github.com/nitric-dev/membrane/plugins/storage/boltdb"
	"github.com/nitric-dev/membrane/sdk"
)

type DevServiceFactory struct {
}

func New() sdk.ServiceFactory {
	return &DevServiceFactory{}
}

// NewEventingPlugin - Returns local dev eventing plugin
func (p *DevServiceFactory) NewEventService() (sdk.EventService, error) {
	return eventing_plugin.New()
}

// NewGatewayPlugin - Returns local dev Gateway plugin
func (p *DevServiceFactory) NewGatewayService() (sdk.GatewayService, error) {
	return gateway_plugin.New()
}

// NewKeyValuePlugin - Returns local dev key value plugin
func (p *DevServiceFactory) NewKeyValueService() (sdk.KeyValueService, error) {
	return kv_plugin.New()
}

// NewQueuePlugin - Returns local dev queue plugin
func (p *DevServiceFactory) NewQueueService() (sdk.QueueService, error) {
	return queue_plugin.New()
}

// NewStoragePlugin - Returns local dev storage plugin
func (p *DevServiceFactory) NewStorageService() (sdk.StorageService, error) {
	return storage_plugin.New()
}
