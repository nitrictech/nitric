package deploy

import (
	deploymentspb "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"
	"github.com/pulumi/pulumi-azure-native-sdk/storage"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func (p *NitricAzurePulumiProvider) KeyValueStore(ctx *pulumi.Context, parent pulumi.Resource, name string, config *deploymentspb.KeyValueStore) error {
	var err error
	opts := []pulumi.ResourceOption{pulumi.Parent(parent)}

	p.keyValueStores[name], err = storage.NewTable(ctx, name, &storage.TableArgs{
		AccountName:       p.storageAccount.Name,
		ResourceGroupName: p.resourceGroup.Name,
		TableName:         pulumi.String(name),
	}, opts...)

	return err
}
