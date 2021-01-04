package main

import (
	storage_plugin "github.com/nitric-dev/membrane/plugins/dev/storage"
	"github.com/nitric-dev/membrane/plugins/sdk"
)

func New() (sdk.StoragePlugin, error) {
	return storage_plugin.New()
}
