package deploy

import (
	deploymentspb "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func (p *NitricGcpPulumiProvider) KeyValueStore(ctx *pulumi.Context, parent pulumi.Resource, name string, config *deploymentspb.KeyValueStore) error {
	// keyvalue stores are created at runtime in GCP, so we don't need to do anything here.
	return nil
}
