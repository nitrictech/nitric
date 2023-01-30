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

package schedule

import (
	"github.com/nitrictech/nitric/cloud/aws/deploy/exec"
	common "github.com/nitrictech/nitric/cloud/common/deploy/tags"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/cloudwatch"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/lambda"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type AwsCloudwatchSchedule struct {
	pulumi.ResourceState
	Name        string
	EventRule   *cloudwatch.EventRule
	EventTarget *cloudwatch.EventTarget
	Permission  *lambda.Permission
}

type AwsCloudwatchScheduleArgs struct {
	StackID pulumi.StringInput
	Exec    *exec.LambdaExecUnit
	Cron    string
}

func NewAwsCloudwatchSchedule(ctx *pulumi.Context, name string, args *AwsCloudwatchScheduleArgs, opts ...pulumi.ResourceOption) (*AwsCloudwatchSchedule, error) {
	res := &AwsCloudwatchSchedule{Name: name}

	err := ctx.RegisterComponentResource("nitric:schedule:AwsCloudwatchSchedule", name, res, opts...)
	if err != nil {
		return nil, err
	}

	awsCronValue, err := ConvertToAWS(args.Cron)
	if err != nil {
		return nil, err
	}

	res.EventRule, err = cloudwatch.NewEventRule(ctx, name, &cloudwatch.EventRuleArgs{
		ScheduleExpression: pulumi.String(awsCronValue),
		Tags:               common.Tags(ctx, args.StackID, name),
	}, opts...)
	if err != nil {
		return nil, err
	}

	// Give the event rule permission to execute the lambda
	res.Permission, err = lambda.NewPermission(ctx, res.Name, &lambda.PermissionArgs{
		Function:  args.Exec.Function.Name,
		Action:    pulumi.String("lambda:InvokeFunction"),
		Principal: pulumi.String("events.amazonaws.com"),
		SourceArn: res.EventRule.Arn,
	})
	if err != nil {
		return nil, err
	}

	res.EventTarget, err = cloudwatch.NewEventTarget(ctx, name+"Target", &cloudwatch.EventTargetArgs{
		Rule: res.EventRule.Name,
		Arn:  args.Exec.Function.Arn,
	}, opts...)
	if err != nil {
		return nil, err
	}

	return res, nil
}
