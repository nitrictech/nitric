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
	"strings"

	deploymentspb "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"
	storagepb "github.com/nitrictech/nitric/core/pkg/proto/storage/v1"
	eventgrid "github.com/pulumi/pulumi-azure-native-sdk/eventgrid/v2"
	"github.com/pulumi/pulumi-azure-native-sdk/storage"
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

// removeWildcard - Remove the trailing wildcard from a prefix filter, they're not supported by Azure
func removeWildcard(prefixFilter string) string {
	return strings.TrimRight(prefixFilter, "*")
}

func (p *NitricAzurePulumiProvider) newAzureBucketNotification(ctx *pulumi.Context, parent pulumi.Resource, bucketName string, config *deploymentspb.BucketListener) error {
	target, ok := p.ContainerApps[config.GetService()]
	if !ok {
		return fmt.Errorf("target container app %s not found", config.GetService())
	}

	bucket, ok := p.Buckets[bucketName]
	if !ok {
		return fmt.Errorf("target bucket %s not found", bucketName)
	}

	opts := []pulumi.ResourceOption{pulumi.Parent(parent), pulumi.DependsOn([]pulumi.Resource{target.App, bucket})}

	hostUrl, err := target.HostUrl()
	if err != nil {
		return fmt.Errorf("unable to determine container app host URL: %w", err)
	}

	_, err = eventgrid.NewEventSubscription(ctx, ResourceName(ctx, bucketName+target.Name, EventSubscriptionRT), &eventgrid.EventSubscriptionArgs{
		Scope: p.StorageAccount.ID(),
		Destination: &eventgrid.WebHookEventSubscriptionDestinationArgs{
			EndpointType:                           pulumi.String("WebHook"),
			EndpointUrl:                            pulumi.Sprintf("%s/%s/x-nitric-notification/bucket/%s", hostUrl, target.EventToken, bucketName),
			MaxEventsPerBatch:                      pulumi.Int(1),
			AzureActiveDirectoryApplicationIdOrUri: target.Sp.ClientID,
			AzureActiveDirectoryTenantId:           target.Sp.TenantID,
		},
		RetryPolicy: eventgrid.RetryPolicyArgs{
			MaxDeliveryAttempts:      pulumi.Int(30),
			EventTimeToLiveInMinutes: pulumi.Int(5),
		},
		Filter: eventgrid.EventSubscriptionFilterArgs{
			SubjectBeginsWith:  pulumi.Sprintf("/blobServices/default/containers/%s/blobs/%s", bucketName, removeWildcard(config.Config.KeyPrefixFilter)),
			IncludedEventTypes: pulumi.ToStringArray(eventTypeToStorageEventType(&config.Config.BlobEventType)),
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

	p.Buckets[name], err = storage.NewBlobContainer(ctx, ResourceName(ctx, name, StorageContainerRT), &storage.BlobContainerArgs{
		ResourceGroupName: p.ResourceGroup.Name,
		AccountName:       p.StorageAccount.Name,
	}, opts...)
	if err != nil {
		return err
	}

	for _, sub := range config.Listeners {
		err = p.newAzureBucketNotification(ctx, parent, name, sub)
		if err != nil {
			return err
		}
	}

	return nil
}
