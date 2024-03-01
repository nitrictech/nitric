// Copyright Nitric Pty Ltd.
//
// SPDX-License-Identifier: Apache-2.0
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
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
	"github.com/nitrictech/nitric/cloud/common/deploy/resources"
	"github.com/nitrictech/nitric/cloud/common/deploy/tags"
	deploymentspb "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"
	v1 "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/dynamodb"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type DynamodbKeyValueStore struct {
	pulumi.ResourceState

	Table *dynamodb.Table
	Name  string
}

type DynamodbKeyValueStoreArgs struct {
	StackID       string
	KeyValueStore *v1.KeyValueStore
}

func (n *NitricAwsPulumiProvider) KeyValueStore(ctx *pulumi.Context, parent pulumi.Resource, name string, keyvalue *deploymentspb.KeyValueStore) error {
	var err error
	opts := []pulumi.ResourceOption{pulumi.Parent(parent)}

	n.KeyValueStores[name], err = dynamodb.NewTable(ctx, name, &dynamodb.TableArgs{
		Attributes: dynamodb.TableAttributeArray{
			&dynamodb.TableAttributeArgs{
				Name: pulumi.String("_pk"),
				Type: pulumi.String("S"),
			},
			&dynamodb.TableAttributeArgs{
				Name: pulumi.String("_sk"),
				Type: pulumi.String("S"),
			},
		},
		HashKey:     pulumi.String("_pk"),
		RangeKey:    pulumi.String("_sk"),
		BillingMode: pulumi.String("PAY_PER_REQUEST"),
		Tags:        pulumi.ToStringMap(tags.Tags(n.StackId, name, resources.Collection)),
	}, opts...)

	return err
}
