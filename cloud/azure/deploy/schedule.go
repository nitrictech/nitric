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

package deploy

import (
	"fmt"
	"strconv"
	"strings"

	deploymentspb "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-azure-native-sdk/app"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type ScheduleArgs struct {
	ResourceGroupName pulumi.StringInput
	Target            *ContainerApp
	Environment       *app.ManagedEnvironment
	Schedule          *deploymentspb.Schedule
}

type Schedule struct {
	pulumi.ResourceState
	Name      string
	Component *app.DaprComponent
}

func (p *NitricAzurePulumiProvider) Schedule(ctx *pulumi.Context, parent pulumi.Resource, name string, config *deploymentspb.Schedule) error {
	var err error
	opts := []pulumi.ResourceOption{pulumi.Parent(parent)}

	target := p.containerApps[config.GetTarget().GetService()]

	normalizedName := strings.ToLower(strings.ReplaceAll(name, " ", "-"))
	cronExpression := ""

	switch t := config.Cadence.(type) {
	case *deploymentspb.Schedule_Cron:
		cronExpression = t.Cron.Expression
	case *deploymentspb.Schedule_Every:
		parts := strings.Split(strings.TrimSpace(t.Every.Rate), " ")
		if len(parts) != 2 {
			return fmt.Errorf("invalid schedule rate: %s", t.Every.Rate)
		}

		initialRate, err := strconv.Atoi(parts[0])
		if err != nil {
			return fmt.Errorf("invalid schedule rate, must start with an integer")
		}

		// Dapr cron bindings only support hours, minutes and seconds. Convert days to hours
		if strings.HasPrefix(parts[1], "day") {
			parts[0] = fmt.Sprintf("%d", initialRate*24)
			parts[1] = "hours"
		}

		cronExpression = fmt.Sprintf("@every %s%c", parts[0], parts[1][0])
	default:
		return fmt.Errorf("unknown schedule type, must be one of: cron, every")
	}

	_, err = app.NewDaprComponent(ctx, normalizedName, &app.DaprComponentArgs{
		ResourceGroupName: p.resourceGroup.Name,
		EnvironmentName:   p.containerEnv.ManagedEnv.Name,
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
				Value: pulumi.Sprintf("/x-nitric-schedule/%s?token=%s", normalizedName, target.EventToken),
			},
		},
		Scopes: pulumi.StringArray{
			// Limit the scope to the target container app
			target.App.Configuration.Dapr().AppId().Elem(),
		},
	}, opts...)

	if err != nil {
		return errors.WithMessage(err, fmt.Sprintf("unable to create nitric schedule %s: failed to create DaprComponent for app", name))
	}

	return nil
}
