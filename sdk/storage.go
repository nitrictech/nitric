package sdk

import "fmt"

type StorageService interface {
	Read(bucket string, key string) ([]byte, error)
	Write(bucket string, key string, object []byte) error
	Delete(bucket string, key string) error
}

type UnimplementedStoragePlugin struct{}

var _ StorageService = (*UnimplementedStoragePlugin)(nil)

func (*UnimplementedStoragePlugin) Read(bucket string, key string) ([]byte, error) {
	return nil, fmt.Errorf("UNIMPLEMENTED")
}

func (*UnimplementedStoragePlugin) Write(bucket string, key string, object []byte) error {
	return fmt.Errorf("UNIMPLEMENTED")
}

func (*UnimplementedStoragePlugin) Delete(bucket string, key string) error {
	return fmt.Errorf("UNIMPLEMENTED")
}
