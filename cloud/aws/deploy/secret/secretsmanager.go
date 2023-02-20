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
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/secretsmanager"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type SecretsManagerSecret struct {
	pulumi.ResourceState
	SecretsManager *secretsmanager.Secret
	Name           string
}

type SecretsManagerSecretArgs struct {
	StackID pulumi.StringInput
	Secret  *v1.Secret
}

// Create a new SecretsManager secret
func NewSecretsManagerSecret(ctx *pulumi.Context, name string, args *SecretsManagerSecretArgs, opts ...pulumi.ResourceOption) (*SecretsManagerSecret, error) {
	res := &SecretsManagerSecret{
		Name: name,
	}

	err := ctx.RegisterComponentResource("nitric:secrete:AwsSecretsManager", name, res, opts...)
	if err != nil {
		return nil, err
	}

	sec, err := secretsmanager.NewSecret(ctx, name, &secretsmanager.SecretArgs{
		Tags: common.Tags(ctx, args.StackID, name),
	})
	if err != nil {
		return nil, err
	}

	res.SecretsManager = sec

	return res, nil
}
