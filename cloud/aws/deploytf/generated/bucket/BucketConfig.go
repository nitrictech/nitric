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
	// The name of the bucket.
	//
	// This must be globally unique.
	BucketName *string `field:"required" json:"bucketName" yaml:"bucketName"`
	// The notification target configurations.
	NotificationTargets interface{} `field:"required" json:"notificationTargets" yaml:"notificationTargets"`
	// The ID of the Nitric stack.
	StackId *string `field:"required" json:"stackId" yaml:"stackId"`
}
