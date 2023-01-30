// Copyright Nitric Pty Ltd.
//
// SPDX-License-Identifier: Apache-2.0
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package storage

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-azure-native-sdk/storage"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"

	"github.com/nitrictech/nitric/cloud/azure/deploy/utils"
	common "github.com/nitrictech/nitric/cloud/common/deploy/tags"
)

type StorageArgs struct {
	ResourceGroupName pulumi.StringInput
	StackID           pulumi.StringInput
	BucketNames       []string
	QueueNames        []string
}

type Storage struct {
	pulumi.ResourceState

	Name       string
	Account    *storage.StorageAccount
	Queues     map[string]*storage.Queue
	Containers map[string]*storage.BlobContainer
}

func NewStorageResources(ctx *pulumi.Context, name string, args *StorageArgs, opts ...pulumi.ResourceOption) (*Storage, error) {
	res := &Storage{
		Name:       name,
		Queues:     map[string]*storage.Queue{},
		Containers: map[string]*storage.BlobContainer{},
	}

	err := ctx.RegisterComponentResource("nitric:storage:AzureStorage", name, res, opts...)
	if err != nil {
		return nil, err
	}

	accName := utils.ResourceName(ctx, name, utils.StorageAccountRT)

	res.Account, err = storage.NewStorageAccount(ctx, accName, &storage.StorageAccountArgs{
		AccessTier:        storage.AccessTierHot,
		ResourceGroupName: args.ResourceGroupName,
		Kind:              pulumi.String("StorageV2"),
		Sku: storage.SkuArgs{
			Name: pulumi.String(storage.SkuName_Standard_LRS),
		},
		Tags: common.Tags(ctx, args.StackID, accName),
	}, pulumi.Parent(res))
	if err != nil {
		return nil, errors.WithMessage(err, "account create")
	}

	for _, bName := range args.BucketNames {
		res.Containers[bName], err = storage.NewBlobContainer(ctx, utils.ResourceName(ctx, bName, utils.StorageContainerRT), &storage.BlobContainerArgs{
			ResourceGroupName: args.ResourceGroupName,
			AccountName:       res.Account.Name,
		}, pulumi.Parent(res))
		if err != nil {
			return nil, errors.WithMessage(err, "container create")
		}
	}

	for _, qName := range args.QueueNames {
		res.Queues[qName], err = storage.NewQueue(ctx, utils.ResourceName(ctx, qName, utils.StorageQueueRT), &storage.QueueArgs{
			ResourceGroupName: args.ResourceGroupName,
			AccountName:       res.Account.Name,
		}, pulumi.Parent(res))
		if err != nil {
			return nil, errors.WithMessage(err, "queue create")
		}
	}

	return res, nil
}
