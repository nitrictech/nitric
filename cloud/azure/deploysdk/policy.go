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
	"github.com/pulumi/pulumi-azure-native-sdk/keyvault"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"

	"github.com/nitrictech/nitric/cloud/azure/deploy/bucket"
	"github.com/nitrictech/nitric/cloud/azure/deploy/exec"
	"github.com/nitrictech/nitric/cloud/azure/deploy/topic"
	deploymentspb "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"
	resourcespb "github.com/nitrictech/nitric/core/pkg/proto/resources/v1"
)

type Policy struct {
	pulumi.ResourceState

	Name         string
	RolePolicies []*iam.RolePolicy
}

type StackResources struct {
	SubscriptionId pulumi.StringInput
	Topics         map[string]*topic.AzureEventGridTopic
	Buckets        map[string]*bucket.AzureStorageBucket
	// Collections    map[string]*documentdb.MongoDBResourceMongoDBCollection
	// The vault that all secrets are stored in
	KeyVault *keyvault.Vault
}

type PrincipalMap = map[resourcespb.ResourceType]map[string]*exec.ServicePrincipal

type PolicyArgs struct {
	ResourceGroupName pulumi.StringInput

	Policy *deploymentspb.Policy
	// Nitric Action to AzureAD role mappings
	// AvailableRoles map[v1.Action]*authorization.RoleDefinition
	// Nitric roles
	Roles *Roles
	// Resources in the stack that must be protected
	Resources *StackResources
	// Resources in the stack that may act as actors
	Principals PrincipalMap
}

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
		if resource.Id.Type == resourcespb.ResourceType_Collection {
			continue
		}

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

// func NewAzureADPolicy(ctx *pulumi.Context, name string, args *PolicyArgs, opts ...pulumi.ResourceOption) (*Policy, error) {
// 	res := &Policy{Name: name, RolePolicies: make([]*iam.RolePolicy, 0)}

// 	err := ctx.RegisterComponentResource("nitric:policy:AazureADPolicy", name, res, opts...)
// 	if err != nil {
// 		return nil, err
// 	}

// 	for _, resource := range args.Policy.Resources {
// 		if resource.Id.Type == resourcespb.ResourceType_Collection {
// 			continue
// 		}

// 		for _, principal := range args.Policy.Principals {
// 			// The roles we need to assign
// 			roles := actionsToAzureRoleDefinitions(args.Roles.RoleDefinitions, args.Policy.Actions)
// 			if len(roles) == 0 {
// 				return nil, fmt.Errorf("policy contained not assignable actions %+v, %+v", args.Policy, args.Roles.RoleDefinitions)
// 			}

// 			sp, ok := args.Principals[principal.Id.Type][principal.Id.Name]
// 			if !ok {
// 				return nil, fmt.Errorf("principal %s of type %s not found", principal.Id.Name, principal.Id.Type)
// 			}

// 			// We have the principal and the roles we need to assign
// 			// just need to scope the resource type to the RoleAssignments
// 			for roleName, role := range roles {
// 				// FIXME: Implement collection and secret least priveledge
// 				scope, err := p.scopeFromResource(resource)
// 				if err != nil {
// 					return nil, err
// 				}

// 				_, err = authorization.NewRoleAssignment(ctx, fmt.Sprintf("%s-%s", principal.Id.Name, roleName), &authorization.RoleAssignmentArgs{
// 					PrincipalId:      sp.ServicePrincipalId,
// 					PrincipalType:    pulumi.String("ServicePrincipal"),
// 					RoleDefinitionId: role.ID(),
// 					// Convert the target resources into a scope
// 					Scope:     scope.scope,
// 					Condition: scope.condition,
// 				}, pulumi.Parent(res))
// 				if err != nil {
// 					return nil, fmt.Errorf("there was an error creating the role assignment: %w", err)
// 				}
// 			}
// 		}
// 	}

// 	return res, nil
// }
