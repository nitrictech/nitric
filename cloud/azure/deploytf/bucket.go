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
	"github.com/nitrictech/nitric/cloud/azure/deploytf/generated/bucket"
	"github.com/nitrictech/nitric/cloud/common/deploy/resources"
	deploymentspb "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"
	storagepb "github.com/nitrictech/nitric/core/pkg/proto/storage/v1"
)

type BucketSubscriber struct {
	EventGridSubscriber `json:",inline"`
	// Url                       *string   `json:"url"`
	// ActiveDirectoryAppIdOrUri *string   `json:"active_directory_app_id_or_uri"`
	// ActiveDirectoryTenantId   *string   `json:"active_directory_tenant_id"`
	// EventToken                *string   `json:"event_token"`
	EventType []*string `json:"event_type"`
}

func eventTypeToStorageEventType(eventType storagepb.BlobEventType) []*string {
	switch eventType {
	case storagepb.BlobEventType_Created:
		return []*string{jsii.String("Microsoft.Storage.BlobCreated")}
	case storagepb.BlobEventType_Deleted:
		return []*string{jsii.String("Microsoft.Storage.BlobDeleted")}
	default:
		return []*string{}
	}
}

// Bucket - Deploy a Storage Bucket
func (n *NitricAzureTerraformProvider) Bucket(stack cdktf.TerraformStack, name string, config *deploymentspb.Bucket) error {
	listeners := map[string]BucketSubscriber{}

	allDependants := []cdktf.ITerraformDependable{}
	for _, v := range config.GetListeners() {
		svc := n.Services[v.GetService()]

		listeners[v.GetService()] = BucketSubscriber{
			EventGridSubscriber: EventGridSubscriber{
				Url:                       svc.EndpointOutput(),
				ActiveDirectoryAppIdOrUri: svc.ClientIdOutput(),
				ActiveDirectoryTenantId:   svc.TenantIdOutput(),
				EventToken:                svc.EventTokenOutput(),
			},
			EventType: eventTypeToStorageEventType(v.GetConfig().BlobEventType),
		}
		allDependants = append(allDependants, svc)
	}

	n.Buckets[name] = bucket.NewBucket(stack, jsii.String(name), &bucket.BucketConfig{
		Name:             jsii.String(name),
		StorageAccountId: n.Stack.StorageAccountIdOutput(),
		Listeners:        listeners,
		DependsOn:        &allDependants,
		Tags:             n.GetTags(*n.Stack.StackIdOutput(), name, resources.Bucket),
	})

	return nil
}
