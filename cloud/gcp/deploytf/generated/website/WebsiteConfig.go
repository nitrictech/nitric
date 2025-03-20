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
	// The base path for the website files.
	BasePath *string `field:"required" json:"basePath" yaml:"basePath"`
	// The local directory containing website files.
	LocalDirectory *string `field:"required" json:"localDirectory" yaml:"localDirectory"`
	// The region where the bucket will be created.
	Region *string `field:"required" json:"region" yaml:"region"`
	// The ID of the stack.
	StackId *string `field:"required" json:"stackId" yaml:"stackId"`
	// The name of the website.
	WebsiteName *string `field:"required" json:"websiteName" yaml:"websiteName"`
	// The error document for the website.
	//
	// 404.html
	ErrorDocument *string `field:"optional" json:"errorDocument" yaml:"errorDocument"`
	// The index document for the website.
	//
	// index.html
	IndexDocument *string `field:"optional" json:"indexDocument" yaml:"indexDocument"`
}

