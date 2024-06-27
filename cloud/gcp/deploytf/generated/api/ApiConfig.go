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
	// The name of the API Gateway.
	Name *string `field:"required" json:"name" yaml:"name"`
	// The OpenAPI spec as a JSON string.
	OpenapiSpec *string `field:"required" json:"openapiSpec" yaml:"openapiSpec"`
	// The ID of the stack.
	StackId *string `field:"required" json:"stackId" yaml:"stackId"`
	// The map of target service names The property type contains a map, they have special handling, please see {@link cdk.tf /module-map-inputs the docs}.
	TargetServices *map[string]*string `field:"required" json:"targetServices" yaml:"targetServices"`
}
