package deploytf

import (
	"fmt"

	"github.com/hashicorp/terraform-cdk-go/cdktf"
	deploymentspb "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"
)

// // KeyValueStore - Deploy a Key Value tioStore
func (a *NitricAwsTerraformProvider) KeyValueStore(stack cdktf.TerraformStack, name string, config *deploymentspb.KeyValueStore) error {
	return fmt.Errorf("nitric AWS terraform provider does not support KeyValueStore deployment")
}
