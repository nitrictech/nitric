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
	// The name of the stack.
	StackName *string `field:"required" json:"stackName" yaml:"stackName"`
	// Enable the creation of a database.
	EnableDatabase *bool `field:"optional" json:"enableDatabase" yaml:"enableDatabase"`
	// Enable the creation of a keyvault.
	EnableKeyvault *bool `field:"optional" json:"enableKeyvault" yaml:"enableKeyvault"`
	// Enable the creation of a storage account.
	EnableStorage *bool `field:"optional" json:"enableStorage" yaml:"enableStorage"`
	// Enable the creation of a website.
	EnableWebsite *bool `field:"optional" json:"enableWebsite" yaml:"enableWebsite"`
	// The id of the subnet to deploy the infrastructure resources.
	InfrastructureSubnetId *string `field:"optional" json:"infrastructureSubnetId" yaml:"infrastructureSubnetId"`
	// The root error document for the website 404.html.
	WebsiteRootErrorDocument *string `field:"optional" json:"websiteRootErrorDocument" yaml:"websiteRootErrorDocument"`
	// The root index document for the website index.html.
	WebsiteRootIndexDocument *string `field:"optional" json:"websiteRootIndexDocument" yaml:"websiteRootIndexDocument"`
}

