package cdn_api_rewrites

import (
	"github.com/hashicorp/terraform-cdk-go/cdktf"
)

type CdnApiRewritesConfig struct {
	// Experimental.
	DependsOn *[]cdktf.ITerraformDependable `field:"optional" json:"dependsOn" yaml:"dependsOn"`
	// Experimental.
	ForEach cdktf.ITerraformIterator `field:"optional" json:"forEach" yaml:"forEach"`
	// Experimental.
	Providers *[]interface{} `field:"optional" json:"providers" yaml:"providers"`
	// Experimental.
	SkipAssetCreationFromLocalModules *bool `field:"optional" json:"skipAssetCreationFromLocalModules" yaml:"skipAssetCreationFromLocalModules"`
	// The host name of the api.
	ApiHostName *string `field:"required" json:"apiHostName" yaml:"apiHostName"`
	// The id of the cdn frontdoor profile to use for the cdn.
	CdnFrontdoorProfileId *string `field:"required" json:"cdnFrontdoorProfileId" yaml:"cdnFrontdoorProfileId"`
	// The id of the default cdn frontdoor rule set to use for the cdn.
	CdnFrontdoorRuleSetId *string `field:"required" json:"cdnFrontdoorRuleSetId" yaml:"cdnFrontdoorRuleSetId"`
	// The name of the api.
	Name *string `field:"required" json:"name" yaml:"name"`
	// The order of the rule to use for the cdn 1.
	RuleOrder *float64 `field:"optional" json:"ruleOrder" yaml:"ruleOrder"`
}

