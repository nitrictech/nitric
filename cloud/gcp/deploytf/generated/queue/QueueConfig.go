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
	// The name of the queue.
	QueueName *string `field:"required" json:"queueName" yaml:"queueName"`
	// The ID of the Nitric stack.
	StackId *string `field:"required" json:"stackId" yaml:"stackId"`
	// The KMS key to use for encryption.
	KmsKey *string `field:"optional" json:"kmsKey" yaml:"kmsKey"`
}

