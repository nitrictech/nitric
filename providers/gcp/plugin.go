package main

import (
	identity_platform_service "github.com/nitric-dev/membrane/plugins/auth/identityplatform"
	pubsub_service "github.com/nitric-dev/membrane/plugins/eventing/pubsub"
	http_service "github.com/nitric-dev/membrane/plugins/gateway/cloudrun"
	firestore_service "github.com/nitric-dev/membrane/plugins/kv/firestore"
	pubsub_queue_service "github.com/nitric-dev/membrane/plugins/queue/pubsub"
	storage_service "github.com/nitric-dev/membrane/plugins/storage/storage"
	"github.com/nitric-dev/membrane/sdk"
)

type GCPServiceFactory struct {
}

func New() sdk.ServiceFactory {
	return &GCPServiceFactory{}
}

// NewAuthService - Returns Google Cloud Identity Platform based auth service
func (p *GCPServiceFactory) NewAuthService() (sdk.UserService, error) {
	return identity_platform_service.New()
}

// NewKeyValueService - Returns Google Cloud Firestore based kv service
func (p *GCPServiceFactory) NewKeyValueService() (sdk.KeyValueService, error) {
	return firestore_service.New()
}

// NewEventService - Returns Google Cloud Pubsub based eventing service
func (p *GCPServiceFactory) NewEventService() (sdk.EventService, error) {
	return pubsub_service.New()
}

// NewGatewayService - Google Cloud Http Gateway service
func (p *GCPServiceFactory) NewGatewayService() (sdk.GatewayService, error) {
	return http_service.New()
}

// NewQueueService - Returns Google Cloud Pubsub based queue service
func (p *GCPServiceFactory) NewQueueService() (sdk.QueueService, error) {
	return pubsub_queue_service.New()
}

// NewStorageService - Returns Google Cloud Storage based storage service
func (p *GCPServiceFactory) NewStorageService() (sdk.StorageService, error) {
	return storage_service.New()
}
