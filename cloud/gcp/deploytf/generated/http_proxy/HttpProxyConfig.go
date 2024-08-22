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
	// The email of the service account that will invoke the API.
	InvokerEmail *string `field:"required" json:"invokerEmail" yaml:"invokerEmail"`
	// The name of the HTTP proxy gateway.
	Name *string `field:"required" json:"name" yaml:"name"`
	// The ID of the stack.
	StackId *string `field:"required" json:"stackId" yaml:"stackId"`
	// The URL of the service being proxied.
	TargetServiceUrl *string `field:"required" json:"targetServiceUrl" yaml:"targetServiceUrl"`
}

