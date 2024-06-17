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
			"s3:ObjectCreated:*",
		}
	case storagepb.BlobEventType_Deleted:
		return []string{
			"s3:ObjectRemoved:*",
		}
	default:
		return []string{}
	}
}

// Bucket - Deploy a Storage Bucket
func (n *NitricGcpTerraformProvider) Bucket(stack cdktf.TerraformStack, name string, config *deploymentspb.Bucket) error {
	notificationTargets := map[string]interface{}{}

	for _, target := range config.Listeners {
		notificationTargets[target.GetService()] = map[string]interface{}{
			"url":    n.Services[target.GetService()].LambdaArnOutput(),
			"prefix": jsii.String(target.Config.KeyPrefixFilter),
			"events": jsii.Strings(eventsForBlobEventType(target.Config.BlobEventType)...),
		}
	}

	n.Buckets[name] = bucket.NewBucket(stack, jsii.Sprintf("bucket_%s", name), &bucket.BucketConfig{
		BucketName:          &name,
		StackId:             n.Stack.StackIdOutput(),
		NotificationTargets: &notificationTargets,
	})

	return nil
}
