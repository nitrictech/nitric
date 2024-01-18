package deploy

import (
	deploymentspb "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func (p *NitricAzurePulumiProvider) Secret(ctx *pulumi.Context, parent pulumi.Resource, name string, config *deploymentspb.Secret) error {
	// Secrets in Azure Key Vaults are unique resources created during deployment, so we don't need to do anything here.
	// Instead, if at least one secret is requested a Key Vault will be created for the stack.
	// Policies are also created which restrict access to the Key Vault and the named secrets inside.
	return nil
}
