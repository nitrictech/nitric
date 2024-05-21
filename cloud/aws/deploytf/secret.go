package deploytf

import (
	"github.com/aws/jsii-runtime-go"
	"github.com/hashicorp/terraform-cdk-go/cdktf"
	"github.com/nitrictech/nitric/cloud/aws/deploytf/generated/secret"
	deploymentspb "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"
)

// // Secret - Deploy a Secret
func (a *NitricAwsTerraformProvider) Secret(stack cdktf.TerraformStack, name string, config *deploymentspb.Secret) error {
	a.Secrets[name] = secret.NewSecret(stack, jsii.String(name), &secret.SecretConfig{
		SecretName: jsii.String(name),
		StackId:    a.Stack.StackIdOutput(),
	})

	return nil
}
