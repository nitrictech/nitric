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
	// The location where the bucket and its contents are stored.
	BucketLocation *string `field:"required" json:"bucketLocation" yaml:"bucketLocation"`
	// The name of the bucket.
	//
	// This must be globally unique.
	BucketName *string `field:"required" json:"bucketName" yaml:"bucketName"`
	// The notification target configurations.
	NotificationTargets interface{} `field:"required" json:"notificationTargets" yaml:"notificationTargets"`
	// The ID of the Google Cloud project where the bucket is created.
	ProjectId *string `field:"required" json:"projectId" yaml:"projectId"`
	// The ID of the Nitric stack.
	StackId *string `field:"required" json:"stackId" yaml:"stackId"`
	// The class of storage used to store the bucket's contents.
	//
	// This can be STANDARD, NEARLINE, COLDLINE, ARCHIVE, or MULTI_REGIONAL.
	// STANDARD.
	StorageClass *string `field:"optional" json:"storageClass" yaml:"storageClass"`
}

