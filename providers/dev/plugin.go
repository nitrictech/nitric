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
	auth_service "github.com/nitric-dev/membrane/plugins/dev/auth"
	documents_plugin "github.com/nitric-dev/membrane/plugins/dev/documents"
	eventing_plugin "github.com/nitric-dev/membrane/plugins/dev/eventing"
	gateway_plugin "github.com/nitric-dev/membrane/plugins/dev/gateway"
	queue_plugin "github.com/nitric-dev/membrane/plugins/dev/queue"
	storage_plugin "github.com/nitric-dev/membrane/plugins/dev/storage"
	"github.com/nitric-dev/membrane/sdk"
)

type DevServiceFactory struct {
}

func New() sdk.ServiceFactory {
	return &DevServiceFactory{}
}

// NewAuthPlugin - Returns AWS Cognito based auth plugin
func (p *DevServiceFactory) NewAuthService() (sdk.UserService, error) {
	return auth_service.New()
}

// NewDocumentPlugin - Returns AWS DynamoDB based document plugin
func (p *DevServiceFactory) NewDocumentService() (sdk.DocumentService, error) {
	return documents_plugin.New()
}

// NewEventingPlugin - Returns AWS SNS based eventing plugin
func (p *DevServiceFactory) NewEventService() (sdk.EventService, error) {
	return eventing_plugin.New()
}

// NewGatewayPlugin - Returns AWS Lambda Gateway plugin
func (p *DevServiceFactory) NewGatewayService() (sdk.GatewayService, error) {
	return gateway_plugin.New()
}

// NewQueuePlugin - Returns AWS SQS based queue plugin
func (p *DevServiceFactory) NewQueueService() (sdk.QueueService, error) {
	return queue_plugin.New()
}

// NewStoragePlugin - Returns AWS S3 based storage plugin
func (p *DevServiceFactory) NewStorageService() (sdk.StorageService, error) {
	return storage_plugin.New()
}
