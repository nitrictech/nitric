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

package queue

import (
	common "github.com/nitrictech/nitric/cloud/common/deploy/tags"
	v1 "github.com/nitrictech/nitric/core/pkg/api/nitric/deploy/v1"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/sqs"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type SQSQueue struct {
	pulumi.ResourceState
	Sqs  *sqs.Queue
	Name string
}

type SQSQueueArgs struct {
	Queue   *v1.Queue
	StackID pulumi.StringInput
}

func NewSQSQueue(ctx *pulumi.Context, name string, args *SQSQueueArgs, opts ...pulumi.ResourceOption) (*SQSQueue, error) {
	res := &SQSQueue{Name: name}

	err := ctx.RegisterComponentResource("nitric:queue:AwsSqsQueue", name, res, opts...)
	if err != nil {
		return nil, err
	}

	queue, err := sqs.NewQueue(ctx, name, &sqs.QueueArgs{
		Tags: common.Tags(ctx, args.StackID, name),
	})
	if err != nil {
		return nil, err
	}

	res.Sqs = queue

	return res, nil
}
