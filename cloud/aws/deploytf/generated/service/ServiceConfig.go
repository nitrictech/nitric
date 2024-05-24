package service

import (
	"github.com/hashicorp/terraform-cdk-go/cdktf"
)

type ServiceConfig struct {
	// Experimental.
	DependsOn *[]cdktf.ITerraformDependable `field:"optional" json:"dependsOn" yaml:"dependsOn"`
	// Experimental.
	ForEach cdktf.ITerraformIterator `field:"optional" json:"forEach" yaml:"forEach"`
	// Experimental.
	Providers *[]interface{} `field:"optional" json:"providers" yaml:"providers"`
	// Experimental.
	SkipAssetCreationFromLocalModules *bool `field:"optional" json:"skipAssetCreationFromLocalModules" yaml:"skipAssetCreationFromLocalModules"`
	// Environment variables to set on the lambda function The property type contains a map, they have special handling, please see {@link cdk.tf /module-map-inputs the docs}.
	Environment *map[string]*string `field:"required" json:"environment" yaml:"environment"`
	// The docker image to deploy.
	Image *string `field:"required" json:"image" yaml:"image"`
	// The name of the service.
	ServiceName *string `field:"required" json:"serviceName" yaml:"serviceName"`
	// The ID of the Nitric stack.
	StackId *string `field:"required" json:"stackId" yaml:"stackId"`
}

