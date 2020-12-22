package main

import (
	documents_plugin "github.com/nitric-dev/membrane/plugins/dev/documents"
	"github.com/nitric-dev/membrane/plugins/sdk"
)

func New() (sdk.DocumentsPlugin, error) {
	return documents_plugin.New()
}
