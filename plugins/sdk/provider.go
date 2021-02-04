package sdk

import "fmt"

// ProviderPluginFactory - interface for Provider Plugin Factories, which
// create provider specific plugin implementations of services.
type ProviderPluginFactory interface {
	NewAuthPlugin() (AuthPlugin, error)
	NewDocumentPlugin() (DocumentsPlugin, error)
	NewEventingPlugin() (EventingPlugin, error)
	NewGatewayPlugin() (GatewayPlugin, error)
	NewQueuePlugin() (QueuePlugin, error)
	NewStoragePlugin() (StoragePlugin, error)
}

// UnimplementedProviderPluginFactory - provides stub methods for a ProviderPluginFactory which return Unimplemented Methods.
type UnimplementedProviderPluginFactory struct {

}

// Ensure UnimplementedProviderPluginFactory implement all methods of the interface
var _ ProviderPluginFactory = (*UnimplementedProviderPluginFactory)(nil)

// NewAuthPlugin - Unimplemented
func (p *UnimplementedProviderPluginFactory) NewAuthPlugin() (AuthPlugin, error) {
	return nil, fmt.Errorf("UNIMPLEMENTED")
}

// NewDocumentPlugin - Unimplemented
func (p *UnimplementedProviderPluginFactory) NewDocumentPlugin() (DocumentsPlugin, error) {
	return nil, fmt.Errorf("UNIMPLEMENTED")
}

// NewEventingPlugin - Unimplemented
func (p *UnimplementedProviderPluginFactory) NewEventingPlugin() (EventingPlugin, error) {
	return nil, fmt.Errorf("UNIMPLEMENTED")
}

// NewGatewayPlugin - Unimplemented
func (p *UnimplementedProviderPluginFactory) NewGatewayPlugin() (GatewayPlugin, error) {
	return nil, fmt.Errorf("UNIMPLEMENTED")
}

// NewQueuePlugin - Unimplemented
func (p *UnimplementedProviderPluginFactory) NewQueuePlugin() (QueuePlugin, error) {
	return nil, fmt.Errorf("UNIMPLEMENTED")
}

// NewStoragePlugin - Unimplemented
func (p *UnimplementedProviderPluginFactory) NewStoragePlugin() (StoragePlugin, error) {
	return nil, fmt.Errorf("UNIMPLEMENTED")
}
