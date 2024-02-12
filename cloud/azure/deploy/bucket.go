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
	target, ok := p.containerApps[config.GetService()]
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
			Url: pulumi.Sprintf("%s/%s/x-nitric-notification/bucket/%s", hostUrl, target.EventToken, bucketName),
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
		err = p.newAzureBucketNotification(ctx, parent, name+sub.GetService(), sub)
		if err != nil {
			return err
		}
	}

	return nil
}
