package queue

import (
	"github.com/nitrictech/nitric/cloud/azure/deploy/utils"
	"github.com/pulumi/pulumi-azure-native-sdk/resources"
	"github.com/pulumi/pulumi-azure-native-sdk/storage"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Topics
type AzureStorageQueue struct {
	pulumi.ResourceState

	Name  string
	Queue *storage.Queue
}

type AzureStorageQueueArgs struct {
	StackID       pulumi.StringInput
	Account       *storage.StorageAccount
	ResourceGroup *resources.ResourceGroup
}

func NewAzureStorageQueue(ctx *pulumi.Context, name string, args *AzureStorageQueueArgs, opts ...pulumi.ResourceOption) (*AzureStorageQueue, error) {
	res := &AzureStorageQueue{Name: name}

	err := ctx.RegisterComponentResource("nitric:queue:AzureStorageQueue", name, res, opts...)
	if err != nil {
		return nil, err
	}

	res.Queue, err = storage.NewQueue(ctx, utils.ResourceName(ctx, name, utils.StorageQueueRT), &storage.QueueArgs{
		AccountName:       args.Account.Name,
		ResourceGroupName: args.ResourceGroup.Name,
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}
