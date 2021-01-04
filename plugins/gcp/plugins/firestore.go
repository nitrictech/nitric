package main

import (
	firestore_plugin "github.com/nitric-dev/membrane/plugins/gcp/documents/firestore"
	"github.com/nitric-dev/membrane/plugins/sdk"
)

func New() (sdk.DocumentsPlugin, error) {
	return firestore_plugin.New()
}
