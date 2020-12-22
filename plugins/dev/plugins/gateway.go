package main

import (
	gateway_plugin "github.com/nitric-dev/membrane/plugins/dev/gateway"
	"github.com/nitric-dev/membrane/plugins/sdk"
)

func New() (sdk.GatewayPlugin, error) {
	return gateway_plugin.New()
}
