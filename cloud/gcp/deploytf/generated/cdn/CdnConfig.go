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
	// A map of API gateway configurations.
	ApiGateways interface{} `field:"required" json:"apiGateways" yaml:"apiGateways"`
	// The CDN domain configuration.
	CdnDomain interface{} `field:"required" json:"cdnDomain" yaml:"cdnDomain"`
	// The project ID where resources will be created.
	ProjectId *string `field:"required" json:"projectId" yaml:"projectId"`
	// The region where resources will be created.
	Region *string `field:"required" json:"region" yaml:"region"`
	// The unique identifier for the stack.
	StackId *string `field:"required" json:"stackId" yaml:"stackId"`
	// A map of website bucket configurations.
	WebsiteBuckets interface{} `field:"required" json:"websiteBuckets" yaml:"websiteBuckets"`
}

