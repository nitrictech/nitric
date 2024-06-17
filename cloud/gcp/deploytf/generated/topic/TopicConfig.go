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
	// A list of lambda ARNs to subscribe to the topic The property type contains a map, they have special handling, please see {@link cdk.tf /module-map-inputs the docs}.
	LambdaSubscribers *map[string]*string `field:"required" json:"lambdaSubscribers" yaml:"lambdaSubscribers"`
	// The ID of the Nitric stack.
	StackId *string `field:"required" json:"stackId" yaml:"stackId"`
	// The name of the bucket.
	//
	// This must be globally unique.
	TopicName *string `field:"required" json:"topicName" yaml:"topicName"`
}

