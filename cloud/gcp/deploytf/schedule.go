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

package deploytf

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/aws/jsii-runtime-go"
	"github.com/hashicorp/terraform-cdk-go/cdktf"
	"github.com/nitrictech/nitric/cloud/gcp/deploytf/generated/schedule"
	deploymentspb "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"
)

// // Schedule - Deploy a Schedule
func (a *NitricGcpTerraformProvider) Schedule(stack cdktf.TerraformStack, name string, config *deploymentspb.Schedule) error {
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
			parts[0] = fmt.Sprintf("%d", initialRate*24)
			parts[1] = "hours"
		}

		cronExpression = fmt.Sprintf("every %s %s", parts[0], parts[1])
	default:
		return fmt.Errorf("unknown schedule type, must be one of: cron, every")
	}

	svc, ok := a.Services[config.Target.GetService()]
	if !ok {
		return fmt.Errorf("service not found: %s", config.Target.GetService())
	}

	a.Schedules[name] = schedule.NewSchedule(stack, jsii.Sprintf("schedule_%s", name), &schedule.ScheduleConfig{
		ScheduleName:              jsii.String(name),
		ScheduleExpression:        jsii.String(cronExpression),
		ScheduleTimezone:          jsii.String(a.GcpConfig.ScheduleTimezone),
		TargetServiceUrl:          svc.ServiceEndpointOutput(),
		TargetServiceInvokerEmail: svc.InvokerServiceAccountEmailOutput(),
		ServiceToken:              svc.EventTokenOutput(),
	})

	return nil
}
