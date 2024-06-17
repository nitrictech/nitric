package api

import (
	"github.com/hashicorp/terraform-cdk-go/cdktf"
)

type ApiConfig struct {
	// Experimental.
	DependsOn *[]cdktf.ITerraformDependable `field:"optional" json:"dependsOn" yaml:"dependsOn"`
	// Experimental.
	ForEach cdktf.ITerraformIterator `field:"optional" json:"forEach" yaml:"forEach"`
	// Experimental.
	Providers *[]interface{} `field:"optional" json:"providers" yaml:"providers"`
	// Experimental.
	SkipAssetCreationFromLocalModules *bool `field:"optional" json:"skipAssetCreationFromLocalModules" yaml:"skipAssetCreationFromLocalModules"`
	// The domains to associate with the API Gateway.
	Domains *[]*string `field:"required" json:"domains" yaml:"domains"`
	// The name of the API Gateway.
	Name *string `field:"required" json:"name" yaml:"name"`
	// Open API spec.
	Spec *string `field:"required" json:"spec" yaml:"spec"`
	// The ID of the stack.
	StackId *string `field:"required" json:"stackId" yaml:"stackId"`
	// The names of the target lambda functions The property type contains a map, they have special handling, please see {@link cdk.tf /module-map-inputs the docs}.
	TargetLambdaFunctions *map[string]*string `field:"required" json:"targetLambdaFunctions" yaml:"targetLambdaFunctions"`
}
