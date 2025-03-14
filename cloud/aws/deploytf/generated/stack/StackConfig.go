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
	// Enable the creation of a website.
	EnableWebsite *bool `field:"optional" json:"enableWebsite" yaml:"enableWebsite"`
	// The root error document for the website 404.html.
	WebsiteRootErrorDocument *string `field:"optional" json:"websiteRootErrorDocument" yaml:"websiteRootErrorDocument"`
	// The root index document for the website index.html.
	WebsiteRootIndexDocument *string `field:"optional" json:"websiteRootIndexDocument" yaml:"websiteRootIndexDocument"`
}

