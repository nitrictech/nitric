package pubsub

import (
	"fmt"

	pubsubpb "github.com/nitrictech/nitric/proto/pubsub/v2"
	"github.com/nitrictech/nitric/server/runtime/plugin"
)

type Pubsub = pubsubpb.PubsubServer

// Available storage plugins for runtime
var pubsubPlugins = make(map[string]Pubsub)

func GetPlugin(namespace string) Pubsub {
	fmt.Println("available plugins", pubsubPlugins)
	return pubsubPlugins[namespace]
}

// Register a new instance of a storage plugin
func Register(namespace string, constructor plugin.Constructor[Pubsub]) error {
	pubsubPlugin, err := constructor()
	if err != nil {
		return err
	}

	pubsubPlugins[namespace] = pubsubPlugin
	return nil
}
