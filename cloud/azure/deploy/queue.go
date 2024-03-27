package deploy

import (
	deploymentspb "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"
	"github.com/pulumi/pulumi-azure-native-sdk/storage"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func (a *NitricAzurePulumiProvider) Queue(ctx *pulumi.Context, parent pulumi.Resource, name string, config *deploymentspb.Queue) error {
	var err error
	opts := []pulumi.ResourceOption{pulumi.Parent(parent)}

	a.Queues[name], err = storage.NewQueue(ctx, ResourceName(ctx, name, StorageQueueRT), &storage.QueueArgs{
		AccountName:       a.StorageAccount.Name,
		ResourceGroupName: a.ResourceGroup.Name,
	}, opts...)

	return err
}

// Copyright 2021 Nitric Technologies Pty Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// package queue

// import (
// 	"github.com/nitrictech/nitric/cloud/azure/deploy/utils"
// 	"github.com/pulumi/pulumi-azure-native-sdk/resources"
// 	"github.com/pulumi/pulumi-azure-native-sdk/storage"
// 	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
// )

// // Topics
// type AzureStorageQueue struct {
// 	pulumi.ResourceState

// 	Name          string
// 	Account       *storage.StorageAccount
// 	ResourceGroup *resources.ResourceGroup
// 	Queue         *storage.Queue
// }

// type AzureStorageQueueArgs struct {
// 	Account       *storage.StorageAccount
// 	ResourceGroup *resources.ResourceGroup
// }

// func NewAzureStorageQueue(ctx *pulumi.Context, name string, args *AzureStorageQueueArgs, opts ...pulumi.ResourceOption) (*AzureStorageQueue, error) {
// 	res := &AzureStorageQueue{
// 		Name:          name,
// 		Account:       args.Account,
// 		ResourceGroup: args.ResourceGroup,
// 	}

// 	err := ctx.RegisterComponentResource("nitric:queue:AzureStorageQueue", name, res, opts...)
// 	if err != nil {
// 		return nil, err
// 	}

// 	res.Queue, err = storage.NewQueue(ctx, utils.ResourceName(ctx, name, utils.StorageQueueRT), &storage.QueueArgs{
// 		AccountName:       args.Account.Name,
// 		ResourceGroupName: args.ResourceGroup.Name,
// 	})
// 	if err != nil {
// 		return nil, err
// 	}

// 	return res, nil
// }
