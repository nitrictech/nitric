package schedule

import (
	"github.com/hashicorp/terraform-cdk-go/cdktf"
)

type ScheduleConfig struct {
	// Experimental.
	DependsOn *[]cdktf.ITerraformDependable `field:"optional" json:"dependsOn" yaml:"dependsOn"`
	// Experimental.
	ForEach cdktf.ITerraformIterator `field:"optional" json:"forEach" yaml:"forEach"`
	// Experimental.
	Providers *[]interface{} `field:"optional" json:"providers" yaml:"providers"`
	// Experimental.
	SkipAssetCreationFromLocalModules *bool `field:"optional" json:"skipAssetCreationFromLocalModules" yaml:"skipAssetCreationFromLocalModules"`
	// The schedule expression.
	ScheduleExpression *string `field:"required" json:"scheduleExpression" yaml:"scheduleExpression"`
	// The name of the schedule.
	ScheduleName *string `field:"required" json:"scheduleName" yaml:"scheduleName"`
	// The ID of the Nitric stack.
	StackId *string `field:"required" json:"stackId" yaml:"stackId"`
	// The ARN of the target lambda function.
	TargetLambdaArn *string `field:"required" json:"targetLambdaArn" yaml:"targetLambdaArn"`
}

