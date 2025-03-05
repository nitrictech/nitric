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
	"github.com/nitrictech/nitric/cloud/azure/deploytf/generated/keyvalue"
	deploymentspb "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"
)

// // KeyValueStore - Deploy a Key Value tioStore
func (a *NitricAzureTerraformProvider) KeyValueStore(stack cdktf.TerraformStack, name string, config *deploymentspb.KeyValueStore) error {
	a.KvStores[name] = keyvalue.NewKeyvalue(stack, jsii.String(name), &keyvalue.KeyvalueConfig{
		Name:               jsii.String(name),
		StorageAccountName: jsii.String(*a.Stack.StorageAccountNameOutput()),
	})

	return nil
}
