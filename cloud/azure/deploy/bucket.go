package deploy

import (
	"fmt"

	deploymentspb "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"
	storagepb "github.com/nitrictech/nitric/core/pkg/proto/storage/v1"
	"github.com/pulumi/pulumi-azure-native-sdk/storage"
	pulumiEventgrid "github.com/pulumi/pulumi-azure/sdk/v4/go/azure/eventgrid"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func eventTypeToStorageEventType(eventType *storagepb.BlobEventType) []string {
	switch *eventType {
	case storagepb.BlobEventType_Created:
		return []string{"Microsoft.Storage.BlobCreated"}
	case storagepb.BlobEventType_Deleted:
		return []string{"Microsoft.Storage.BlobDeleted"}
	default:
		return []string{}
	}
}

func (p *NitricAzurePulumiProvider) newAzureBucketNotification(ctx *pulumi.Context, parent pulumi.Resource, bucketName string, config *deploymentspb.BucketListener) error {

	target, ok := p.containerApps[config.GetExecutionUnit()]
	if !ok {
		return fmt.Errorf("")
	}

	bucket, ok := p.buckets[bucketName]
	if !ok {
		return fmt.Errorf("")
	}

	opts := []pulumi.ResourceOption{pulumi.Parent(parent), pulumi.DependsOn([]pulumi.Resource{target.App, bucket})}

	hostUrl, err := target.HostUrl()
	if err != nil {
		return err
	}

	_, err = pulumiEventgrid.NewEventSubscription(ctx, ResourceName(ctx, bucketName+target.Name, EventSubscriptionRT), &pulumiEventgrid.EventSubscriptionArgs{
		Scope: p.storageAccount.ID(),
		WebhookEndpoint: pulumiEventgrid.EventSubscriptionWebhookEndpointArgs{
			Url: pulumi.Sprintf("%s/x-nitric-notification/bucket/%s?token=%s", hostUrl, bucketName, target.EventToken),
			// TODO: Reduce event chattiness here and handle internally in the Azure AppService HTTP Gateway?
			MaxEventsPerBatch:         pulumi.Int(1),
			ActiveDirectoryAppIdOrUri: target.Sp.ClientID,
			ActiveDirectoryTenantId:   target.Sp.TenantID,
		},
		RetryPolicy: pulumiEventgrid.EventSubscriptionRetryPolicyArgs{
			MaxDeliveryAttempts: pulumi.Int(30),
			EventTimeToLive:     pulumi.Int(5),
		},
		IncludedEventTypes: pulumi.ToStringArray(eventTypeToStorageEventType(&config.Config.BlobEventType)),
		SubjectFilter: pulumiEventgrid.EventSubscriptionSubjectFilterArgs{
			SubjectBeginsWith: pulumi.Sprintf("/blobServices/default/containers/%s/blobs/%s", bucketName, config.Config.KeyPrefixFilter),
		},
	}, opts...)
	if err != nil {
		return err
	}

	return nil
}

func (p *NitricAzurePulumiProvider) Bucket(ctx *pulumi.Context, parent pulumi.Resource, name string, config *deploymentspb.Bucket) error {
	var err error
	opts := []pulumi.ResourceOption{pulumi.Parent(parent)}

	p.buckets[name], err = storage.NewBlobContainer(ctx, ResourceName(ctx, name, StorageContainerRT), &storage.BlobContainerArgs{
		ResourceGroupName: p.resourceGroup.Name,
		AccountName:       p.storageAccount.Name,
	}, opts...)
	if err != nil {
		return err
	}

	for _, sub := range config.Listeners {
		err = p.newAzureBucketNotification(ctx, parent, name+sub.GetExecutionUnit(), sub)
		if err != nil {
			return err
		}
	}

	return nil
}
