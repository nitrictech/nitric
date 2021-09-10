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
	boltdb_service "github.com/nitric-dev/membrane/pkg/plugins/document/boltdb"
	"github.com/nitric-dev/membrane/pkg/plugins/emails"
	email_service "github.com/nitric-dev/membrane/pkg/plugins/emails/dev"
	"github.com/nitric-dev/membrane/pkg/plugins/events"
	events_service "github.com/nitric-dev/membrane/pkg/plugins/events/dev"
	"github.com/nitric-dev/membrane/pkg/plugins/gateway"
	gateway_plugin "github.com/nitric-dev/membrane/pkg/plugins/gateway/dev"
	"github.com/nitric-dev/membrane/pkg/plugins/queue"
	queue_service "github.com/nitric-dev/membrane/pkg/plugins/queue/dev"
	"github.com/nitric-dev/membrane/pkg/plugins/secret"
	secret_service "github.com/nitric-dev/membrane/pkg/plugins/secret/dev"
	"github.com/nitric-dev/membrane/pkg/plugins/storage"
	boltdb_storage_service "github.com/nitric-dev/membrane/pkg/plugins/storage/boltdb"
	"github.com/nitric-dev/membrane/pkg/providers"
)

type DevServiceFactory struct {
}

func New() providers.ServiceFactory {
	return &DevServiceFactory{}
}

// NewDocumentService - Returns local dev document plugin
func (p *DevServiceFactory) NewDocumentService() (document.DocumentService, error) {
	return boltdb_service.New()
}

// NewEmailService - Returns local dev emails plugin
func (p *DevServiceFactory) NewEmailService() (emails.EmailService, error) {
	return email_service.New()
}

// NewEventService - Returns local dev events plugin
func (p *DevServiceFactory) NewEventService() (events.EventService, error) {
	return events_service.New()
}

// NewGatewayService - Returns local dev gateway plugin
func (p *DevServiceFactory) NewGatewayService() (gateway.GatewayService, error) {
	return gateway_plugin.New()
}

// NewQueueService - Returns local dev queue plugin
func (p *DevServiceFactory) NewQueueService() (queue.QueueService, error) {
	return queue_service.New()
}

// NewSecretService - Returns local dev secret plugin
func (p *DevServiceFactory) NewSecretService() (secret.SecretService, error) {
	return secret_service.New()
}

// NewStorageService - Returns local dev storage plugin
func (p *DevServiceFactory) NewStorageService() (storage.StorageService, error) {
	return boltdb_storage_service.New()
}
