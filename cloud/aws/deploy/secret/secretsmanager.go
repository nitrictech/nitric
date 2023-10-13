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
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/resourcegroupstaggingapi"
	"github.com/nitrictech/nitric/cloud/common/deploy/resources"
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
	StackID string
	// Import an existing secret
	Import string
	Secret *v1.Secret
	Client *resourcegroupstaggingapi.ResourceGroupsTaggingAPI
}

// Create a new SecretsManager secret
func NewSecretsManagerSecret(ctx *pulumi.Context, name string, args *SecretsManagerSecretArgs, opts ...pulumi.ResourceOption) (*SecretsManagerSecret, error) {
	res := &SecretsManagerSecret{
		Name: name,
	}

	err := ctx.RegisterComponentResource("nitric:secret:AwsSecretsManager", name, res, opts...)
	if err != nil {
		return nil, err
	}

	tags := common.Tags(args.StackID, name, resources.Secret)

	if args.Import != "" {
		secretLookup, err := secretsmanager.LookupSecret(ctx, &secretsmanager.LookupSecretArgs{
			Arn: aws.String(args.Import),
		})
		if err != nil {
			return nil, err
		}

		_, err = args.Client.TagResources(&resourcegroupstaggingapi.TagResourcesInput{
			ResourceARNList: aws.StringSlice([]string{secretLookup.Arn}),
			Tags:            aws.StringMap(tags),
		})

		if err != nil {
			return nil, err
		}

		// import an existing secret
		res.SecretsManager, err = secretsmanager.GetSecret(
			ctx,
			name,
			pulumi.ID(secretLookup.Id),
			nil,
			// not our resource so we'll keep it around
			pulumi.RetainOnDelete(true),
		)
		if err != nil {
			return nil, err
		}
	} else {
		// create a new secret
		res.SecretsManager, err = secretsmanager.NewSecret(ctx, name, &secretsmanager.SecretArgs{
			Tags: pulumi.ToStringMap(common.Tags(args.StackID, name, resources.Secret)),
		})
		if err != nil {
			return nil, err
		}
	}

	return res, nil
}
