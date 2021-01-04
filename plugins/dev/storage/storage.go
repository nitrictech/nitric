package storage_plugin

import (
	"fmt"

	"github.com/nitric-dev/membrane/plugins/sdk"
)

type LocalStoragePlugin struct {
	sdk.UnimplementedStoragePlugin
}

func (s *LocalStoragePlugin) Put(bucket string, key string, payload []byte) error {
	return fmt.Errorf("UNIMPLEMENTED")
}

// Retrieve an item from a bucket
func (s *LocalStoragePlugin) Get(bucket string, key string) ([]byte, error) {
	return nil, fmt.Errorf("UNIMPLEMENTED")
}

// Create new DynamoDB documents server
// XXX: No External Args for function atm (currently the plugin loader does not pass any argument information)
func New() (sdk.StoragePlugin, error) {
	return &LocalStoragePlugin{}, nil
}
