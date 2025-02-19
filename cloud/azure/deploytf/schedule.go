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
	"github.com/aws/jsii-runtime-go"
	"github.com/hashicorp/terraform-cdk-go/cdktf"
	"github.com/nitrictech/nitric/cloud/azure/deploy"
	"github.com/nitrictech/nitric/cloud/azure/deploytf/generated/schedule"
	deploymentspb "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"
)

// // Schedule - Deploy a Schedule
func (a *NitricAzureTerraformProvider) Schedule(stack cdktf.TerraformStack, name string, config *deploymentspb.Schedule) error {
	targetService := a.Services[config.GetTarget().GetService()]

	//	If this instance contains a schedule set the minimum instances to 1
	// schedules rely on the Dapr Runtime to trigger the function, without a running instance the Dapr Runtime will not execute, so the schedule won't trigger.
	if targetService.MinReplicas() == nil || *targetService.MinReplicas() < 1 {
		targetService.SetMinReplicas(jsii.Number(1))
	}

	cronExpression, err := deploy.GenerateCronExpression(config)
	if err != nil {
		return err
	}

	schedule.NewSchedule(stack, jsii.String(name), &schedule.ScheduleConfig{
		Name:                      jsii.String(name),
		ContainerAppEnvironmentId: a.Stack.ContainerAppEnvironmentIdOutput(),
		TargetEventToken:          targetService.EventTokenOutput(),
		TargetAppId:               targetService.DaprAppIdOutput(),
		CronExpression:            jsii.String(cronExpression),
	})

	return nil
}
