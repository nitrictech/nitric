package main

import (
	dynamodb_plugin "github.com/nitric-dev/membrane/plugins/aws/documents/dynamodb"
	"github.com/nitric-dev/membrane/plugins/sdk"
)

func New() (sdk.DocumentsPlugin, error) {
	return dynamodb_plugin.New()
}
