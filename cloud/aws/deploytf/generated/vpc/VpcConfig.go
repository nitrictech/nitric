package vpc

import (
	"github.com/hashicorp/terraform-cdk-go/cdktf"
)

type VpcConfig struct {
	// Experimental.
	DependsOn *[]cdktf.ITerraformDependable `field:"optional" json:"dependsOn" yaml:"dependsOn"`
	// Experimental.
	ForEach cdktf.ITerraformIterator `field:"optional" json:"forEach" yaml:"forEach"`
	// Experimental.
	Providers *[]interface{} `field:"optional" json:"providers" yaml:"providers"`
	// Experimental.
	SkipAssetCreationFromLocalModules *bool `field:"optional" json:"skipAssetCreationFromLocalModules" yaml:"skipAssetCreationFromLocalModules"`
	// The CIDR block for the VPC 10.0.0.0/16.
	CidrBlock *string `field:"optional" json:"cidrBlock" yaml:"cidrBlock"`
	// Private Subnet CIDR values 10.0.4.0/24 10.0.5.0/24 10.0.6.0/24.
	PrivateSubnetCidrs *[]*string `field:"optional" json:"privateSubnetCidrs" yaml:"privateSubnetCidrs"`
	// Public Subnet CIDR values 10.0.1.0/24 10.0.2.0/24 10.0.3.0/24.
	PublicSubnetCidrs *[]*string `field:"optional" json:"publicSubnetCidrs" yaml:"publicSubnetCidrs"`
}
