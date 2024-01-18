package deploy

import (
	"fmt"

	deploymentspb "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func (a *NitricGcpPulumiProvider) Websocket(ctx *pulumi.Context, parent pulumi.Resource, name string, config *deploymentspb.Websocket) error {
	return fmt.Errorf("Websockets not implemented for GCP")
}
