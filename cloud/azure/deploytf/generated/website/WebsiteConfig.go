package website

import (
	"github.com/hashicorp/terraform-cdk-go/cdktf"
)

type WebsiteConfig struct {
	// Experimental.
	DependsOn *[]cdktf.ITerraformDependable `field:"optional" json:"dependsOn" yaml:"dependsOn"`
	// Experimental.
	ForEach cdktf.ITerraformIterator `field:"optional" json:"forEach" yaml:"forEach"`
	// Experimental.
	Providers *[]interface{} `field:"optional" json:"providers" yaml:"providers"`
	// Experimental.
	SkipAssetCreationFromLocalModules *bool `field:"optional" json:"skipAssetCreationFromLocalModules" yaml:"skipAssetCreationFromLocalModules"`
	// The base path for the website.
	BasePath *string `field:"required" json:"basePath" yaml:"basePath"`
	// The local directory to deploy the website from.
	LocalDirectory *string `field:"required" json:"localDirectory" yaml:"localDirectory"`
	// The name of the website.
	Name *string `field:"required" json:"name" yaml:"name"`
	// The connection string for the storage account.
	StorageAccountConnectionString *string `field:"required" json:"storageAccountConnectionString" yaml:"storageAccountConnectionString"`
	// The name of the storage account.
	StorageAccountName *string `field:"required" json:"storageAccountName" yaml:"storageAccountName"`
}

