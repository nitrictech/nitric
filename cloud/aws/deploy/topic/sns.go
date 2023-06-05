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

package topic

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/iam"
	awslambda "github.com/pulumi/pulumi-aws/sdk/v5/go/aws/lambda"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/sfn"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/sns"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"

	"github.com/nitrictech/nitric/cloud/aws/deploy/exec"
	common "github.com/nitrictech/nitric/cloud/common/deploy/tags"
	v1 "github.com/nitrictech/nitric/core/pkg/api/nitric/deploy/v1"
)

type SNSTopic struct {
	pulumi.ResourceState
	Name string
	Sns  *sns.Topic
	Sfn  *sfn.StateMachine
}

type SNSTopicArgs struct {
	StackID string
	Topic   *v1.Topic
}

func NewSNSTopic(ctx *pulumi.Context, name string, args *SNSTopicArgs, opts ...pulumi.ResourceOption) (*SNSTopic, error) {
	res := &SNSTopic{Name: name}

	err := ctx.RegisterComponentResource("nitric:topic:AwsSnsTopic", name, res, opts...)
	if err != nil {
		return nil, err
	}

	// create the SNS topic
	res.Sns, err = sns.NewTopic(ctx, name, &sns.TopicArgs{
		Tags: pulumi.ToStringMap(common.Tags(args.StackID, name)),
	}, pulumi.Parent(res))
	if err != nil {
		return nil, err
	}

	// create a State Machine to support delayed messaging
	// unfortunately we cannot create a single dynamic state machine that uses
	// the topicArn as input so we need to create one per topic
	// Note this is going to be better for security
	r, _ := json.Marshal(map[string]interface{}{
		"Version": "2012-10-17",
		"Statement": []map[string]interface{}{
			{
				"Sid":    "",
				"Effect": "Allow",
				"Principal": map[string]interface{}{
					"Service": "states.amazonaws.com",
				},
				"Action": "sts:AssumeRole",
			},
		},
	})

	sfnRole, err := iam.NewRole(ctx, fmt.Sprintf("%s-delay-ctrl", name), &iam.RoleArgs{
		AssumeRolePolicy: pulumi.String(r),
	}, pulumi.Parent(res))
	if err != nil {
		return nil, errors.WithMessage(err, "topic delay controller role")
	}

	policy := res.Sns.Arn.ApplyT(func(arn string) (string, error) {
		rp, err := json.Marshal(map[string]interface{}{
			"Version": "2012-10-17",
			"Statement": []map[string]interface{}{
				{
					"Sid":      "",
					"Effect":   "Allow",
					"Action":   []string{"sns:Publish"},
					"Resource": arn,
				},
			},
		})

		return string(rp), err
	})

	// Enable a role with publish access to this stacks topics only
	_, err = iam.NewRolePolicy(ctx, fmt.Sprintf("%s-delay-ctrl", name), &iam.RolePolicyArgs{
		Role: sfnRole,
		// TODO: Limit to only this stacks topics (deployed above)
		Policy: policy,
	}, pulumi.Parent(res))
	if err != nil {
		return nil, errors.WithMessage(err, "topic delay controller role policy")
	}

	sfnDef := res.Sns.Arn.ApplyT(func(arn string) (string, error) {
		def, err := json.Marshal(map[string]interface{}{
			"Comment": "",
			"StartAt": "Wait",
			"States": map[string]interface{}{
				"Wait": map[string]string{
					"Type":        "Wait",
					"SecondsPath": "$.seconds",
					"Next":        "Publish",
				},
				"Publish": map[string]interface{}{
					"Type":     "Task",
					"Resource": "arn:aws:states:::sns:publish",
					"Parameters": map[string]string{
						"TopicArn":  arn,
						"Message.$": "$.message",
					},
					"End": true,
				},
			},
		})

		return string(def), err
	}).(pulumi.StringOutput)

	// Deploy a delay manager using AWS step functions
	// This will enable runtime delaying of event
	res.Sfn, err = sfn.NewStateMachine(ctx, name, &sfn.StateMachineArgs{
		RoleArn: sfnRole.Arn,
		// Apply the same name as the topic to the state machine
		Tags:       pulumi.ToStringMap(common.Tags(args.StackID, name)),
		Definition: sfnDef,
	}, pulumi.Parent(res))
	if err != nil {
		return nil, errors.WithMessage(err, "topic delay controller")
	}

	return res, nil
}

type SNSTopicSubscription struct {
	pulumi.ResourceState
	Name string
	Sns  *sns.TopicSubscription
}

type SNSTopicSubscriptionArgs struct {
	Name   string
	Topic  *SNSTopic
	Lambda *exec.LambdaExecUnit
}

// Create a new subscription for an SNS Topic
func NewSNSTopicSubscription(ctx *pulumi.Context, name string, args *SNSTopicSubscriptionArgs, opts ...pulumi.ResourceOption) (*SNSTopicSubscription, error) {
	res := &SNSTopicSubscription{Name: name}

	err := ctx.RegisterComponentResource("nitric:topic:AwsSnsTopicSubscription", name, res, opts...)
	if err != nil {
		return nil, err
	}

	_, err = awslambda.NewPermission(ctx, name+"Permission", &awslambda.PermissionArgs{
		SourceArn: args.Topic.Sns.Arn,
		Function:  args.Lambda.Function.Name,
		Principal: pulumi.String("sns.amazonaws.com"),
		Action:    pulumi.String("lambda:InvokeFunction"),
	}, opts...)
	if err != nil {
		return nil, err
	}

	res.Sns, err = sns.NewTopicSubscription(ctx, name+"Subscription", &sns.TopicSubscriptionArgs{
		Endpoint: args.Lambda.Function.Arn,
		Protocol: pulumi.String("lambda"),
		Topic:    args.Topic.Sns.ID(),
	}, opts...)
	if err != nil {
		return nil, err
	}

	return res, nil
}
