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

package api

import (
	"github.com/nitrictech/nitric/cloud/common/deploy/pulumix"
	resourcespb "github.com/nitrictech/nitric/core/pkg/proto/resources/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type NitricApi struct {
	pulumi.ResourceState

	Name string
}

type NitricApiArgs struct{}

type ApiDeploymentFunc = func() error

func NewNitricApi(ctx *pulumi.Context, name string, args NitricApiArgs, opts ...pulumi.ResourceOption) (*NitricApi, error) {
	res := &NitricApi{Name: name}

	err := ctx.RegisterComponentResource(pulumix.PulumiUrn(resourcespb.ResourceType_Api), name, res, opts...)
	if err != nil {
		return nil, err
	}

	return res, err
}
