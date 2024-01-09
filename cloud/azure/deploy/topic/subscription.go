package topic

import (
	"github.com/nitrictech/nitric/cloud/azure/deploy/exec"
	"github.com/nitrictech/nitric/cloud/azure/deploy/utils"
	pulumiEventgrid "github.com/pulumi/pulumi-azure/sdk/v4/go/azure/eventgrid"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type AzureEventGridTopicSubscription struct {
	pulumi.ResourceState

	TopicResourceName string
	Subscription      *pulumiEventgrid.EventSubscription
}

type AzureEventGridTopicSubscriptionArgs struct {
	// The topic we want to source events from
	Topic *AzureEventGridTopic
	// The target we want to send events to
	Target *exec.ContainerApp
}

func (a *AzureEventGridTopic) AddSubscription(ctx *pulumi.Context, name string, args *AzureEventGridTopicSubscriptionArgs, opts ...pulumi.ResourceOption) error {
	parent := a

	hostUrl, err := args.Target.HostUrl()
	if err != nil {
		return err
	}

	_, err = pulumiEventgrid.NewEventSubscription(ctx, utils.ResourceName(ctx, name, utils.EventSubscriptionRT), &pulumiEventgrid.EventSubscriptionArgs{
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
	}, pulumi.Parent(parent))

	return err
}
