package main

import (
	storage_plugin "github.com/nitric-dev/membrane/plugins/gcp/storage/storage"
	"github.com/nitric-dev/membrane/plugins/sdk"
)

func New() (sdk.GatewayPlugin, error) {
	return storage_plugin.New()
}
