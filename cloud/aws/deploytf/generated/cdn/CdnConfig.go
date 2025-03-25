package cdn

import (
	"github.com/hashicorp/terraform-cdk-go/cdktf"
)

type CdnConfig struct {
	// Experimental.
	DependsOn *[]cdktf.ITerraformDependable `field:"optional" json:"dependsOn" yaml:"dependsOn"`
	// Experimental.
	ForEach cdktf.ITerraformIterator `field:"optional" json:"forEach" yaml:"forEach"`
	// Experimental.
	Providers *[]interface{} `field:"optional" json:"providers" yaml:"providers"`
	// Experimental.
	SkipAssetCreationFromLocalModules *bool `field:"optional" json:"skipAssetCreationFromLocalModules" yaml:"skipAssetCreationFromLocalModules"`
	// information about the root website for default behaviour.
	RootWebsite interface{} `field:"required" json:"rootWebsite" yaml:"rootWebsite"`
	// The name of the stack.
	StackName *string `field:"required" json:"stackName" yaml:"stackName"`
	// Map of websites and their storage information.
	Websites interface{} `field:"required" json:"websites" yaml:"websites"`
	// Map of APIs and their gateway information.
	Apis interface{} `field:"optional" json:"apis" yaml:"apis"`
}

