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
	"strings"

	"github.com/aws/jsii-runtime-go"
	"github.com/hashicorp/terraform-cdk-go/cdktf"
	"github.com/nitrictech/nitric/cloud/azure/deploytf/generated/topic"
	deploymentspb "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"
)

type WebhookSubscriber struct {
	ClientId            *string `json:"client_id"`
	ClientSecret        *string `json:"client_secret"`
	TenantId            *string `json:"tenant_id"`
	EventGridSubscriber `json:",inline"`
}

func (a *NitricAzureTerraformProvider) Topic(stack cdktf.TerraformStack, name string, config *deploymentspb.Topic) error {
	listeners := map[string]WebhookSubscriber{}

	allDependants := []cdktf.ITerraformDependable{}
	for _, v := range config.GetSubscriptions() {
		svc := a.Services[v.GetService()]

		normalizedServiceName := strings.Replace(v.GetService(), "_", "-", -1)

		listeners[normalizedServiceName] = WebhookSubscriber{
			ClientId:     svc.ClientIdOutput(),
			ClientSecret: svc.ClientSecretOutput(),
			TenantId:     svc.TenantIdOutput(),
			EventGridSubscriber: EventGridSubscriber{
				Url:                       svc.EndpointOutput(),
				ActiveDirectoryAppIdOrUri: svc.ClientIdOutput(),
				ActiveDirectoryTenantId:   svc.TenantIdOutput(),
				EventToken:                svc.EventTokenOutput(),
			},
		}

		allDependants = append(allDependants, svc)
	}

	a.Topics[name] = topic.NewTopic(stack, jsii.String(name), &topic.TopicConfig{
		Name:              jsii.String(name),
		StackName:         a.Stack.StackNameOutput(),
		ResourceGroupName: a.Stack.ResourceGroupNameOutput(),
		Location:          jsii.String(a.Region),
		Listeners:         listeners,
		DependsOn:         &allDependants,
	})

	return nil
}
