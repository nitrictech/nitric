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

package topic

import (
	"github.com/nitrictech/nitric/cloud/azure/deploy/exec"
	"github.com/nitrictech/nitric/cloud/azure/deploy/utils"
	common "github.com/nitrictech/nitric/cloud/common/deploy/tags"
	"github.com/pulumi/pulumi-azure-native-sdk/eventgrid"
	"github.com/pulumi/pulumi-azure-native-sdk/resources"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"

	pulumiEventgrid "github.com/pulumi/pulumi-azure/sdk/v4/go/azure/eventgrid"
)

// Topics
type AzureEventGridTopic struct {
	pulumi.ResourceState

	Name  string
	Topic *eventgrid.Topic
}

type AzureEventGridTopicArgs struct {
	StackID       string
	ResourceGroup *resources.ResourceGroup
}

func NewAzureEventGridTopic(ctx *pulumi.Context, name string, args *AzureEventGridTopicArgs, opts ...pulumi.ResourceOption) (*AzureEventGridTopic, error) {
	res := &AzureEventGridTopic{Name: name}

	err := ctx.RegisterComponentResource("nitric:topic:AzureEventGridTopic", name, res, opts...)
	if err != nil {
		return nil, err
	}

	res.Topic, err = eventgrid.NewTopic(ctx, utils.ResourceName(ctx, res.Name, utils.EventGridRT), &eventgrid.TopicArgs{
		ResourceGroupName: args.ResourceGroup.Name,
		Location:          args.ResourceGroup.Location,
		Tags:              pulumi.ToStringMap(common.Tags(ctx, args.StackID, res.Name)),
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

// Subscriptions
type AzureEventGridTopicSubscription struct {
	pulumi.ResourceState

	Name         string
	Subscription *pulumiEventgrid.EventSubscription
}

type AzureEventGridTopicSubscriptionArgs struct {
	// The topic we want to source events from
	Topic *AzureEventGridTopic
	// The target we want to send events to
	Target *exec.ContainerApp
}

func NewAzureEventGridTopicSubscription(ctx *pulumi.Context, name string, args *AzureEventGridTopicSubscriptionArgs, opts ...pulumi.ResourceOption) (*AzureEventGridTopicSubscription, error) {
	res := &AzureEventGridTopicSubscription{Name: name}

	err := ctx.RegisterComponentResource("nitric:topic:AzureEventGridTopicSubscription", name, res, opts...)
	if err != nil {
		return nil, err
	}

	hostUrl, err := args.Target.HostUrl()
	if err != nil {
		return nil, err
	}

	res.Subscription, err = pulumiEventgrid.NewEventSubscription(ctx, utils.ResourceName(ctx, name, utils.EventSubscriptionRT), &pulumiEventgrid.EventSubscriptionArgs{
		Scope: args.Topic.Topic.ID(),
		WebhookEndpoint: pulumiEventgrid.EventSubscriptionWebhookEndpointArgs{
			Url: pulumi.Sprintf("%s/x-nitric-topic/%s?token=%s", hostUrl, args.Topic.Name, args.Target.EventToken),
			// TODO: Reduce event chattiness here and handle internally in the Azure AppService HTTP Gateway?
			MaxEventsPerBatch:         pulumi.Int(1),
			ActiveDirectoryAppIdOrUri: args.Target.Sp.ClientID,
			ActiveDirectoryTenantId:   args.Target.Sp.TenantID,
		},
		RetryPolicy: pulumiEventgrid.EventSubscriptionRetryPolicyArgs{
			MaxDeliveryAttempts: pulumi.Int(30),
			EventTimeToLive:     pulumi.Int(5),
		},
	}, pulumi.Parent(res))
	if err != nil {
		return nil, err
	}

	return res, nil
}
