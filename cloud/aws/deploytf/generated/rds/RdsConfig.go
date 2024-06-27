package rds

import (
	"github.com/hashicorp/terraform-cdk-go/cdktf"
)

type RdsConfig struct {
	// Experimental.
	DependsOn *[]cdktf.ITerraformDependable `field:"optional" json:"dependsOn" yaml:"dependsOn"`
	// Experimental.
	ForEach cdktf.ITerraformIterator `field:"optional" json:"forEach" yaml:"forEach"`
	// Experimental.
	Providers *[]interface{} `field:"optional" json:"providers" yaml:"providers"`
	// Experimental.
	SkipAssetCreationFromLocalModules *bool `field:"optional" json:"skipAssetCreationFromLocalModules" yaml:"skipAssetCreationFromLocalModules"`
	// the maximum capacity of the RDS cluster.
	MaxCapacity *float64 `field:"required" json:"maxCapacity" yaml:"maxCapacity"`
	// the minimum capacity of the RDS cluster.
	MinCapacity *float64 `field:"required" json:"minCapacity" yaml:"minCapacity"`
	// private subnets to assign to the RDS cluster.
	PrivateSubnetIds *[]*string `field:"required" json:"privateSubnetIds" yaml:"privateSubnetIds"`
	// The nitric stack ID.
	StackId *string `field:"required" json:"stackId" yaml:"stackId"`
	// the VPC to assign to the RDS cluster.
	VpcId *string `field:"required" json:"vpcId" yaml:"vpcId"`
}
