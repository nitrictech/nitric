package stack

import (
	"github.com/hashicorp/terraform-cdk-go/cdktf"
)

type StackConfig struct {
	// Experimental.
	DependsOn *[]cdktf.ITerraformDependable `field:"optional" json:"dependsOn" yaml:"dependsOn"`
	// Experimental.
	ForEach cdktf.ITerraformIterator `field:"optional" json:"forEach" yaml:"forEach"`
	// Experimental.
	Providers *[]interface{} `field:"optional" json:"providers" yaml:"providers"`
	// Experimental.
	SkipAssetCreationFromLocalModules *bool `field:"optional" json:"skipAssetCreationFromLocalModules" yaml:"skipAssetCreationFromLocalModules"`
	// Enable customer managed encryption keys.
	CmekEnabled *bool `field:"required" json:"cmekEnabled" yaml:"cmekEnabled"`
	// The location to deploy the stack.
	Location *string `field:"required" json:"location" yaml:"location"`
	// The name of the nitric stack.
	StackName *string `field:"required" json:"stackName" yaml:"stackName"`
}

