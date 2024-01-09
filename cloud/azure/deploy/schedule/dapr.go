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
	"strings"

	"github.com/nitrictech/nitric/cloud/azure/deploy/exec"
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-azure-native-sdk/app"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type ScheduleArgs struct {
	ResourceGroupName pulumi.StringInput
	Target            *exec.ContainerApp
	Environment       *app.ManagedEnvironment
	Cron              string
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

	res.Component, err = app.NewDaprComponent(ctx, normalizedName, &app.DaprComponentArgs{
		ResourceGroupName: args.ResourceGroupName,
		EnvironmentName:   args.Environment.Name,
		ComponentName:     pulumi.String(normalizedName),
		ComponentType:     pulumi.String("bindings.cron"),
		Version:           pulumi.String("v1"),
		Metadata: app.DaprMetadataArray{
			app.DaprMetadataArgs{
				Name:  pulumi.String("schedule"),
				Value: pulumi.String(args.Cron),
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
