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

	nitricresources "github.com/nitrictech/nitric/cloud/common/deploy/resources"
	"github.com/nitrictech/nitric/cloud/common/deploy/tags"
	deploymentspb "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"
	eventgrid "github.com/pulumi/pulumi-azure-native-sdk/eventgrid/v2"
	"github.com/pulumi/pulumi-azure-native-sdk/resources"
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

	target, ok := p.ContainerApps[config.GetService()]
	if !ok {
		return fmt.Errorf("Unable to find container app for service: %s", config.GetService())
	}

	topic, ok := p.Topics[topicName]
	if !ok {
		return fmt.Errorf("Unable to find topic: %s", topicName)
	}

	hostUrl, err := target.HostUrl()
	if err != nil {
		return err
	}

	subName := topicName + "-" + config.GetService()

	_, err = eventgrid.NewEventSubscription(ctx, ResourceName(ctx, subName, EventSubscriptionRT), &eventgrid.EventSubscriptionArgs{
		Scope: topic.ID(),
		Destination: &eventgrid.WebHookEventSubscriptionDestinationArgs{
			EndpointType:                           pulumi.String("WebHook"),
			EndpointUrl:                            pulumi.Sprintf("%s/%s/x-nitric-topic/%s", hostUrl, target.EventToken, topicName),
			MaxEventsPerBatch:                      pulumi.Int(1),
			AzureActiveDirectoryApplicationIdOrUri: target.Sp.ClientID,
			AzureActiveDirectoryTenantId:           target.Sp.TenantID,
		},
	}, opts...)

	return err
}

func (p *NitricAzurePulumiProvider) Topic(ctx *pulumi.Context, parent pulumi.Resource, name string, config *deploymentspb.Topic) error {
	var err error
	opts := []pulumi.ResourceOption{pulumi.Parent(parent)}

	p.Topics[name], err = eventgrid.NewTopic(ctx, ResourceName(ctx, name, EventGridRT), &eventgrid.TopicArgs{
		ResourceGroupName: p.ResourceGroup.Name,
		Location:          p.ResourceGroup.Location,
		Tags:              pulumi.ToStringMap(tags.Tags(p.StackId, name, nitricresources.Topic)),
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
