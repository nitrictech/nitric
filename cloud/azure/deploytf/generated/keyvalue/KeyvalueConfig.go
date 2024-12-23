package keyvalue

import (
	"github.com/hashicorp/terraform-cdk-go/cdktf"
)

type KeyvalueConfig struct {
	// Experimental.
	DependsOn *[]cdktf.ITerraformDependable `field:"optional" json:"dependsOn" yaml:"dependsOn"`
	// Experimental.
	ForEach cdktf.ITerraformIterator `field:"optional" json:"forEach" yaml:"forEach"`
	// Experimental.
	Providers *[]interface{} `field:"optional" json:"providers" yaml:"providers"`
	// Experimental.
	SkipAssetCreationFromLocalModules *bool `field:"optional" json:"skipAssetCreationFromLocalModules" yaml:"skipAssetCreationFromLocalModules"`
	// The name of the kv store.
	Name *string `field:"required" json:"name" yaml:"name"`
	// The name of the storage account.
	StorageAccountName *string `field:"required" json:"storageAccountName" yaml:"storageAccountName"`
}

