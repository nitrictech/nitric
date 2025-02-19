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

package deploytf

import (
	"fmt"

	"github.com/aws/jsii-runtime-go"
	"github.com/hashicorp/terraform-cdk-go/cdktf"
	"github.com/nitrictech/nitric/cloud/azure/deploytf/generated/policy"
	deploymentspb "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"
	resourcespb "github.com/nitrictech/nitric/core/pkg/proto/resources/v1"
)

func (p *NitricAzureTerraformProvider) actionsToAzureRoleDefinitions(actions []resourcespb.Action) map[string]*string {
	azureRoles := map[string]*string{}

	for _, a := range actions {
		switch a {
		case resourcespb.Action_BucketFileGet:
			azureRoles[resourcespb.Action_BucketFileGet.String()] = p.Roles.BucketReadOutput()
		case resourcespb.Action_BucketFilePut:
			azureRoles[resourcespb.Action_BucketFilePut.String()] = p.Roles.BucketWriteOutput()
		case resourcespb.Action_BucketFileDelete:
			azureRoles[resourcespb.Action_BucketFileDelete.String()] = p.Roles.BucketDeleteOutput()
		case resourcespb.Action_BucketFileList:
			azureRoles[resourcespb.Action_BucketFileList.String()] = p.Roles.BucketListOutput()
		case resourcespb.Action_KeyValueStoreRead:
			azureRoles[resourcespb.Action_KeyValueStoreRead.String()] = p.Roles.KvReadOutput()
		case resourcespb.Action_KeyValueStoreWrite:
			azureRoles[resourcespb.Action_KeyValueStoreWrite.String()] = p.Roles.KvWriteOutput()
		case resourcespb.Action_KeyValueStoreDelete:
			azureRoles[resourcespb.Action_KeyValueStoreDelete.String()] = p.Roles.KvDeleteOutput()
		case resourcespb.Action_TopicPublish:
			azureRoles[resourcespb.Action_TopicPublish.String()] = p.Roles.TopicPublishOutput()
		case resourcespb.Action_QueueEnqueue:
			azureRoles[resourcespb.Action_QueueEnqueue.String()] = p.Roles.QueueEnqueueOutput()
		case resourcespb.Action_QueueDequeue:
			azureRoles[resourcespb.Action_QueueDequeue.String()] = p.Roles.QueueDequeueOutput()
		case resourcespb.Action_SecretAccess:
			azureRoles[resourcespb.Action_SecretAccess.String()] = p.Roles.SecretAccessOutput()
		case resourcespb.Action_SecretPut:
			azureRoles[resourcespb.Action_SecretPut.String()] = p.Roles.SecretPutOutput()
		}
	}

	return azureRoles
}

type ResourceScope struct {
	Scope      *string `json:"scope"`
	Condition  *string `json:"condition"`
	Dependency cdktf.ITerraformDependable
}

