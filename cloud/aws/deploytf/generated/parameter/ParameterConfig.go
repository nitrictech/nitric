package parameter

import (
	"github.com/hashicorp/terraform-cdk-go/cdktf"
)

type ParameterConfig struct {
	// Experimental.
	DependsOn *[]cdktf.ITerraformDependable `field:"optional" json:"dependsOn" yaml:"dependsOn"`
	// Experimental.
	ForEach cdktf.ITerraformIterator `field:"optional" json:"forEach" yaml:"forEach"`
	// Experimental.
	Providers *[]interface{} `field:"optional" json:"providers" yaml:"providers"`
	// Experimental.
	SkipAssetCreationFromLocalModules *bool `field:"optional" json:"skipAssetCreationFromLocalModules" yaml:"skipAssetCreationFromLocalModules"`
	// The names of the roles that can access the parameter.
	AccessRoleNames *[]*string `field:"required" json:"accessRoleNames" yaml:"accessRoleNames"`
	// The name of the parameter.
	ParameterName *string `field:"required" json:"parameterName" yaml:"parameterName"`
	// The text value of the parameter.
	ParameterValue *string `field:"required" json:"parameterValue" yaml:"parameterValue"`
	// The tier of the SSM parameter (Standard or Advanced).
	// Docs: https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/ssm_parameter#tier
	ParameterTier *string `field:"optional" json:"parameterTier" yaml:"parameterTier"`
}

