package sql

import (
	"github.com/hashicorp/terraform-cdk-go/cdktf"
)

type SqlConfig struct {
	// Experimental.
	DependsOn *[]cdktf.ITerraformDependable `field:"optional" json:"dependsOn" yaml:"dependsOn"`
	// Experimental.
	ForEach cdktf.ITerraformIterator `field:"optional" json:"forEach" yaml:"forEach"`
	// Experimental.
	Providers *[]interface{} `field:"optional" json:"providers" yaml:"providers"`
	// Experimental.
	SkipAssetCreationFromLocalModules *bool `field:"optional" json:"skipAssetCreationFromLocalModules" yaml:"skipAssetCreationFromLocalModules"`
	// The region of the codebuild project.
	CodebuildRegion *string `field:"required" json:"codebuildRegion" yaml:"codebuildRegion"`
	// The arn of the codebuild role.
	CodebuildRoleArn *string `field:"required" json:"codebuildRoleArn" yaml:"codebuildRoleArn"`
	// The name of the create database codebuild project.
	CreateDatabaseProjectName *string `field:"required" json:"createDatabaseProjectName" yaml:"createDatabaseProjectName"`
	// The name of the database to create.
	DbName *string `field:"required" json:"dbName" yaml:"dbName"`
	// The URI of the docker image to use for the codebuild project.
	ImageUri *string `field:"required" json:"imageUri" yaml:"imageUri"`
	// The command to run to migrate the database.
	MigrateCommand *string `field:"required" json:"migrateCommand" yaml:"migrateCommand"`
	// The endpoint of the RDS cluster to connect to.
	RdsClusterEndpoint *string `field:"required" json:"rdsClusterEndpoint" yaml:"rdsClusterEndpoint"`
	// The password to connect to the RDS cluster.
	RdsClusterPassword *string `field:"required" json:"rdsClusterPassword" yaml:"rdsClusterPassword"`
	// The username to connect to the RDS cluster.
	RdsClusterUsername *string `field:"required" json:"rdsClusterUsername" yaml:"rdsClusterUsername"`
	// The security group ids to use for the codebuild project.
	SecurityGroupIds *[]*string `field:"required" json:"securityGroupIds" yaml:"securityGroupIds"`
	// The subnet ids to use for the codebuild project.
	SubnetIds *[]*string `field:"required" json:"subnetIds" yaml:"subnetIds"`
	// The vpc id to use for the codebuild project.
	VpcId *string `field:"required" json:"vpcId" yaml:"vpcId"`
	// The working directory for the codebuild project.
	WorkDir *string `field:"required" json:"workDir" yaml:"workDir"`
}
