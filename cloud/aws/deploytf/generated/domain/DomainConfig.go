package domain

import (
	"github.com/hashicorp/terraform-cdk-go/cdktf"
)

type DomainConfig struct {
	// Experimental.
	DependsOn *[]cdktf.ITerraformDependable `field:"optional" json:"dependsOn" yaml:"dependsOn"`
	// Experimental.
	ForEach cdktf.ITerraformIterator `field:"optional" json:"forEach" yaml:"forEach"`
	// Experimental.
	Providers *[]interface{} `field:"optional" json:"providers" yaml:"providers"`
	// Experimental.
	SkipAssetCreationFromLocalModules *bool `field:"optional" json:"skipAssetCreationFromLocalModules" yaml:"skipAssetCreationFromLocalModules"`
	// The ID of the API.
	ApiId *string `field:"required" json:"apiId" yaml:"apiId"`
	// The stage for the API.
	ApiStageName *string `field:"required" json:"apiStageName" yaml:"apiStageName"`
	// The name of the domain.
	//
	// This must be globally unique.
	DomainName *string `field:"required" json:"domainName" yaml:"domainName"`
	// The ID of the Nitric stack.
	StackId *string `field:"required" json:"stackId" yaml:"stackId"`
}

