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
	common "github.com/nitrictech/nitric/cloud/common/deploy/tags"
	deploymentspb "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/secretmanager"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type SecretManagerSecret struct {
	pulumi.ResourceState

	Name   string
	Secret *secretmanager.Secret
}

func (p *NitricGcpPulumiProvider) Secret(ctx *pulumi.Context, parent pulumi.Resource, name string, config *deploymentspb.Secret) error {
	var err error
	opts := append([]pulumi.ResourceOption{}, pulumi.Parent(parent))

	secId := pulumi.Sprintf("%s-%s", p.stackName, name)

	p.secrets[name], err = secretmanager.NewSecret(ctx, name, &secretmanager.SecretArgs{
		Replication: secretmanager.SecretReplicationArgs{
			Automatic: pulumi.Bool(true),
		},
		SecretId: secId,
		Labels:   pulumi.ToStringMap(common.Tags(p.stackId, name, resources.Secret)),
	}, opts...)
	if err != nil {
		return err
	}

	return nil
}
