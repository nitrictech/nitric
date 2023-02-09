package bucket

import (
	"github.com/nitrictech/nitric/cloud/azure/deploy/utils"
	"github.com/pulumi/pulumi-azure-native-sdk/resources"
	"github.com/pulumi/pulumi-azure-native-sdk/storage"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Topics
type AzureStorageBucket struct {
	pulumi.ResourceState

	Name      string
	Container *storage.BlobContainer
}

type AzureStorageBucketArgs struct {
	StackID       pulumi.StringInput
	Account       *storage.StorageAccount
	ResourceGroup *resources.ResourceGroup
}

func NewAzureStorageBucket(ctx *pulumi.Context, name string, args *AzureStorageBucketArgs, opts ...pulumi.ResourceOption) (*AzureStorageBucket, error) {
	res := &AzureStorageBucket{Name: name}

	err := ctx.RegisterComponentResource("nitric:bucket:AzureStorageBucket", name, res, opts...)
	if err != nil {
		return nil, err
	}

	res.Container, err = storage.NewBlobContainer(ctx, utils.ResourceName(ctx, name, utils.StorageContainerRT), &storage.BlobContainerArgs{
		ResourceGroupName: args.ResourceGroup.Name,
		AccountName:       args.Account.Name,
	}, pulumi.Parent(res))
	if err != nil {
		return nil, err
	}

	return res, nil
}
