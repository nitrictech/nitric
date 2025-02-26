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
	deploymentspb "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"
	"github.com/pulumi/pulumi-azure-native-sdk/storage"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func (p *NitricAzurePulumiProvider) KeyValueStore(ctx *pulumi.Context, parent pulumi.Resource, name string, config *deploymentspb.KeyValueStore) error {
	var err error
	opts := []pulumi.ResourceOption{pulumi.Parent(parent)}

	p.KeyValueStores[name], err = storage.NewTable(ctx, name, &storage.TableArgs{
		AccountName:       p.StorageAccount.Name,
		ResourceGroupName: p.ResourceGroup.Name,
		TableName:         pulumi.String(name),
	}, opts...)

	return err
}
