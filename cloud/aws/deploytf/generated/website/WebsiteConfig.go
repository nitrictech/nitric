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
	// The production website output directory.
	LocalDirectory *string `field:"required" json:"localDirectory" yaml:"localDirectory"`
	// The name of the website.
	Name *string `field:"required" json:"name" yaml:"name"`
	// The unique ID for this stack.
	StackId *string `field:"required" json:"stackId" yaml:"stackId"`
}

