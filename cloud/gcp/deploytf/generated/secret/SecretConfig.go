package secret

import (
	"github.com/hashicorp/terraform-cdk-go/cdktf"
)

type SecretConfig struct {
	// Experimental.
	DependsOn *[]cdktf.ITerraformDependable `field:"optional" json:"dependsOn" yaml:"dependsOn"`
	// Experimental.
	ForEach cdktf.ITerraformIterator `field:"optional" json:"forEach" yaml:"forEach"`
	// Experimental.
	Providers *[]interface{} `field:"optional" json:"providers" yaml:"providers"`
	// Experimental.
	SkipAssetCreationFromLocalModules *bool `field:"optional" json:"skipAssetCreationFromLocalModules" yaml:"skipAssetCreationFromLocalModules"`
	// location of the secret.
	Location *string `field:"required" json:"location" yaml:"location"`
	// The name of the secret.
	SecretName *string `field:"required" json:"secretName" yaml:"secretName"`
	// The ID of the Nitric stack.
	StackId *string `field:"required" json:"stackId" yaml:"stackId"`
	// The KMS key to use for encryption.
	CmekKey *string `field:"optional" json:"cmekKey" yaml:"cmekKey"`
}

