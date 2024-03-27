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
	"fmt"

	"github.com/nitrictech/nitric/cloud/aws/deploy/embeds"
	deploymentspb "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/iam"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/scheduler"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type AwsEventbridgeSchedule struct {
	pulumi.ResourceState
	Name     string
	Schedule *scheduler.Schedule
}

func (a *NitricAwsPulumiProvider) Schedule(ctx *pulumi.Context, parent pulumi.Resource, name string, config *deploymentspb.Schedule) error {
	var err error
	opts := []pulumi.ResourceOption{pulumi.Parent(parent)}

	awsScheduleExpression := ""
	switch config.Cadence.(type) {
	case *deploymentspb.Schedule_Cron:
		// handle cron
		awsScheduleExpression, err = ConvertToAWS(config.GetCron().Expression)
	case *deploymentspb.Schedule_Every:
		// handle rate
		awsScheduleExpression = fmt.Sprintf("rate(%s)", config.GetEvery().Rate)
	default:
		return fmt.Errorf("unknown schedule type, must be one of: cron, every")
	}

	if err != nil {
		return err
	}

	// create a new role
	role, err := iam.NewRole(ctx, fmt.Sprintf("schedule-%s-role", name), &iam.RoleArgs{
		AssumeRolePolicy: embeds.GetScheduleRoleDocument(),
	}, opts...)
	if err != nil {
		return err
	}

	target, ok := a.Lambdas[config.Target.GetService()]
	if !ok {
		return fmt.Errorf("unable to find target lambda: %s", config.Target.GetService())
	}

	_, err = iam.NewRolePolicy(ctx, fmt.Sprintf("schedule-%s-policy", name), &iam.RolePolicyArgs{
		Policy: embeds.GetSchedulePermissionDocument(target.Arn),
		Role:   role,
	}, opts...)
	if err != nil {
		return err
	}

	// Create a new eventbridge schedule
	_, err = scheduler.NewSchedule(ctx, name, &scheduler.ScheduleArgs{
		ScheduleExpression:         pulumi.String(awsScheduleExpression),
		ScheduleExpressionTimezone: pulumi.String(a.AwsConfig.ScheduleTimezone),
		FlexibleTimeWindow: &scheduler.ScheduleFlexibleTimeWindowArgs{
			Mode: pulumi.String("OFF"),
		},
		Target: &scheduler.ScheduleTargetArgs{
			Arn:     target.Arn,
			RoleArn: role.Arn,
			RetryPolicy: &scheduler.ScheduleTargetRetryPolicyArgs{
				MaximumEventAgeInSeconds: pulumi.Int(60),
				MaximumRetryAttempts:     pulumi.Int(5),
			},
			Input: embeds.GetScheduleInputDocument(pulumi.String(name)),
		},
	}, opts...)
	if err != nil {
		return err
	}

	return nil
}
