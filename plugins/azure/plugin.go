package main

import (
	http_service "github.com/nitric-dev/membrane/plugins/aws/gateway/http"
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
func (p *AzureServiceFactory) NewDocumentService() (sdk.DocumentService, error) {
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
