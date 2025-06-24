package storage

import (
	"fmt"

	storagepb "github.com/nitrictech/nitric/proto/storage/v2"
	"github.com/nitrictech/nitric/server/runtime/plugin"
)

// Define the interface for a storage plugin here
type Storage = storagepb.StorageServer

// Available storage plugins for runtime
var storagePlugins = make(map[string]Storage)

func GetPlugin(namespace string) Storage {
	fmt.Println("available plugins", storagePlugins)
	return storagePlugins[namespace]
}

// Register a new instance of a storage plugin
func Register(namespace string, constructor plugin.Constructor[Storage]) error {
	storagePlugin, err := constructor()
	if err != nil {
		return err
	}

	storagePlugins[namespace] = storagePlugin
	return nil
}
