package deploytf

import (
	"github.com/aws/aws-sdk-go/service/apigatewayv2"
	"github.com/hashicorp/terraform-cdk-go/cdktf"
)

type domainNameArgs struct {
	domainName string
	api        *apigatewayv2.Api
	stage      *apigatewayv2.Stage
}

func (a *NitricAwsTerraformProvider) newDomainName(stack cdktf.TerraformStack, name string, args domainNameArgs) error {
	return nil
}
