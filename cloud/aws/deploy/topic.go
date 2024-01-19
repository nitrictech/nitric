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
	"encoding/json"
	"fmt"

	"github.com/nitrictech/nitric/cloud/common/deploy/resources"
	common "github.com/nitrictech/nitric/cloud/common/deploy/tags"
	deploymentspb "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/iam"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/lambda"
	awslambda "github.com/pulumi/pulumi-aws/sdk/v5/go/aws/lambda"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/sfn"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/sns"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type topic struct {
	sns *sns.Topic
	sfn *sfn.StateMachine
}

func createSubscription(ctx *pulumi.Context, parent pulumi.Resource, name string, topic *sns.Topic, target *lambda.Function) error {
	var err error

	opts := []pulumi.ResourceOption{pulumi.Parent(parent)}

	_, err = awslambda.NewPermission(ctx, name+"Permission", &awslambda.PermissionArgs{
		SourceArn: topic.Arn,
		Function:  target.Name,
		Principal: pulumi.String("sns.amazonaws.com"),
		Action:    pulumi.String("lambda:InvokeFunction"),
	}, opts...)
	if err != nil {
		return err
	}

	_, err = sns.NewTopicSubscription(ctx, name+"Subscription", &sns.TopicSubscriptionArgs{
		Endpoint: target.Arn,
		Protocol: pulumi.String("lambda"),
		Topic:    topic.ID(),
	}, opts...)
	if err != nil {
		return err
	}

	return nil
}

func (a *NitricAwsPulumiProvider) Topic(ctx *pulumi.Context, parent pulumi.Resource, name string, config *deploymentspb.Topic) error {
	var err error

	a.topics[name] = &topic{}

	opts := []pulumi.ResourceOption{pulumi.Parent(parent)}

	// create the SNS topic
	a.topics[name].sns, err = sns.NewTopic(ctx, name, &sns.TopicArgs{
		Tags: pulumi.ToStringMap(common.Tags(a.stackId, name, resources.Topic)),
	}, opts...)
	if err != nil {
		return err
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
	}, opts...)
	if err != nil {
		return errors.WithMessage(err, "topic delay controller role")
	}

	policy := a.topics[name].sns.Arn.ApplyT(func(arn string) (string, error) {
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
	}, opts...)
	if err != nil {
		return errors.WithMessage(err, "topic delay controller role policy")
	}

	sfnDef := a.topics[name].sns.Arn.ApplyT(func(arn string) (string, error) {
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
	a.topics[name].sfn, err = sfn.NewStateMachine(ctx, name, &sfn.StateMachineArgs{
		RoleArn: sfnRole.Arn,
		// Apply the same name as the topic to the state machine
		Tags:       pulumi.ToStringMap(common.Tags(a.stackId, name, resources.Topic)),
		Definition: sfnDef,
	}, opts...)
	if err != nil {
		return errors.WithMessage(err, "topic delay controller")
	}

	for _, sub := range config.Subscriptions {
		targetLambda, ok := a.lambdas[sub.GetService()]
		if !ok {
			return fmt.Errorf("unable to find lambda %s for subscription", sub.GetService())
		}

		createSubscription(ctx, parent, name, a.topics[name].sns, targetLambda)
	}

	return nil
}
