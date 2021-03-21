package sdk

import "fmt"

// The base Documents Plugin interface
// Use this over proto definitions to remove dependency on protobuf in the plugin internally
// and open options to adding additional non-grpc interfaces
type KeyValueService interface {
	Put(string, string, map[string]interface{}) error
	Get(string, string) (map[string]interface{}, error)
	Delete(string, string) error
}

type UnimplementedKeyValuePlugin struct {
	DocumentService
}

func (p *UnimplementedKeyValuePlugin) Put(collection string, key string, value map[string]interface{}) error {
	return fmt.Errorf("UNIMPLEMENTED")
}

func (p *UnimplementedKeyValuePlugin) Get(collection string, key string) (map[string]interface{}, error) {
	return nil, fmt.Errorf("UNIMPLEMENTED")
}

func (p *UnimplementedKeyValuePlugin) Delete(collection string, key string) error {
	return fmt.Errorf("UNIMPLEMENTED")
}
