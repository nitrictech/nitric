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
	"github.com/nitric-dev/membrane/pkg/plugins/document"
	mongodb_service "github.com/nitric-dev/membrane/pkg/plugins/document/mongodb"
	"github.com/nitric-dev/membrane/pkg/plugins/events"
	"github.com/nitric-dev/membrane/pkg/plugins/gateway"
	http_service "github.com/nitric-dev/membrane/pkg/plugins/gateway/appservice"
	"github.com/nitric-dev/membrane/pkg/plugins/queue"
	"github.com/nitric-dev/membrane/pkg/plugins/storage"
	"github.com/nitric-dev/membrane/pkg/providers"
)

type AzureServiceFactory struct {
}

func New() providers.ServiceFactory {
	return &AzureServiceFactory{}
}

// NewDocumentService - Returns a MongoDB based document service
func (p *AzureServiceFactory) NewDocumentService() (document.DocumentService, error) {
	return mongodb_service.New()
}

// NewEventService - Returns Azure _ based events plugin
func (p *AzureServiceFactory) NewEventService() (events.EventService, error) {
	return &events.UnimplementedeventsPlugin{}, nil
}

// NewGatewayService - Returns Azure _ Gateway plugin
func (p *AzureServiceFactory) NewGatewayService() (gateway.GatewayService, error) {
	return http_service.New()
}

// NewQueueService - Returns Azure _ based queue plugin
func (p *AzureServiceFactory) NewQueueService() (queue.QueueService, error) {
	return &queue.UnimplementedQueuePlugin{}, nil
}

// NewStorageService - Returns Azure _ based storage plugin
func (p *AzureServiceFactory) NewStorageService() (storage.StorageService, error) {
	return &storage.UnimplementedStoragePlugin{}, nil
}
