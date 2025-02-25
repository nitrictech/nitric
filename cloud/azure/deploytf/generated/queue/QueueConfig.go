package queue

import (
	"github.com/hashicorp/terraform-cdk-go/cdktf"
)

type QueueConfig struct {
	// Experimental.
	DependsOn *[]cdktf.ITerraformDependable `field:"optional" json:"dependsOn" yaml:"dependsOn"`
	// Experimental.
	ForEach cdktf.ITerraformIterator `field:"optional" json:"forEach" yaml:"forEach"`
	// Experimental.
	Providers *[]interface{} `field:"optional" json:"providers" yaml:"providers"`
	// Experimental.
	SkipAssetCreationFromLocalModules *bool `field:"optional" json:"skipAssetCreationFromLocalModules" yaml:"skipAssetCreationFromLocalModules"`
	// the name of the queue.
	Name *string `field:"required" json:"name" yaml:"name"`
	// the name of the storage account.
	StorageAccountName *string `field:"required" json:"storageAccountName" yaml:"storageAccountName"`
}

