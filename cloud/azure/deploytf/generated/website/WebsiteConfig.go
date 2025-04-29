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
	// The location/region where the resources will be created.
	Location *string `field:"required" json:"location" yaml:"location"`
	// The name of the resource group to use for the cdn.
	ResourceGroupName *string `field:"required" json:"resourceGroupName" yaml:"resourceGroupName"`
	// The id of the stack.
	StackId *string `field:"required" json:"stackId" yaml:"stackId"`
	// The error document for the website 404.html.
	ErrorDocument *string `field:"optional" json:"errorDocument" yaml:"errorDocument"`
	// The index document for the website index.html.
	IndexDocument *string `field:"optional" json:"indexDocument" yaml:"indexDocument"`
}

