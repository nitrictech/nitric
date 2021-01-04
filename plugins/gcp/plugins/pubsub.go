package main

import (
	pubsub_plugin "github.com/nitric-dev/membrane/plugins/gcp/eventing/pubsub"
	"github.com/nitric-dev/membrane/plugins/sdk"
)

func New() (sdk.GatewayPlugin, error) {
	return pubsub_plugin.New()
}
