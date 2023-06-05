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

package collection

import (
	"github.com/nitrictech/nitric/cloud/common/deploy/tags"
	v1 "github.com/nitrictech/nitric/core/pkg/api/nitric/deploy/v1"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/dynamodb"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type DynamodbCollection struct {
	pulumi.ResourceState

	Table *dynamodb.Table
	Name  string
}

type DynamodbCollectionArgs struct {
	StackID    string
	Collection *v1.Collection
}

func NewDynamodbCollection(ctx *pulumi.Context, name string, args *DynamodbCollectionArgs, opts ...pulumi.ResourceOption) (*DynamodbCollection, error) {
	res := &DynamodbCollection{Name: name}

	err := ctx.RegisterComponentResource("nitric:collection:Dynamodb", name, res, opts...)
	if err != nil {
		return nil, err
	}

	res.Table, err = dynamodb.NewTable(ctx, name, &dynamodb.TableArgs{
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
		Tags:        pulumi.ToStringMap(tags.Tags(args.StackID, name)),
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}
