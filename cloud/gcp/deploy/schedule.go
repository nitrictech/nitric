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
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/nitrictech/nitric/cloud/gcp/deploy/exec"
	deploymentspb "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/cloudscheduler"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type CloudScheduler struct {
	pulumi.ResourceState

	Name string
	Job  *cloudscheduler.Job
}

type CloudSchedulerArgs struct {
	Location string
	Exec     *exec.CloudRunner
	Schedule *deploymentspb.Schedule
	Tz       string
}

type ScheduleEvent struct {
	PayloadType string                 `yaml:"payloadType"`
	Payload     map[string]interface{} `yaml:"payload,omitempty"`
}

func (p *NitricGcpPulumiProvider) Schedule(ctx *pulumi.Context, parent pulumi.Resource, name string, config *deploymentspb.Schedule) error {
	opts := append([]pulumi.ResourceOption{}, pulumi.Parent(parent))

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

		// Google App Engine cron syntax only support hours, minutes and seconds. Convert days to hours
		if strings.HasPrefix(parts[1], "day") {
			// TODO: verify that rates exceeding 24 hours are supported.
			parts[0] = fmt.Sprintf("%d", initialRate*24)
			parts[1] = "hours"
		}

		cronExpression = fmt.Sprintf("every %s %s", parts[0], parts[1])
	default:
		return fmt.Errorf("unknown schedule type, must be one of: cron, every")
	}

	eventJSON, err := json.Marshal(map[string]interface{}{
		"schedule": name,
	})
	if err != nil {
		return err
	}

	targetService := p.cloudRunServices[config.Target.GetExecutionUnit()]

	payload := base64.StdEncoding.EncodeToString(eventJSON)

	_, err = cloudscheduler.NewJob(ctx, name, &cloudscheduler.JobArgs{
		TimeZone: pulumi.String(p.config.ScheduleTimezone),
		HttpTarget: &cloudscheduler.JobHttpTargetArgs{
			Uri: pulumi.Sprintf("%s/x-nitric-schedule/%s?token=%s", targetService.Url, name, targetService.EventToken),
			OidcToken: &cloudscheduler.JobHttpTargetOidcTokenArgs{
				ServiceAccountEmail: targetService.Invoker.Email,
			},
			Body: pulumi.String(payload),
		},
		Schedule: pulumi.String(cronExpression),
	}, opts...)

	return err
}
