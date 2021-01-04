package main

import (
	sns_plugin "github.com/nitric-dev/membrane/plugins/aws/eventing/sns"
	"github.com/nitric-dev/membrane/plugins/sdk"
)

func New() (sdk.EventingPlugin, error) {
	return sns_plugin.New()
}
