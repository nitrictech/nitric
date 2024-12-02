package topic

import (
	"github.com/hashicorp/terraform-cdk-go/cdktf"
)

type TopicConfig struct {
	// Experimental.
	DependsOn *[]cdktf.ITerraformDependable `field:"optional" json:"dependsOn" yaml:"dependsOn"`
	// Experimental.
	ForEach cdktf.ITerraformIterator `field:"optional" json:"forEach" yaml:"forEach"`
	// Experimental.
	Providers *[]interface{} `field:"optional" json:"providers" yaml:"providers"`
	// Experimental.
	SkipAssetCreationFromLocalModules *bool `field:"optional" json:"skipAssetCreationFromLocalModules" yaml:"skipAssetCreationFromLocalModules"`
	// The ID of the Nitric stack.
	StackId *string `field:"required" json:"stackId" yaml:"stackId"`
	// The services to create subscriptions for.
	SubscriberServices interface{} `field:"required" json:"subscriberServices" yaml:"subscriberServices"`
	// The name of the topic.
	TopicName *string `field:"required" json:"topicName" yaml:"topicName"`
	// The KMS key to use for encryption.
	KmsKey *string `field:"optional" json:"kmsKey" yaml:"kmsKey"`
}

