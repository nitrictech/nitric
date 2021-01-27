package sdk

import "fmt"

type StoragePlugin interface {
	Get(bucket string, key string) ([]byte, error)
	Put(bucket string, key string, object []byte) error
	Delete(bucket string, key string) error
}

type UnimplementedStoragePlugin struct {
}

var _ StoragePlugin = (*UnimplementedStoragePlugin)(nil)

func (*UnimplementedStoragePlugin) Get(bucket string, key string) ([]byte, error) {
	return nil, fmt.Errorf("UNIMPLEMENTED")
}

func (*UnimplementedStoragePlugin) Put(bucket string, key string, object []byte) error {
	return fmt.Errorf("UNIMPLEMENTED")
}

func (*UnimplementedStoragePlugin) Delete(bucket string, key string) error {
	return fmt.Errorf("UNIMPLEMENTED")
}
