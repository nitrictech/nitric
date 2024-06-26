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
	"github.com/nitrictech/nitric/cloud/gcp/deploytf/generated/bucket"
	deploymentspb "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"
	storagepb "github.com/nitrictech/nitric/core/pkg/proto/storage/v1"
)

func eventsForBlobEventType(blobEventType storagepb.BlobEventType) []string {
	switch blobEventType {
	case storagepb.BlobEventType_Created:
		return []string{
			"OBJECT_FINALIZE",
		}
	case storagepb.BlobEventType_Deleted:
		return []string{
			"OBJECT_DELETE",
		}
	default:
		return []string{}
	}
}

type NotifiedService struct {
	// Explicit JSON names required for JSII serialization
	Name                       string   `json:"name"`
	Url                        string   `json:"url"`
	InvokerServiceAccountEmail string   `json:"invoker_service_account_email"`
	EventToken                 string   `json:"event_token"`
	Prefix                     string   `json:"prefix"`
	Events                     []string `json:"events"`
}

// Bucket - Deploy a Storage Bucket
func (n *NitricGcpTerraformProvider) Bucket(stack cdktf.TerraformStack, name string, config *deploymentspb.Bucket) error {
	notificationTargets := map[string]*NotifiedService{}

	for _, target := range config.Listeners {
		notificationTargets[target.GetService()] = &NotifiedService{
			Name:                       target.GetService(),
			Url:                        *n.Services[target.GetService()].ServiceEndpointOutput(),
			InvokerServiceAccountEmail: *n.Services[target.GetService()].InvokerServiceAccountEmailOutput(),
			EventToken:                 *n.Services[target.GetService()].EventTokenOutput(),
			Events:                     eventsForBlobEventType(target.Config.BlobEventType),
			Prefix:                     target.Config.KeyPrefixFilter,
		}
	}

	n.Buckets[name] = bucket.NewBucket(stack, jsii.Sprintf("bucket_%s", name), &bucket.BucketConfig{
		BucketName:          &name,
		StackId:             n.Stack.StackIdOutput(),
		NotificationTargets: &notificationTargets,
	})

	return nil
}
