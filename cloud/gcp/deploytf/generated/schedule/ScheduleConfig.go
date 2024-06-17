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
	// The timezone for the schedule.
	ScheduleTimezone *string `field:"required" json:"scheduleTimezone" yaml:"scheduleTimezone"`
	// The token to authenticate with the target service.
	ServiceToken *string `field:"required" json:"serviceToken" yaml:"serviceToken"`
	// The URL of the target service.
	TargetServiceUrl *string `field:"required" json:"targetServiceUrl" yaml:"targetServiceUrl"`
}

