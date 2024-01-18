// Copyright 2021 Nitric Technologies Pty Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package deploy

import (
	"fmt"

	"github.com/nitrictech/nitric/cloud/azure/deploy/utils"
	nitricresources "github.com/nitrictech/nitric/cloud/common/deploy/resources"
	"github.com/nitrictech/nitric/cloud/common/deploy/tags"
	deploymentspb "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"
	"github.com/pulumi/pulumi-azure-native-sdk/eventgrid"
	"github.com/pulumi/pulumi-azure-native-sdk/resources"
	pulumiEventgrid "github.com/pulumi/pulumi-azure/sdk/v4/go/azure/eventgrid"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Topics
type AzureEventGridTopic struct {
	pulumi.ResourceState

	Name               string
	SourceResourceName string
	Topic              *eventgrid.Topic
}

type AzureEventGridTopicArgs struct {
	StackID       string
	ResourceGroup *resources.ResourceGroup
}

func (p *NitricAzurePulumiProvider) newEventGridTopicSubscription(ctx *pulumi.Context, parent pulumi.Resource, topicName string, config *deploymentspb.SubscriptionTarget) error {
	opts := []pulumi.ResourceOption{pulumi.Parent(parent)}

	target, ok := p.containerApps[config.GetExecutionUnit()]
	if !ok {
		return fmt.Errorf("Unable to find container app for execution unit: %s", config.GetExecutionUnit())
	}

	topic, ok := p.topics[topicName]
	if !ok {
		return fmt.Errorf("Unable to find topic: %s", topicName)
	}

	hostUrl, err := target.HostUrl()
	if err != nil {
		return err
	}

	subName := topicName + "-" + config.GetExecutionUnit()

	_, err = pulumiEventgrid.NewEventSubscription(ctx, utils.ResourceName(ctx, subName, utils.EventSubscriptionRT), &pulumiEventgrid.EventSubscriptionArgs{
		Scope: topic.ID(),
		WebhookEndpoint: pulumiEventgrid.EventSubscriptionWebhookEndpointArgs{
			Url: pulumi.Sprintf("%s/x-nitric-topic/%s?token=%s", hostUrl, topicName, target.EventToken),
			// TODO: Reduce event chattiness here and handle internally in the Azure AppService HTTP Gateway?
			MaxEventsPerBatch:         pulumi.Int(1),
			ActiveDirectoryAppIdOrUri: target.Sp.ClientID,
			ActiveDirectoryTenantId:   target.Sp.TenantID,
		},
		RetryPolicy: pulumiEventgrid.EventSubscriptionRetryPolicyArgs{
			MaxDeliveryAttempts: pulumi.Int(30),
			EventTimeToLive:     pulumi.Int(5),
		},
	}, opts...)

	return err
}

func (p *NitricAzurePulumiProvider) Topic(ctx *pulumi.Context, parent pulumi.Resource, name string, config *deploymentspb.Topic) error {
	var err error
	opts := []pulumi.ResourceOption{pulumi.Parent(parent)}

	p.topics[name], err = eventgrid.NewTopic(ctx, utils.ResourceName(ctx, name, utils.EventGridRT), &eventgrid.TopicArgs{
		ResourceGroupName: p.resourceGroup.Name,
		Location:          p.resourceGroup.Location,
		Tags:              pulumi.ToStringMap(tags.Tags(p.stackId, name, nitricresources.Topic)),
	}, opts...)
	if err != nil {
		return err
	}

	for _, sub := range config.Subscriptions {
		err = p.newEventGridTopicSubscription(ctx, parent, name, sub)
		if err != nil {
			return err
		}
	}

	return nil
}
