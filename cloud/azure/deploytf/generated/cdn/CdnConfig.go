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
	// The name of the resource group to use for the cdn.
	ResourceGroupName *string `field:"required" json:"resourceGroupName" yaml:"resourceGroupName"`
	// The name of the stack.
	StackName *string `field:"required" json:"stackName" yaml:"stackName"`
	// The primary web host of the storage account to use for the cdn.
	StorageAccountPrimaryWebHost *string `field:"required" json:"storageAccountPrimaryWebHost" yaml:"storageAccountPrimaryWebHost"`
	// Map of APIs and their gateway information.
	Apis interface{} `field:"optional" json:"apis" yaml:"apis"`
	// Map of content paths to purge from the CDN The property type contains a map, they have special handling, please see {@link cdk.tf /module-map-inputs the docs}.
	CdnPurgePaths *map[string]*string `field:"optional" json:"cdnPurgePaths" yaml:"cdnPurgePaths"`
}

