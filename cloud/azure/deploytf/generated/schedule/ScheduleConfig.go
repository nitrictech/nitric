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
	// The container app environment id.
	ContainerAppEnvironmentId *string `field:"required" json:"containerAppEnvironmentId" yaml:"containerAppEnvironmentId"`
	// The cron expression for the schedule.
	CronExpression *string `field:"required" json:"cronExpression" yaml:"cronExpression"`
	// The name of the schedule.
	Name *string `field:"required" json:"name" yaml:"name"`
	// The target app id for the schedule.
	TargetAppId *string `field:"required" json:"targetAppId" yaml:"targetAppId"`
	// The target event token for the schedule.
	TargetEventToken *string `field:"required" json:"targetEventToken" yaml:"targetEventToken"`
}

