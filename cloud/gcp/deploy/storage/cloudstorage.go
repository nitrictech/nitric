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
	"github.com/nitrictech/nitric/cloud/common/deploy/resources"
	common "github.com/nitrictech/nitric/cloud/common/deploy/tags"
	v1 "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/storage"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type CloudStorageBucket struct {
	pulumi.ResourceState

	Name         string
	CloudStorage *storage.Bucket
}

type CloudStorageBucketArgs struct {
	Location string
	StackID  string

	Bucket *v1.Bucket
}

func NewCloudStorageBucket(ctx *pulumi.Context, name string, args *CloudStorageBucketArgs, opts ...pulumi.ResourceOption) (*CloudStorageBucket, error) {
	res := &CloudStorageBucket{
		Name: name,
	}

	err := ctx.RegisterComponentResource("nitric:bucket:GCPCloudStorage", name, res, opts...)
	if err != nil {
		return nil, err
	}

	res.CloudStorage, err = storage.NewBucket(ctx, name, &storage.BucketArgs{
		Location: pulumi.String(args.Location),
		Labels:   pulumi.ToStringMap(common.Tags(args.StackID, name, resources.Bucket)),
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}
