package stack

import (
	"github.com/hashicorp/terraform-cdk-go/cdktf"
)

type StackConfig struct {
	// Experimental.
	DependsOn *[]cdktf.ITerraformDependable `field:"optional" json:"dependsOn" yaml:"dependsOn"`
	// Experimental.
	ForEach cdktf.ITerraformIterator `field:"optional" json:"forEach" yaml:"forEach"`
	// Experimental.
	Providers *[]interface{} `field:"optional" json:"providers" yaml:"providers"`
	// Experimental.
	SkipAssetCreationFromLocalModules *bool `field:"optional" json:"skipAssetCreationFromLocalModules" yaml:"skipAssetCreationFromLocalModules"`
	// The location/region where the resources will be created.
	Location *string `field:"required" json:"location" yaml:"location"`
	// Deploy compatible services with private endpoints.
	PrivateEndpoints *bool `field:"required" json:"privateEndpoints" yaml:"privateEndpoints"`
	// The name of the stack.
	StackName *string `field:"required" json:"stackName" yaml:"stackName"`
	// The id of the subnet to deploy the infrastructure resources.
	SubnetId *string `field:"required" json:"subnetId" yaml:"subnetId"`
	// The id of the subscription to deploy the stack.
	SubscriptionId *string `field:"required" json:"subscriptionId" yaml:"subscriptionId"`
	// The tags to apply to the stack The property type contains a map, they have special handling, please see {@link cdk.tf /module-map-inputs the docs}.
	Tags *map[string]*string `field:"required" json:"tags" yaml:"tags"`
	// The name of the vnet to deploy the infrastructure resources.
	VnetId *string `field:"required" json:"vnetId" yaml:"vnetId"`
	// Whether to create private DNS zones for private endpoints.
	CreateDnsZones *bool `field:"optional" json:"createDnsZones" yaml:"createDnsZones"`
	// Enable the creation of a database.
	EnableDatabase *bool `field:"optional" json:"enableDatabase" yaml:"enableDatabase"`
	// Enable the creation of a keyvault.
	EnableKeyvault *bool `field:"optional" json:"enableKeyvault" yaml:"enableKeyvault"`
	// Enable the creation of a storage account.
	EnableStorage *bool `field:"optional" json:"enableStorage" yaml:"enableStorage"`
	// The name of the resource group to reuse.
	ResourceGroupName *string `field:"optional" json:"resourceGroupName" yaml:"resourceGroupName"`
}

