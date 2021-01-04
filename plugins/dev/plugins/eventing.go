package main

import (
	eventing_plugin "github.com/nitric-dev/membrane/plugins/dev/eventing"
	"github.com/nitric-dev/membrane/plugins/sdk"
)

func New() (sdk.EventingPlugin, error) {
	return eventing_plugin.New()
}
