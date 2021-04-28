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
	http_service "github.com/nitric-dev/membrane/plugins/gateway/appservice"
	"github.com/nitric-dev/membrane/sdk"
)

type AzureServiceFactory struct {
}

func New() sdk.ServiceFactory {
	return &AzureServiceFactory{}
}

// NewAuthPlugin - Returns Azure _ based auth plugin
func (p *AzureServiceFactory) NewAuthService() (sdk.UserService, error) {
	return &sdk.UnimplementedAuthPlugin{}, nil
}

// NewDocumentPlugin - Returns Azure _ based document plugin
func (p *AzureServiceFactory) NewKeyValueService() (sdk.KeyValueService, error) {
	return &sdk.UnimplementedDocumentsPlugin{}, nil
}

// NewEventingPlugin - Returns Azure _ based eventing plugin
func (p *AzureServiceFactory) NewEventService() (sdk.EventService, error) {
	return &sdk.UnimplementedEventingPlugin{}, nil
}

// NewGatewayPlugin - Returns Azure _ Gateway plugin
func (p *AzureServiceFactory) NewGatewayService() (sdk.GatewayService, error) {
	return http_service.New()
}

// NewQueuePlugin - Returns Azure _ based queue plugin
func (p *AzureServiceFactory) NewQueueService() (sdk.QueueService, error) {
	return &sdk.UnimplementedQueuePlugin{}, nil
}

// NewStoragePlugin - Returns Azure _ based storage plugin
func (p *AzureServiceFactory) NewStorageService() (sdk.StorageService, error) {
	return &sdk.UnimplementedStoragePlugin{}, nil
}
