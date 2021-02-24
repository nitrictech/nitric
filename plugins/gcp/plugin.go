package main

import (
	identity_platform_service "github.com/nitric-dev/membrane/plugins/gcp/auth/identityplatform"
	firestore_service "github.com/nitric-dev/membrane/plugins/gcp/documents/firestore"
	pubsub_service "github.com/nitric-dev/membrane/plugins/gcp/eventing/pubsub"
	http_service "github.com/nitric-dev/membrane/plugins/gcp/gateway/http"
	pubsub_queue_service "github.com/nitric-dev/membrane/plugins/gcp/queue/pubsub"
	storage_service "github.com/nitric-dev/membrane/plugins/gcp/storage/storage"
	"github.com/nitric-dev/membrane/plugins/sdk"
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

// NewDocumentService - Returns Google Cloud Firestore based document service
func (p *GCPServiceFactory) NewDocumentService() (sdk.DocumentService, error) {
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