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
	"github.com/nitrictech/nitric/cloud/azure/deploytf/generated/queue"
	"github.com/nitrictech/nitric/cloud/common/deploy/resources"
	deploymentspb "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"
)

// // Queue - Deploy a Queue
func (a *NitricAzureTerraformProvider) Queue(stack cdktf.TerraformStack, name string, config *deploymentspb.Queue) error {
	a.Queues[name] = queue.NewQueue(stack, jsii.String(name), &queue.QueueConfig{
		Name:               jsii.String(name),
		StorageAccountName: a.Stack.StorageAccountNameOutput(),
		Tags:               a.GetTags(*a.Stack.StackIdOutput(), name, resources.Queue),
	})
	return nil
}
