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
	// The role definition id.
	RoleDefinitionId *string `field:"required" json:"roleDefinitionId" yaml:"roleDefinitionId"`
	// The scope of the role assignment.
	Scope *string `field:"required" json:"scope" yaml:"scope"`
	// The service principal id.
	ServicePrincipalId *string `field:"required" json:"servicePrincipalId" yaml:"servicePrincipalId"`
}

