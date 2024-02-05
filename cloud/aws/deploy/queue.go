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

package deploy

import (
	"github.com/nitrictech/nitric/cloud/common/deploy/resources"
	common "github.com/nitrictech/nitric/cloud/common/deploy/tags"
	deploymentspb "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/sqs"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Bucket - Implements deployments of Nitric Buckets using AWS S3
func (a *NitricAwsPulumiProvider) Queue(ctx *pulumi.Context, parent pulumi.Resource, name string, config *deploymentspb.Queue) error {
	opts := []pulumi.ResourceOption{pulumi.Parent(parent)}

	queue, err := sqs.NewQueue(ctx, name, &sqs.QueueArgs{
		Tags: pulumi.ToStringMap(common.Tags(a.stackId, name, resources.Queue)),
	}, opts...)
	if err != nil {
		return err
	}

	a.queues[name] = queue

	return nil
}
