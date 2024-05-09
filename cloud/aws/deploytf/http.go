package deploytf

import (
	"fmt"

	"github.com/hashicorp/terraform-cdk-go/cdktf"
	deploymentspb "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"
)

func (a *NitricAwsTerraformProvider) Http(tack cdktf.TerraformStack, name string, config *deploymentspb.Http) error {
	return fmt.Errorf("nitric AWS terraform provider does not support HTTP deployment")
}
