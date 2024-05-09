package deploytf

import (
	"fmt"

	"github.com/hashicorp/terraform-cdk-go/cdktf"
	deploymentspb "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"
)

func (n *NitricAwsTerraformProvider) Api(cdktf.TerraformStack, string, *deploymentspb.Api) error {
	return fmt.Errorf("nitric AWS terraform provider does not support API deployment")
}
