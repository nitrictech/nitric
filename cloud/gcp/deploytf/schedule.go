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

	"github.com/aws/jsii-runtime-go"
	"github.com/hashicorp/terraform-cdk-go/cdktf"
	"github.com/nitrictech/nitric/cloud/aws/deploy"
	"github.com/nitrictech/nitric/cloud/aws/deploytf/generated/schedule"
	deploymentspb "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"
)

// // Schedule - Deploy a Schedule
func (a *NitricGcpTerraformProvider) Schedule(stack cdktf.TerraformStack, name string, config *deploymentspb.Schedule) error {
	var err error

	awsScheduleExpression := ""
	switch config.Cadence.(type) {
	case *deploymentspb.Schedule_Cron:
		// handle cron
		awsScheduleExpression, err = deploy.ConvertToAWS(config.GetCron().Expression)
		if err != nil {
			return err
		}
	case *deploymentspb.Schedule_Every:
		// handle rate
		awsScheduleExpression = fmt.Sprintf("rate(%s)", config.GetEvery().Rate)
	default:
		return fmt.Errorf("unknown schedule type, must be one of: cron, every")
	}

	svc, ok := a.Services[config.Target.GetService()]
	if !ok {
		return fmt.Errorf("service not found: %s", config.Target.GetService())
	}

	a.Schedules[name] = schedule.NewSchedule(stack, jsii.Sprintf("schedule_%s", name), &schedule.ScheduleConfig{
		ScheduleName:       jsii.String(name),
		ScheduleExpression: jsii.String(awsScheduleExpression),
		ScheduleTimezone:   jsii.String(a.GcpConfig.ScheduleTimezone),
		TargetLambdaArn:    svc.LambdaArnOutput(),
		StackId:            a.Stack.StackIdOutput(),
	})

	return nil
}
