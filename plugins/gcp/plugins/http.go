package main

import (
	http_plugin "github.com/nitric-dev/membrane/plugins/gcp/gateway/http"
	"github.com/nitric-dev/membrane/plugins/sdk"
)

func New() (sdk.GatewayPlugin, error) {
	return http_plugin.New()
}
