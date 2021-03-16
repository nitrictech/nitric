package sdk

import "fmt"

// The base Documents Plugin interface
// Use this over proto definitions to remove dependency on protobuf in the plugin internally
// and open options to adding additional non-grpc interfaces
type DocumentService interface {
	Create(collection string, key string, document map[string]interface{}) error
	Get(collection string, key string) (map[string]interface{}, error)
	Update(collection string, key string, document map[string]interface{}) error
	Delete(collection string, key string) error
}

type UnimplementedDocumentsPlugin struct {
	DocumentService
}

func (p *UnimplementedDocumentsPlugin) Create(collection string, key string, document map[string]interface{}) error {
	return fmt.Errorf("UNIMPLEMENTED")
}

func (p *UnimplementedDocumentsPlugin) Get(collection string, key string) (map[string]interface{}, error) {
	return nil, fmt.Errorf("UNIMPLEMENTED")
}

func (p *UnimplementedDocumentsPlugin) Update(collection string, key string, document map[string]interface{}) error {
	return fmt.Errorf("UNIMPLEMENTED")
}

func (p *UnimplementedDocumentsPlugin) Delete(collection string, key string) error {
	return fmt.Errorf("UNIMPLEMENTED")
}
