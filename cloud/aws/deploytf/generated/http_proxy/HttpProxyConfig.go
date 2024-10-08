package http_proxy

import (
	"github.com/hashicorp/terraform-cdk-go/cdktf"
)

type HttpProxyConfig struct {
	// Experimental.
	DependsOn *[]cdktf.ITerraformDependable `field:"optional" json:"dependsOn" yaml:"dependsOn"`
	// Experimental.
	ForEach cdktf.ITerraformIterator `field:"optional" json:"forEach" yaml:"forEach"`
	// Experimental.
	Providers *[]interface{} `field:"optional" json:"providers" yaml:"providers"`
	// Experimental.
	SkipAssetCreationFromLocalModules *bool `field:"optional" json:"skipAssetCreationFromLocalModules" yaml:"skipAssetCreationFromLocalModules"`
	// The name of the HTTP proxy gateway.
	Name *string `field:"required" json:"name" yaml:"name"`
	// The ID of the stack.
	StackId *string `field:"required" json:"stackId" yaml:"stackId"`
	// The name or arn of the lambda function being proxied.
	TargetLambdaFunction *string `field:"required" json:"targetLambdaFunction" yaml:"targetLambdaFunction"`
}

