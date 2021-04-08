package sdk

// ServiceFactory - interface for Service Factory Plugins, which instantiate provider specific service implementations.
type ServiceFactory interface {
	NewAuthService() (UserService, error)
	NewKeyValueService() (KeyValueService, error)
	NewEventService() (EventService, error)
	NewGatewayService() (GatewayService, error)
	NewQueueService() (QueueService, error)
	NewStorageService() (StorageService, error)
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

// NewAuthService - Unimplemented
func (p *UnimplementedServiceFactory) NewAuthService() (UserService, error) {
	return nil, nil
}

// NewKeyValueService - Unimplemented
func (p *UnimplementedServiceFactory) NewKeyValueService() (KeyValueService, error) {
	return nil, nil
}

// NewEventService - Unimplemented
func (p *UnimplementedServiceFactory) NewEventService() (EventService, error) {
	return nil, nil
}

// NewGatewayService - Unimplemented
func (p *UnimplementedServiceFactory) NewGatewayService() (GatewayService, error) {
	return nil, nil
}

// NewQueueService - Unimplemented
func (p *UnimplementedServiceFactory) NewQueueService() (QueueService, error) {
	return nil, nil
}

// NewStorageService - Unimplemented
func (p *UnimplementedServiceFactory) NewStorageService() (StorageService, error) {
	return nil, nil
}
