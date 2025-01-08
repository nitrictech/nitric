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
	"fmt"

	"github.com/aws/jsii-runtime-go"
	"github.com/hashicorp/terraform-cdk-go/cdktf"
	"github.com/nitrictech/nitric/cloud/azure/deploytf/generated/bucket"
	deploymentspb "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"
)

type BucketListener struct {
	Url                       *string `json:"url"`
	ActiveDirectoryAppIdOrUri *string `json:"active_directory_app_id_or_uri"`
	ActiveDirectoryTenantId   *string `json:"active_directory_tenant_id"`
}

// Bucket - Deploy a Storage Bucket
func (n *NitricAzureTerraformProvider) Bucket(stack cdktf.TerraformStack, name string, config *deploymentspb.Bucket) error {
	listeners := map[string]BucketListener{}

	for _, v := range config.GetListeners() {
		svc := n.Services[v.GetService()]

		listeners[v.GetService()] = BucketListener{
			Url:                       svc.EndpointOutput(),
			ActiveDirectoryAppIdOrUri: svc.ClientIdOutput(),
			ActiveDirectoryTenantId:   svc.TenantIdOutput(),
		}
	}

	bucket.NewBucket(stack, jsii.String(name), &bucket.BucketConfig{
		Name:      jsii.String(name),
		StackName: n.Stack.StackNameOutput(),
		Listeners: listeners,
	})

	return fmt.Errorf("not implemented")
}
