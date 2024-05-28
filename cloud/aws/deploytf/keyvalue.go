package deploytf

import (
	"fmt"

	"github.com/aws/jsii-runtime-go"
	"github.com/hashicorp/terraform-cdk-go/cdktf"
	"github.com/nitrictech/nitric/cloud/aws/deploytf/generated/keyvalue"
	deploymentspb "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"
)

// // KeyValueStore - Deploy a Key Value tioStore
func (a *NitricAwsTerraformProvider) KeyValueStore(stack cdktf.TerraformStack, name string, config *deploymentspb.KeyValueStore) error {
	a.KeyValueStores[name] = keyvalue.NewKeyvalue(stack, jsii.String(name), &keyvalue.KeyvalueConfig{
		KvstoreName: jsii.String(name),
		StackId:     a.Stack.StackIdOutput(),
	})

	return fmt.Errorf("nitric AWS terraform provider does not support KeyValueStore deployment")
}
