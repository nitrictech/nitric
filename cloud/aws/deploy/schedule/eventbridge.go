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
	_ "embed"
	"fmt"

	"github.com/nitrictech/nitric/cloud/aws/deploy/exec"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/iam"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/scheduler"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type AwsEventbridgeSchedule struct {
	pulumi.ResourceState
	Name     string
	Schedule *scheduler.Schedule
}

type AwsEventbridgeScheduleArgs struct {
	StackID pulumi.StringInput
	Exec    *exec.LambdaExecUnit
	Cron    string
	Tz      string
}

//go:embed scheduler-execution-permissions.json
var permissionsTemplate string

//go:embed scheduler-execution-role.json
var roleTemplate string

//go:embed scheduler-input.json
var scheduleInputTemplate string

func NewAwsEventbridgeSchedule(ctx *pulumi.Context, name string, args *AwsEventbridgeScheduleArgs, opts ...pulumi.ResourceOption) (*AwsEventbridgeSchedule, error) {
	res := &AwsEventbridgeSchedule{Name: name}

	err := ctx.RegisterComponentResource("nitric:schedule:AwsCloudwatchSchedule", name, res, opts...)
	if err != nil {
		return nil, err
	}

	awsCronValue, err := ConvertToAWS(args.Cron)
	if err != nil {
		return nil, err
	}

	// create a new role
	role, err := iam.NewRole(ctx, fmt.Sprintf("schedule-%s-role", name), &iam.RoleArgs{
		AssumeRolePolicy: pulumi.String(roleTemplate),
	})
	if err != nil {
		return nil, err
	}

	_, err = iam.NewRolePolicy(ctx, fmt.Sprintf("schedule-%s-policy", name), &iam.RolePolicyArgs{
		Policy: pulumi.Sprintf(permissionsTemplate, args.Exec.Function.Arn),
		Role:   role,
	})
	if err != nil {
		return nil, err
	}

	// Create a new eventbridge schedule
	res.Schedule, err = scheduler.NewSchedule(ctx, name, &scheduler.ScheduleArgs{
		ScheduleExpression:         pulumi.String(awsCronValue),
		ScheduleExpressionTimezone: pulumi.String(args.Tz),
		FlexibleTimeWindow: &scheduler.ScheduleFlexibleTimeWindowArgs{
			Mode: pulumi.String("OFF"),
		},
		Target: &scheduler.ScheduleTargetArgs{
			Arn:     args.Exec.Function.Arn,
			RoleArn: role.Arn,
			Input:   pulumi.Sprintf(scheduleInputTemplate, name),
		},
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}
