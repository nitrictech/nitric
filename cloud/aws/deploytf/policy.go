package deploytf

import (
	"github.com/hashicorp/terraform-cdk-go/cdktf"
	deploymentspb "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"
)

// // Policy - Deploy a Policy
func (a *NitricAwsTerraformProvider) Policy(stack cdktf.TerraformStack, name string, config *deploymentspb.Policy) error {
	// return fmt.Errorf("nitric AWS terraform provider does not support Policy deployment")
	return nil
}
