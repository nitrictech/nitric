package deploytf

import (
	"fmt"

	"github.com/hashicorp/terraform-cdk-go/cdktf"
	deploymentspb "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"
)

// // Schedule - Deploy a Schedule
func (a *NitricAwsTerraformProvider) Schedule(stack cdktf.TerraformStack, name string, config *deploymentspb.Schedule) error {
	return fmt.Errorf("nitric AWS terraform provider does not support Schedule deployment")
}
