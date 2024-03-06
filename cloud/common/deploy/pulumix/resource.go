package pulumix

import (
	deploymentspb "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"
	resourcespb "github.com/nitrictech/nitric/core/pkg/proto/resources/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type NitricPulumiResource[T any] struct {
	Id     *resourcespb.ResourceIdentifier
	Config T
}

type NitricPulumiServiceResource = NitricPulumiResource[*NitricPulumiServiceConfig]

type NitricPulumiServiceConfig struct {
	*deploymentspb.Service

	// Allow for pulumi strings to be used as service environment variables
	env pulumi.StringMap
}

func (n *NitricPulumiServiceConfig) SetEnv(key string, value pulumi.StringInput) {
	if n.env == nil {
		n.env = pulumi.StringMap{}
	}

	n.env[key] = value
}

func (n *NitricPulumiServiceConfig) Env() pulumi.StringMap {
	envMap := pulumi.StringMap{}

	for k, v := range n.Service.GetEnv() {
		envMap[k] = pulumi.String(v)
	}

	for k, v := range n.env {
		envMap[k] = v
	}

	return envMap
}
