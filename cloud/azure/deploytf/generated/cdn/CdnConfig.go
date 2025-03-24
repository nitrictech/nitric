package cdn

import (
	"github.com/hashicorp/terraform-cdk-go/cdktf"
)

type CdnConfig struct {
	// Experimental.
	DependsOn *[]cdktf.ITerraformDependable `field:"optional" json:"dependsOn" yaml:"dependsOn"`
	// Experimental.
	ForEach cdktf.ITerraformIterator `field:"optional" json:"forEach" yaml:"forEach"`
	// Experimental.
	Providers *[]interface{} `field:"optional" json:"providers" yaml:"providers"`
	// Experimental.
	SkipAssetCreationFromLocalModules *bool `field:"optional" json:"skipAssetCreationFromLocalModules" yaml:"skipAssetCreationFromLocalModules"`
	// The primary host for the CDN.
	PrimaryWebHost *string `field:"required" json:"primaryWebHost" yaml:"primaryWebHost"`
	// The name of the resource group to use for the cdn.
	ResourceGroupName *string `field:"required" json:"resourceGroupName" yaml:"resourceGroupName"`
	// The name of the stack.
	StackName *string `field:"required" json:"stackName" yaml:"stackName"`
	// Enable API rewrites.
	EnableApiRewrites *bool `field:"optional" json:"enableApiRewrites" yaml:"enableApiRewrites"`
	// Map of uploaded files with their MD5 hashes The property type contains a map, they have special handling, please see {@link cdk.tf /module-map-inputs the docs}.
	UploadedFiles *map[string]*string `field:"optional" json:"uploadedFiles" yaml:"uploadedFiles"`
}