func (p *NitricAzureTerraformProvider) scopeFromResource(resource *deploymentspb.Resource) (*ResourceScope, error) {
	switch resource.Id.Type {
	case resourcespb.ResourceType_Topic:
		topic, ok := p.Topics[resource.Id.Name]
		if !ok {
			return nil, fmt.Errorf("topic %s not found", resource.Id.Name)
		}

		return &ResourceScope{
			Scope: jsii.Sprintf(
				"subscriptions/%s/resourceGroups/%s/providers/Microsoft.EventGrid/topics/%s",
				*p.Stack.SubscriptionIdOutput(),
				*p.Stack.ResourceGroupNameOutput(),
				*topic.Name(),
			),
			Dependency: topic,
		}, nil
	case resourcespb.ResourceType_KeyValueStore:
		kv, ok := p.KvStores[resource.Id.Name]
		if !ok {
			return nil, fmt.Errorf("key value store %s not found", resource.Id.Name)
		}

		return &ResourceScope{
			Scope: jsii.Sprintf(
				"subscriptions/%s/resourceGroups/%s/providers/Microsoft.Storage/storageAccounts/%s/tableServices/default/tables/%s",
				*p.Stack.SubscriptionIdOutput(),
				*p.Stack.ResourceGroupNameOutput(),
				*p.Stack.StorageAccountNameOutput(),
				*kv.Name(),
			),
			Dependency: kv,
		}, nil
	case resourcespb.ResourceType_Bucket:
		bucket, ok := p.Buckets[resource.Id.Name]
		if !ok {
			return nil, fmt.Errorf("bucket %s not found", resource.Id.Name)
		}
		return &ResourceScope{
			Scope: jsii.Sprintf(
				"subscriptions/%s/resourceGroups/%s/providers/Microsoft.Storage/storageAccounts/%s/blobServices/default/containers/%s",
				*p.Stack.SubscriptionIdOutput(),
				*p.Stack.ResourceGroupNameOutput(),
				*p.Stack.StorageAccountNameOutput(),
				*bucket.Name(),
			),
			Dependency: bucket,
		}, nil
	case resourcespb.ResourceType_Queue:
		queue, ok := p.Queues[resource.Id.Name]
		if !ok {
			return nil, fmt.Errorf("queue %s not found", resource.Id.Name)
		}

		return &ResourceScope{
			Scope: jsii.Sprintf(
				"subscriptions/%s/resourceGroups/%s/providers/Microsoft.Storage/storageAccounts/%s/queueServices/default/queues/%s",
				*p.Stack.SubscriptionIdOutput(),
				*p.Stack.ResourceGroupNameOutput(),
				*p.Stack.StorageAccountNameOutput(),
				*queue.Name(),
			),
			Dependency: queue,
		}, nil
	case resourcespb.ResourceType_Secret:
		if !*p.Stack.EnableKeyvault() {
			return nil, fmt.Errorf("secret %s not found", resource.Id.Name)
		}

		return &ResourceScope{
			Scope: jsii.Sprintf(
				"subscriptions/%s/resourcegroups/%s/providers/Microsoft.KeyVault/vaults/%s/secrets/%s",
				*p.Stack.SubscriptionIdOutput(),
				*p.Stack.ResourceGroupNameOutput(),
				*p.Stack.KeyvaultNameOutput(),
				resource.Id.Name,
			),
			Dependency: p.Stack,
		}, nil
	default:
		return nil, fmt.Errorf("unknown resource type %s", resource.Id.Type)
	}
}

func (a *NitricAzureTerraformProvider) Policy(stack cdktf.TerraformStack, name string, config *deploymentspb.Policy) error {
	for _, resource := range config.Resources {
		for _, principal := range config.Principals {
			if principal.Id.Type != resourcespb.ResourceType_Service {
				return fmt.Errorf("only service principals are supported")
			}

			// The roles we need to assign
			roles := a.actionsToAzureRoleDefinitions(config.Actions)
			if len(roles) == 0 {
				return fmt.Errorf("policy contained not assignable actions %+v, %+v", config, a.Roles)
			}

			svc, ok := a.Services[principal.Id.Name]
			if !ok {
				return fmt.Errorf("principal %s of type %s not found", principal.Id.Name, principal.Id.Type)
			}

			spId := svc.ServicePrincipalIdOutput()

			// We have the principal and the roles we need to assign
			// just need to scope the resource type to the RoleAssignments
			for roleName, role := range roles {
				scope, err := a.scopeFromResource(resource)
				if err != nil {
					return err
				}

				policy.NewPolicy(stack, jsii.Sprintf("%s-%s", principal.Id.Name, roleName), &policy.PolicyConfig{
					ServicePrincipalId: spId,
					Scope:              scope.Scope,
					RoleDefinitionId:   role,
					DependsOn:          &[]cdktf.ITerraformDependable{scope.Dependency, a.Roles},
				})
			}
		}
	}

	return nil
}
