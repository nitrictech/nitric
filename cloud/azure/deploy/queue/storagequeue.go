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
