package policy

import (
	"github.com/hashicorp/terraform-cdk-go/cdktf"
)

type PolicyConfig struct {
	// Experimental.
	DependsOn *[]cdktf.ITerraformDependable `field:"optional" json:"dependsOn" yaml:"dependsOn"`
	// Experimental.
	ForEach cdktf.ITerraformIterator `field:"optional" json:"forEach" yaml:"forEach"`
	// Experimental.
	Providers *[]interface{} `field:"optional" json:"providers" yaml:"providers"`
	// Experimental.
	SkipAssetCreationFromLocalModules *bool `field:"optional" json:"skipAssetCreationFromLocalModules" yaml:"skipAssetCreationFromLocalModules"`
	// actions to allow.
	Actions *[]*string `field:"required" json:"actions" yaml:"actions"`
	// principals (roles) to apply the policies to.
	Principals *[]*string `field:"required" json:"principals" yaml:"principals"`
	// resources to apply the policies to.
	Resources *[]*string `field:"required" json:"resources" yaml:"resources"`
}

