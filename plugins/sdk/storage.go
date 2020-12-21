package sdk

import "fmt"

type StoragePlugin interface {
	Get(bucket string, key string) []byte, error
	Put(bucket string, key string, object []byte) error
}

type UnimplementedStoragePlugin struct {
	StoragePlugin
}

func(*UnimplementedStoragePlugin) Get(bucket string, key string) []byte, error {
	return nil, fmt.Errorf("UNIMPLEMENTED")
}

func(*UnimplementedStoragePlugin) Put(bucket string, key string, object []byte) error {
	return fmt.Errorf("UNIMPLEMENTED")
}
