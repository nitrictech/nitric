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
	// The actions to apply to the policy.
	Actions *[]*string `field:"required" json:"actions" yaml:"actions"`
	// The IAM roles available to the policy.
	IamRoles interface{} `field:"required" json:"iamRoles" yaml:"iamRoles"`
	// The name of the resource.
	ResourceName *string `field:"required" json:"resourceName" yaml:"resourceName"`
	// The type of the resource (Bucket, Secret, KeyValueStore, Queue).
	ResourceType *string `field:"required" json:"resourceType" yaml:"resourceType"`
	// The service account to apply the policy to.
	ServiceAccountEmail *string `field:"required" json:"serviceAccountEmail" yaml:"serviceAccountEmail"`
}
