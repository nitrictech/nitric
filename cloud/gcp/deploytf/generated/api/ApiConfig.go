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
	// The email of the service account that will invoke the API.
	InvokerEmail *string `field:"required" json:"invokerEmail" yaml:"invokerEmail"`
	// The name of the API Gateway.
	Name *string `field:"required" json:"name" yaml:"name"`
	// The OpenAPI spec as a JSON string.
	OpenapiSpec *string `field:"required" json:"openapiSpec" yaml:"openapiSpec"`
	// The GCP project ID.
	ProjectId *string `field:"required" json:"projectId" yaml:"projectId"`
	// The ID of the stack.
	StackId *string `field:"required" json:"stackId" yaml:"stackId"`
	// The list of target services.
	TargetServices interface{} `field:"required" json:"targetServices" yaml:"targetServices"`
}

