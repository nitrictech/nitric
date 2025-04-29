package cdn_subsites

import (
	"github.com/hashicorp/terraform-cdk-go/cdktf"
)

type CdnSubsitesConfig struct {
	// Experimental.
	DependsOn *[]cdktf.ITerraformDependable `field:"optional" json:"dependsOn" yaml:"dependsOn"`
	// Experimental.
	ForEach cdktf.ITerraformIterator `field:"optional" json:"forEach" yaml:"forEach"`
	// Experimental.
	Providers *[]interface{} `field:"optional" json:"providers" yaml:"providers"`
	// Experimental.
	SkipAssetCreationFromLocalModules *bool `field:"optional" json:"skipAssetCreationFromLocalModules" yaml:"skipAssetCreationFromLocalModules"`
	// The host name of the api.
	BasePath *string `field:"required" json:"basePath" yaml:"basePath"`
	// The id of the default cdn frontdoor rule set to use for the cdn.
	CdnDefaultFrontdoorRuleSetId *string `field:"required" json:"cdnDefaultFrontdoorRuleSetId" yaml:"cdnDefaultFrontdoorRuleSetId"`
	// The id of the cdn frontdoor profile to use for the cdn.
	CdnFrontdoorProfileId *string `field:"required" json:"cdnFrontdoorProfileId" yaml:"cdnFrontdoorProfileId"`
	// The name of the api.
	Name *string `field:"required" json:"name" yaml:"name"`
	// The primary host for the website.
	PrimaryWebHost *string `field:"required" json:"primaryWebHost" yaml:"primaryWebHost"`
	// The id of the stack.
	StackId *string `field:"required" json:"stackId" yaml:"stackId"`
	// The order of the rule to use for the cdn 1.
	RuleOrder *float64 `field:"optional" json:"ruleOrder" yaml:"ruleOrder"`
}

