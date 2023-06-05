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

package secret

import (
	common "github.com/nitrictech/nitric/cloud/common/deploy/tags"
	v1 "github.com/nitrictech/nitric/core/pkg/api/nitric/deploy/v1"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/secretmanager"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type SecretManagerSecret struct {
	pulumi.ResourceState

	Name   string
	Secret *secretmanager.Secret
}

type SecretManagerSecretArgs struct {
	Location  string
	StackID   string
	StackName string

	Secret *v1.Secret
}

func NewSecretManagerSecret(ctx *pulumi.Context, name string, args *SecretManagerSecretArgs, opts ...pulumi.ResourceOption) (*SecretManagerSecret, error) {
	res := &SecretManagerSecret{
		Name: name,
	}

	err := ctx.RegisterComponentResource("nitric:secret:GCPSecretManager", name, res, opts...)
	if err != nil {
		return nil, err
	}

	secId := pulumi.Sprintf("%s-%s", args.StackName, name)

	res.Secret, err = secretmanager.NewSecret(ctx, name, &secretmanager.SecretArgs{
		Replication: secretmanager.SecretReplicationArgs{
			Automatic: pulumi.Bool(true),
		},
		SecretId: secId,
		Labels:   pulumi.ToStringMap(common.Tags(args.StackID, name)),
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}
