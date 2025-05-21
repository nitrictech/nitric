package pubsub

import (
	"fmt"

	pubsubpb "github.com/nitrictech/nitric/core/pkg/proto/topics/v1"
	"github.com/nitrictech/nitric/server/runtime"
)

type Pubsub = pubsubpb.TopicsServer

// Available storage plugins for runtime
var pubsubPlugins = make(map[string]Pubsub)

func GetPlugin(namespace string) Pubsub {
	fmt.Println("available plugins", pubsubPlugins)
	return pubsubPlugins[namespace]
}

// Register a new instance of a storage plugin
func Register(namespace string, constructor runtime.PluginConstructor[Pubsub]) {
	pubsubPlugins[namespace] = constructor()
}
