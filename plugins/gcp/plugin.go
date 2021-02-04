package main

import (
	identity_platform_plugin "github.com/nitric-dev/membrane/plugins/gcp/auth/identityplatform"
	firestore_plugin "github.com/nitric-dev/membrane/plugins/gcp/documents/firestore"
	pubsub_plugin "github.com/nitric-dev/membrane/plugins/gcp/eventing/pubsub"
	http_plugin "github.com/nitric-dev/membrane/plugins/gcp/gateway/http"
	pubsub_queue_plugin "github.com/nitric-dev/membrane/plugins/gcp/queue/pubsub"
	storage_plugin "github.com/nitric-dev/membrane/plugins/gcp/storage/storage"
	"github.com/nitric-dev/membrane/plugins/sdk"
)

type GCPPluginFactory struct {

}

func New() sdk.ProviderPluginFactory  {
	return &GCPPluginFactory{}
}

// NewAuthPlugin - Returns Google Cloud Identity Platform based auth plugin
func (p *GCPPluginFactory) NewAuthPlugin() (sdk.AuthPlugin, error) {
	return identity_platform_plugin.New()
}

// NewDocumentPlugin - Returns Google Cloud Firestore based document plugin
func (p *GCPPluginFactory) NewDocumentPlugin() (sdk.DocumentsPlugin, error) {
	return firestore_plugin.New()
}

// NewEventingPlugin - Returns Google Cloud Pubsub based eventing plugin
func (p *GCPPluginFactory) NewEventingPlugin() (sdk.EventingPlugin, error) {
	return pubsub_plugin.New()
}

// NewGatewayPlugin - Google Cloud Http Gateway plugin
func (p *GCPPluginFactory) NewGatewayPlugin() (sdk.GatewayPlugin, error) {
	return http_plugin.New()
}

// NewQueuePlugin - Returns Google Cloud Pubsub based queue plugin
func (p *GCPPluginFactory) NewQueuePlugin() (sdk.QueuePlugin, error) {
	return pubsub_queue_plugin.New()
}

// NewStoragePlugin - Returns Google Cloud Storage based storage plugin
func (p *GCPPluginFactory) NewStoragePlugin() (sdk.StoragePlugin, error) {
	return storage_plugin.New()
}