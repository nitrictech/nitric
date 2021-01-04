package main

import (
	s3_plugin "github.com/nitric-dev/membrane/plugins/aws/storage/s3"
	"github.com/nitric-dev/membrane/plugins/sdk"
)

func New() (sdk.StoragePlugin, error) {
	return s3_plugin.New()
}
