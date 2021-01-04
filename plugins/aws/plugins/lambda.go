package main

import (
	lambda_plugin "github.com/nitric-dev/membrane/plugins/aws/gateway/lambda"
	"github.com/nitric-dev/membrane/plugins/sdk"
)

func New() (sdk.GatewayPlugin, error) {
	return lambda_plugin.New()
}
