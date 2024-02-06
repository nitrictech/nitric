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
	//#nosec G501 -- md5 used only to produce a unique ID from non-sensistive information (policy IDs)

	"fmt"

	iam "github.com/pulumi/pulumi-aws/sdk/v5/go/aws/iam"
	"github.com/pulumi/pulumi-azure-native-sdk/authorization"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"

	deploymentspb "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"
	resourcespb "github.com/nitrictech/nitric/core/pkg/proto/resources/v1"
)

type Policy struct {
	pulumi.ResourceState

	Name         string
	RolePolicies []*iam.RolePolicy
}

type PrincipalMap = map[resourcespb.ResourceType]map[string]*ServicePrincipal

func actionsToAzureRoleDefinitions(roles map[resourcespb.Action]*authorization.RoleDefinition, actions []resourcespb.Action) map[string]*authorization.RoleDefinition {
	azureRoles := map[string]*authorization.RoleDefinition{}

	for _, a := range actions {
		if role, ok := roles[a]; ok {
			azureRoles[resourcespb.Action_name[int32(a)]] = role
		}
	}

	return azureRoles
}

type resourceScope struct {
	scope     pulumi.StringInput
	condition pulumi.StringInput
}

// "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/[{parentResourcePath}/]{resourceType}/{resourceName}"
func (p *NitricAzurePulumiProvider) scopeFromResource(resource *deploymentspb.Resource) (*resourceScope, error) {
	switch resource.Id.Type {
	case resourcespb.ResourceType_Topic:
		topic, ok := p.topics[resource.Id.Name]
		if !ok {
			return nil, fmt.Errorf("topic %s not found", resource.Id.Name)
		}

		return &resourceScope{
			scope: pulumi.Sprintf(
				"subscriptions/%s/resourceGroups/%s/providers/Microsoft.EventGrid/topics/%s",
				p.clientConfig.SubscriptionId,
				p.resourceGroup.Name,
				topic.Name,
			),
		}, nil
	case resourcespb.ResourceType_KeyValueStore:
		kv, ok := p.keyValueStores[resource.Id.Name]
		if !ok {
			return nil, fmt.Errorf("key value store %s not found", resource.Id.Name)
		}

		return &resourceScope{
			scope: pulumi.Sprintf(
				"subscriptions/%s/resourceGroups/%s/providers/Microsoft.Storage/storageAccounts/%s/tableServices/default/tables/%s",
				p.clientConfig.SubscriptionId,
				p.resourceGroup.Name,
				p.storageAccount.Name,
				kv.Name,
			),
		}, nil
	case resourcespb.ResourceType_Bucket:
		bucket, ok := p.buckets[resource.Id.Name]
		if !ok {
			return nil, fmt.Errorf("bucket %s not found", resource.Id.Name)
		}

		// return pulumi.Sprintf(
		// 	"/subscriptions/%s/resourceGroups/%s",
		// 	deployedResources.SubscriptionId,
		// 	bucket.ResourceGroup.Name,
		// ), nil

		return &resourceScope{
			scope: pulumi.Sprintf(
				"subscriptions/%s/resourceGroups/%s/providers/Microsoft.Storage/storageAccounts/%s/blobServices/default/containers/%s",
				p.clientConfig.SubscriptionId,
				p.resourceGroup.Name,
				p.storageAccount.Name,
				bucket.Name,
			),
		}, nil
	case resourcespb.ResourceType_Queue:
		queue, ok := p.queues[resource.Id.Name]
		if !ok {
			return nil, fmt.Errorf("queue %s not found", resource.Id.Name)
		}

		return &resourceScope{
			scope: pulumi.Sprintf(
				"subscriptions/%s/resourceGroups/%s/providers/Microsoft.Storage/storageAccounts/%s/queueServices/default/queues/%s",
				p.clientConfig.SubscriptionId,
				p.resourceGroup.Name,
				p.storageAccount.Name,
				queue.Name,
			),
		}, nil
	case resourcespb.ResourceType_Secret:
		if p.keyVault == nil {
			return nil, fmt.Errorf("secret %s not found", resource.Id.Name)
		}

		return &resourceScope{
			scope: pulumi.Sprintf(
				"subscriptions/%s/resourcegroups/%s/providers/Microsoft.KeyVault/vaults/%s/secrets/%s",
				p.clientConfig.SubscriptionId,
				p.resourceGroup.Name,
				p.keyVault.Name,
				resource.Id.Name,
			),
			// condition: pulumi.Sprintf("@Resource[Microsoft.KeyVault/vaults/secrets].name equals %s'", resource.Name),
		}, nil
		// TODO
	default:
		return nil, fmt.Errorf("unknown resource type %s", resource.Id.Type)
	}
}

func (p *NitricAzurePulumiProvider) Policy(ctx *pulumi.Context, parent pulumi.Resource, name string, policy *deploymentspb.Policy) error {
	opts := []pulumi.ResourceOption{pulumi.Parent(parent)}

	for _, resource := range policy.Resources {
		for _, principal := range policy.Principals {
			// The roles we need to assign
			roles := actionsToAzureRoleDefinitions(p.roles.RoleDefinitions, policy.Actions)
			if len(roles) == 0 {
				return fmt.Errorf("policy contained not assignable actions %+v, %+v", policy, p.roles.RoleDefinitions)
			}

			sp, ok := p.principals[principal.Id.Type][principal.Id.Name]
			if !ok {
				return fmt.Errorf("principal %s of type %s not found", principal.Id.Name, principal.Id.Type)
			}

			// We have the principal and the roles we need to assign
			// just need to scope the resource type to the RoleAssignments
			for roleName, role := range roles {
				// FIXME: Implement collection and secret least priveledge
				scope, err := p.scopeFromResource(resource)
				if err != nil {
					return err
				}

				_, err = authorization.NewRoleAssignment(ctx, fmt.Sprintf("%s-%s", principal.Id.Name, roleName), &authorization.RoleAssignmentArgs{
					PrincipalId:      sp.ServicePrincipalId,
					PrincipalType:    pulumi.String("ServicePrincipal"),
					RoleDefinitionId: role.ID(),
					// Convert the target resources into a scope
					Scope:     scope.scope,
					Condition: scope.condition,
				}, opts...)
				if err != nil {
					return fmt.Errorf("there was an error creating the role assignment: %w", err)
				}
			}
		}
	}

	return nil
}
