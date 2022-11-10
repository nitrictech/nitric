package apigatewayv2iface

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/apigatewayv2"
)

type ApiGatewayV2API interface {
	GetApi(ctx context.Context, params *apigatewayv2.GetApiInput, optFns ...func(*apigatewayv2.Options)) (*apigatewayv2.GetApiOutput, error)
}
