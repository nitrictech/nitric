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

package providers

import (
	"github.com/nitric-dev/membrane/pkg/plugins/document"
	"github.com/nitric-dev/membrane/pkg/plugins/emails"
	"github.com/nitric-dev/membrane/pkg/plugins/events"
	"github.com/nitric-dev/membrane/pkg/plugins/gateway"
	"github.com/nitric-dev/membrane/pkg/plugins/queue"
	"github.com/nitric-dev/membrane/pkg/plugins/storage"
)

// ServiceFactory - interface for Service Factory Plugins, which instantiate provider specific service implementations.
type ServiceFactory interface {
	NewDocumentService() (document.DocumentService, error)
	NewEventService() (events.EventService, error)
	NewGatewayService() (gateway.GatewayService, error)
	NewQueueService() (queue.QueueService, error)
	NewStorageService() (storage.StorageService, error)
	NewEmailService() (emails.EmailService, error)
}

// UnimplementedServiceFactory - provides stub methods for a ServiceFactory which return Unimplemented Methods.
//
// Returning nil from a New service method is a valid response. Without an accompanying error, this will be
// interpreted as the method being explicitly unimplemented.
//
// Plugin Factories with unimplemented New methods are only supported when the TOLERATE_MISSING_SERVICE option is
// set to true when executing the pluggable membrane.
type UnimplementedServiceFactory struct {
}

// Ensure UnimplementedServiceFactory implement all methods of the interface
var _ ServiceFactory = (*UnimplementedServiceFactory)(nil)

// NewDocumentService - Unimplemented
func (p *UnimplementedServiceFactory) NewDocumentService() (document.DocumentService, error) {
	return nil, nil
}

// NewEventService - Unimplemented
func (p *UnimplementedServiceFactory) NewEventService() (events.EventService, error) {
	return nil, nil
}

// NewGatewayService - Unimplemented
func (p *UnimplementedServiceFactory) NewGatewayService() (gateway.GatewayService, error) {
	return nil, nil
}

// NewQueueService - Unimplemented
func (p *UnimplementedServiceFactory) NewQueueService() (queue.QueueService, error) {
	return nil, nil
}

// NewStorageService - Unimplemented
func (p *UnimplementedServiceFactory) NewStorageService() (storage.StorageService, error) {
	return nil, nil
}

// NewEmailService - Unimplemented
func (p *UnimplementedServiceFactory) NewEmailService() (emails.EmailService, error) {
	return nil, nil
}
