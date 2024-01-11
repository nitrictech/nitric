// Copyright 2021 Nitric Technologies Pty Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
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
	"fmt"
	"strconv"
	"strings"

	"github.com/nitrictech/nitric/cloud/azure/deploy/exec"
	deploymentspb "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-azure-native-sdk/app"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type ScheduleArgs struct {
	ResourceGroupName pulumi.StringInput
	Target            *exec.ContainerApp
	Environment       *app.ManagedEnvironment
	Schedule          *deploymentspb.Schedule
}

type Schedule struct {
	pulumi.ResourceState
	Name      string
	Component *app.DaprComponent
}

func NewDaprCronBindingSchedule(ctx *pulumi.Context, name string, args *ScheduleArgs, opts ...pulumi.ResourceOption) (*Schedule, error) {
	res := &Schedule{
		Name: name,
	}
	normalizedName := strings.ToLower(strings.ReplaceAll(name, " ", "-"))
	err := ctx.RegisterComponentResource("nitric:func:ContainerApp", name, res, opts...)
	if err != nil {
		return nil, err
	}

	cronExpression := ""

	switch t := args.Schedule.Cadence.(type) {
	case *deploymentspb.Schedule_Cron:
		cronExpression = t.Cron.Expression
	case *deploymentspb.Schedule_Every:
		parts := strings.Split(strings.TrimSpace(t.Every.Rate), " ")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid schedule rate: %s", t.Every.Rate)
		}

		initialRate, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, fmt.Errorf("invalid schedule rate, must start with an integer")
		}

		// Dapr cron bindings only support hours, minutes and seconds. Convert days to hours
		if strings.HasPrefix(parts[1], "day") {
			parts[0] = fmt.Sprintf("%d", initialRate*24)
			parts[1] = "hours"
		}

		cronExpression = fmt.Sprintf("@every %s%c", parts[0], parts[1][0])
	default:
		return nil, fmt.Errorf("unknown schedule type, must be one of: cron, every")
	}

	res.Component, err = app.NewDaprComponent(ctx, normalizedName, &app.DaprComponentArgs{
		ResourceGroupName: args.ResourceGroupName,
		EnvironmentName:   args.Environment.Name,
		ComponentName:     pulumi.String(normalizedName),
		ComponentType:     pulumi.String("bindings.cron"),
		Version:           pulumi.String("v1"),
		Metadata: app.DaprMetadataArray{
			app.DaprMetadataArgs{
				Name:  pulumi.String("schedule"),
				Value: pulumi.String(cronExpression),
			},
			app.DaprMetadataArgs{
				Name:  pulumi.String("route"),
				Value: pulumi.Sprintf("/x-nitric-schedule/%s?token=%s", normalizedName, args.Target.EventToken),
			},
		},
		Scopes: pulumi.StringArray{
			// Limit the scope to the target container app
			args.Target.App.Configuration.Dapr().AppId().Elem(),
		},
	})

	if err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("unable to create nitric schedule %s: failed to create DaprComponent for app", name))
	}

	return res, nil
}
