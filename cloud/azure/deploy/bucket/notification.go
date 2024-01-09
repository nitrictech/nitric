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

package bucket

import (
	"github.com/nitrictech/nitric/cloud/azure/deploy/exec"
	"github.com/nitrictech/nitric/cloud/azure/deploy/utils"
	v1 "github.com/nitrictech/nitric/core/pkg/proto/storage/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"

	"github.com/pulumi/pulumi-azure-native-sdk/storage"
	pulumiEventgrid "github.com/pulumi/pulumi-azure/sdk/v4/go/azure/eventgrid"
)

func eventTypeToStorageEventType(eventType *v1.BlobEventType) []string {
	switch *eventType {
	case v1.BlobEventType_Created:
		return []string{"Microsoft.Storage.BlobCreated"}
	case v1.BlobEventType_Deleted:
		return []string{"Microsoft.Storage.BlobDeleted"}
	default:
		return []string{}
	}
}

// Subscriptions
type AzureBucketNotification struct {
	pulumi.ResourceState

	Name         string
	Subscription *pulumiEventgrid.EventSubscription
}

type AzureBucketNotificationArgs struct {
	Bucket         *AzureStorageBucket
	StorageAccount *storage.StorageAccount
	Config         *v1.RegistrationRequest
	Target         *exec.ContainerApp
}

func NewAzureBucketNotification(ctx *pulumi.Context, name string, args *AzureBucketNotificationArgs, opts ...pulumi.ResourceOption) (*AzureBucketNotification, error) {
	res := &AzureBucketNotification{
		Name: name,
	}

	err := ctx.RegisterComponentResource("nitric:bucket:AzureBucketNotification", name, res, opts...)
	if err != nil {
		return nil, err
	}

	hostUrl, err := args.Target.HostUrl()
	if err != nil {
		return nil, err
	}

	res.Subscription, err = pulumiEventgrid.NewEventSubscription(ctx, utils.ResourceName(ctx, name, utils.EventSubscriptionRT), &pulumiEventgrid.EventSubscriptionArgs{
		Scope: args.StorageAccount.ID(),
		WebhookEndpoint: pulumiEventgrid.EventSubscriptionWebhookEndpointArgs{
			Url: pulumi.Sprintf("%s/x-nitric-notification/bucket/%s?token=%s", hostUrl, args.Bucket.Name, args.Target.EventToken),
			// TODO: Reduce event chattiness here and handle internally in the Azure AppService HTTP Gateway?
			MaxEventsPerBatch:         pulumi.Int(1),
			ActiveDirectoryAppIdOrUri: args.Target.Sp.ClientID,
			ActiveDirectoryTenantId:   args.Target.Sp.TenantID,
		},
		RetryPolicy: pulumiEventgrid.EventSubscriptionRetryPolicyArgs{
			MaxDeliveryAttempts: pulumi.Int(30),
			EventTimeToLive:     pulumi.Int(5),
		},
		IncludedEventTypes: pulumi.ToStringArray(eventTypeToStorageEventType(&args.Config.BlobEventType)),
		SubjectFilter: pulumiEventgrid.EventSubscriptionSubjectFilterArgs{
			SubjectBeginsWith: pulumi.Sprintf("/blobServices/default/containers/%s/blobs/%s", args.Bucket.Name, args.Config.KeyPrefixFilter),
		},
	}, pulumi.Parent(res), pulumi.DependsOn([]pulumi.Resource{args.Target.App, args.Bucket.Container}))
	if err != nil {
		return nil, err
	}

	return res, nil
}
