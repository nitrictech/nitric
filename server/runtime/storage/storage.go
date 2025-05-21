package storage

import (
	"fmt"

	storagepb "github.com/nitrictech/nitric/core/pkg/proto/storage/v1"
	"github.com/nitrictech/nitric/server/runtime"
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
func Register(namespace string, constructor runtime.PluginConstructor[Storage]) {
	storagePlugins[namespace] = constructor()
}
