package deploytf

import (
	"fmt"

	"github.com/hashicorp/terraform-cdk-go/cdktf"
	"github.com/nitrictech/nitric/cloud/common/deploy/provider"
	deploymentspb "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"
)

func (a *NitricAwsTerraformProvider) Service(stack cdktf.TerraformStack, name string, config *deploymentspb.Service, runtimeProvider provider.RuntimeProvider) error {
	return fmt.Errorf("nitric AWS terraform provider does not support Service deployment")
}
