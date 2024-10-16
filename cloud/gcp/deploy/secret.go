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
	gcpsecretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"github.com/nitrictech/nitric/cloud/common/deploy/resources"
	common "github.com/nitrictech/nitric/cloud/common/deploy/tags"
	deploymentspb "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"
	"github.com/pulumi/pulumi-gcp/sdk/v8/go/gcp/secretmanager"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

type SecretManagerSecret struct {
	pulumi.ResourceState

	Name   string
	Secret *secretmanager.Secret
}

// tagSecret - tags an existing secret in GCP and adds it to the stack.
func tagSecret(ctx *pulumi.Context, name string, projectId string, secretId string, tags map[string]string, client *gcpsecretmanager.Client, opts []pulumi.ResourceOption) (*secretmanager.Secret, error) {
	secretLookup, err := secretmanager.LookupSecret(ctx, &secretmanager.LookupSecretArgs{
		Project:  &projectId,
		SecretId: secretId,
	})
	if err != nil {
		return nil, err
	}

	_, err = client.UpdateSecret(ctx.Context(), &secretmanagerpb.UpdateSecretRequest{
		Secret: &secretmanagerpb.Secret{
			Name:   secretLookup.Name,
			Labels: tags,
		},
		UpdateMask: &fieldmaskpb.FieldMask{
			Paths: []string{"labels"},
		},
	})
	if err != nil {
		return nil, err
	}

	sec, err := secretmanager.GetSecret(
		ctx,
		name,
		pulumi.ID(secretLookup.Id),
		nil,
		// nitric didn't create this resource, so it shouldn't delete it either.
		append(opts, pulumi.RetainOnDelete(true))...,
	)
	if err != nil {
		return nil, err
	}
	return sec, nil
}

// createSecret - creates a new secret in GCP Secret Manager, using the provided name and tags.
func createSecret(ctx *pulumi.Context, name string, stackName string, tags map[string]string, opts []pulumi.ResourceOption) (*secretmanager.Secret, error) {
	secId := pulumi.Sprintf("%s-%s", stackName, name)
	sec, err := secretmanager.NewSecret(ctx, name, &secretmanager.SecretArgs{
		SecretId: secId,
		Labels:   pulumi.ToStringMap(tags),
	}, opts...)
	if err != nil {
		return nil, err
	}

	return sec, nil
}

func (p *NitricGcpPulumiProvider) Secret(ctx *pulumi.Context, parent pulumi.Resource, name string, config *deploymentspb.Secret) error {
	var err error
	opts := append([]pulumi.ResourceOption{}, pulumi.Parent(parent))

	secretLabels := common.Tags(p.StackId, name, resources.Secret)

	var secret *secretmanager.Secret

	importId := ""
	if p.GcpConfig.Import.Secrets != nil {
		importId = p.GcpConfig.Import.Secrets[name]
	}

	if importId != "" {
		secret, err = tagSecret(ctx, name, p.GcpConfig.ProjectId, importId, secretLabels, p.SecretManagerClient, p.WithDefaultResourceOptions(opts...))
	} else {
		secret, err = createSecret(ctx, name, p.StackName, secretLabels, p.WithDefaultResourceOptions(opts...))
	}

	if err != nil {
		return err
	}

	p.Secrets[name] = secret

	return nil
}
