package bucket

import (
	"github.com/hashicorp/terraform-cdk-go/cdktf"
)

type BucketConfig struct {
	// Experimental.
	DependsOn *[]cdktf.ITerraformDependable `field:"optional" json:"dependsOn" yaml:"dependsOn"`
	// Experimental.
	ForEach cdktf.ITerraformIterator `field:"optional" json:"forEach" yaml:"forEach"`
	// Experimental.
	Providers *[]interface{} `field:"optional" json:"providers" yaml:"providers"`
	// Experimental.
	SkipAssetCreationFromLocalModules *bool `field:"optional" json:"skipAssetCreationFromLocalModules" yaml:"skipAssetCreationFromLocalModules"`
	// The list of listeners to notify.
	Listeners interface{} `field:"required" json:"listeners" yaml:"listeners"`
	// The name of the bucket.
	Name *string `field:"required" json:"name" yaml:"name"`
	// The name of the stack.
	StackName *string `field:"required" json:"stackName" yaml:"stackName"`
	// The id of the storage account.
	StorageAccountId *string `field:"required" json:"storageAccountId" yaml:"storageAccountId"`
}

