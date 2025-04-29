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
	// The database master password.
	DatabaseMasterPassword *string `field:"required" json:"databaseMasterPassword" yaml:"databaseMasterPassword"`
	// The database server fully qualified domain name.
	DatabaseServerFqdn *string `field:"required" json:"databaseServerFqdn" yaml:"databaseServerFqdn"`
	// The image registry password.
	ImageRegistryPassword *string `field:"required" json:"imageRegistryPassword" yaml:"imageRegistryPassword"`
	// The image registry server.
	ImageRegistryServer *string `field:"required" json:"imageRegistryServer" yaml:"imageRegistryServer"`
	// The image registry username.
	ImageRegistryUsername *string `field:"required" json:"imageRegistryUsername" yaml:"imageRegistryUsername"`
	// The location/region the migration container should be deployed.
	Location *string `field:"required" json:"location" yaml:"location"`
	// The subnet id to deploy the migration container.
	MigrationContainerSubnetId *string `field:"required" json:"migrationContainerSubnetId" yaml:"migrationContainerSubnetId"`
	// The migration image to use.
	MigrationImageUri *string `field:"required" json:"migrationImageUri" yaml:"migrationImageUri"`
	// The name of the database.
	Name *string `field:"required" json:"name" yaml:"name"`
	// The name of the resource group.
	ResourceGroupName *string `field:"required" json:"resourceGroupName" yaml:"resourceGroupName"`
	// The id of the postgresql flexible server.
	ServerId *string `field:"required" json:"serverId" yaml:"serverId"`
	// The id of the stack.
	StackId *string `field:"required" json:"stackId" yaml:"stackId"`
}

