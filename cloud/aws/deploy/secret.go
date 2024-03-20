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
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/resourcegroupstaggingapi"
	"github.com/nitrictech/nitric/cloud/common/deploy/resources"
	common "github.com/nitrictech/nitric/cloud/common/deploy/tags"
	deploymentspb "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/secretsmanager"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// tagSecret - tags an existing secret in AWS and adds it to the stack.
func tagSecret(ctx *pulumi.Context, name string, importArn string, tags map[string]string, client *resourcegroupstaggingapi.ResourceGroupsTaggingAPI) (*secretsmanager.Secret, error) {
	secretLookup, err := secretsmanager.LookupSecret(ctx, &secretsmanager.LookupSecretArgs{
		Arn: aws.String(importArn),
	})
	if err != nil {
		return nil, err
	}

	_, err = client.TagResources(&resourcegroupstaggingapi.TagResourcesInput{
		ResourceARNList: aws.StringSlice([]string{secretLookup.Arn}),
		Tags:            aws.StringMap(tags),
	})
	if err != nil {
		return nil, err
	}

	sec, err := secretsmanager.GetSecret(
		ctx,
		name,
		pulumi.ID(secretLookup.Id),
		nil,
		// nitric didn't create this resource, so it shouldn't delete it either.
		pulumi.RetainOnDelete(true),
	)
	if err != nil {
		return nil, err
	}
	return sec, nil
}

// createSecret - creates a new secret in AWS, using the provided name and tags.
func createSecret(ctx *pulumi.Context, name string, tags map[string]string) (*secretsmanager.Secret, error) {
	sec, err := secretsmanager.NewSecret(ctx, name, &secretsmanager.SecretArgs{
		Tags: pulumi.ToStringMap(tags),
	})
	if err != nil {
		return nil, err
	}

	return sec, nil
}

// Secret - Implements deployments of Nitric Secrets using AWS Secrets Manager
func (a *NitricAwsPulumiProvider) Secret(ctx *pulumi.Context, parent pulumi.Resource, name string, config *deploymentspb.Secret) error {
	awsTags := common.Tags(a.stackId, name, resources.Secret)

	var err error
	var secret *secretsmanager.Secret

	importArn := ""
	if a.config.Import.Secrets != nil {
		importArn = a.config.Import.Secrets[name]
	}

	if importArn != "" {
		secret, err = tagSecret(ctx, name, importArn, awsTags, a.resourceTaggingClient)
	} else {
		secret, err = createSecret(ctx, name, awsTags)
	}

	if err != nil {
		return err
	}

	a.secrets[name] = secret

	return nil
}
