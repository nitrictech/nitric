package sdk

import "fmt"

// The base Documents Plugin interface
// Use this over proto definitions to remove dependency on protobuf in the plugin internally
// and open options to adding additional non-grpc interfaces
type DocumentService interface {
	CreateDocument(collection string, key string, document map[string]interface{}) error
	GetDocument(collection string, key string) (map[string]interface{}, error)
	UpdateDocument(collection string, key string, document map[string]interface{}) error
	DeleteDocument(collection string, key string) error
}

type UnimplementedDocumentsPlugin struct {
	DocumentService
}

func (p *UnimplementedDocumentsPlugin) CreateDocument(collection string, key string, document map[string]interface{}) error {
	return fmt.Errorf("UNIMPLEMENTED")
}

func (p *UnimplementedDocumentsPlugin) GetDocument(collection string, key string) (map[string]interface{}, error) {
	return nil, fmt.Errorf("UNIMPLEMENTED")
}

func (p *UnimplementedDocumentsPlugin) UpdateDocument(collection string, key string, document map[string]interface{}) error {
	return fmt.Errorf("UNIMPLEMENTED")
}

func (p *UnimplementedDocumentsPlugin) DeleteDocument(collection string, key string) error {
	return fmt.Errorf("UNIMPLEMENTED")
}
