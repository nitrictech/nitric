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

package stack

import (
	"encoding/json"

	"github.com/nitrictech/nitric/cloud/common/deploy/resources"
	"github.com/nitrictech/nitric/cloud/common/deploy/tags"

	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/resourcegroups"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type AwsResourceGroup struct {
	pulumi.ResourceState
	Name          string
	ResourceGroup *resourcegroups.Group
}

type AwsResourceGroupArgs struct {
	StackID string
}

func NewAwsResourceGroup(ctx *pulumi.Context, name string, args *AwsResourceGroupArgs, opts ...pulumi.ResourceOption) (*AwsResourceGroup, error) {
	res := &AwsResourceGroup{Name: name}

	err := ctx.RegisterComponentResource("nitric:stack:AwsResourceGroup", name, res, opts...)
	if err != nil {
		return nil, err
	}

	rgQueryJSON := pulumi.String(args.StackID).ToStringOutput().ApplyT(func(sid string) (string, error) {
		b, err := json.Marshal(map[string]interface{}{
			"ResourceTypeFilters": []string{"AWS::AllSupported"},
			"TagFilters": []interface{}{
				map[string]interface{}{
					//"Key":    "x-nitric-stack",
					//"Values": []string{sid},
					// TODO: validate this key only filter works
					"Key": tags.GetResourceNameKey(sid),
				},
			},
		})
		if err != nil {
			return "", err
		}

		return string(b), nil
	}).(pulumi.StringOutput)

	res.ResourceGroup, err = resourcegroups.NewGroup(ctx, "rg-"+ctx.Stack(), &resourcegroups.GroupArgs{
		Description: pulumi.Sprintf("Nitric stack %s resources", res.Name),
		Tags:        pulumi.ToStringMap(tags.Tags(args.StackID, name, resources.Stack)),
		ResourceQuery: &resourcegroups.GroupResourceQueryArgs{
			Query: rgQueryJSON,
		},
	}, pulumi.Parent(res))
	if err != nil {
		return nil, err
	}

	return res, nil
}
