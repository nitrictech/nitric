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
	"github.com/nitrictech/nitric/cloud/gcp/deploytf/generated/topic"
	deploymentspb "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"
)

type SubscriberService struct {
	Name                string
	Url                 string
	ServiceAccountEmail string
	EventToken          string
}

func (a *NitricGcpTerraformProvider) Topic(stack cdktf.TerraformStack, name string, config *deploymentspb.Topic) error {
	serviceEndpoints := map[string]*string{}

	for _, subscriber := range config.Subscriptions {
		subscriberSvc := a.Services[subscriber.GetService()]
		serviceEndpoints[subscriber.GetService()] = subscriberSvc.ServiceEndpointOutput()
	}

	subscriberInput := map[string]*SubscriberService{}

	for _, subscriber := range config.Subscriptions {
		subscriberSvc := a.Services[subscriber.GetService()]

		subscriberInput[subscriber.GetService()] = &SubscriberService{
			Name:                subscriber.GetService(),
			Url:                 *serviceEndpoints[subscriber.GetService()],
			ServiceAccountEmail: *subscriberSvc.ServiceAccountEmailOutput(),
			EventToken:          *subscriberSvc.EventTokenOutput(),
		}
	}

	a.Topics[name] = topic.NewTopic(stack, jsii.Sprintf("topic_%s", name), &topic.TopicConfig{
		StackId:            a.Stack.StackIdOutput(),
		TopicName:          jsii.String(name),
		SubscriberServices: subscriberInput,
	})

	return nil
}
