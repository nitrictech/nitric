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
	"fmt"

	deploymentspb "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"
	resourcespb "github.com/nitrictech/nitric/core/pkg/proto/resources/v1"
	v1 "github.com/nitrictech/nitric/core/pkg/proto/resources/v1"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/projects"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/pubsub"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/secretmanager"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/serviceaccount"
	gcpstorage "github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/storage"
	"github.com/pulumi/pulumi-random/sdk/v4/go/random"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Policy struct {
	pulumi.ResourceState

	Name         string
	RolePolicies []*projects.IAMMember
}

var gcpActionsMap map[v1.Action][]string = map[v1.Action][]string{
	v1.Action_BucketFileList: {
		"storage.objects.list",
	},
	v1.Action_BucketFileGet: {
		"storage.objects.get",
	},
	v1.Action_BucketFilePut: {
		"orgpolicy.policy.get",
		"storage.multipartUploads.abort",
		"storage.multipartUploads.create",
		"storage.multipartUploads.listParts",
		"storage.objects.create",
	},
	v1.Action_BucketFileDelete: {
		"storage.objects.delete",
	},
	v1.Action_TopicDetail: {
		"pubsub.topics.get",
	},
	v1.Action_TopicEventPublish: {
		"pubsub.topics.publish",
	},
	v1.Action_TopicList: {}, // see above in gcpListActions
	v1.Action_CollectionDocumentDelete: {
		"appengine.applications.get",
		"datastore.databases.get",
		"datastore.indexes.get",
		"datastore.namespaces.get",
		"datastore.entities.delete",
	},
	v1.Action_CollectionDocumentRead: {
		"appengine.applications.get",
		"datastore.databases.get",
		"datastore.entities.get",
		"datastore.indexes.get",
		"datastore.namespaces.get",
		"datastore.entities.list",
	},
	v1.Action_CollectionDocumentWrite: {
		"appengine.applications.get",
		"datastore.indexes.list",
		"datastore.namespaces.list",
		"datastore.entities.create",
		"datastore.entities.update",
	},
	v1.Action_CollectionQuery: {
		"appengine.applications.get",
		"datastore.databases.get",
		"datastore.entities.get",
		"datastore.entities.list",
		"datastore.indexes.get",
		"datastore.namespaces.get",
	},
	v1.Action_CollectionList: {
		"appengine.applications.get",
	},
	v1.Action_SecretAccess: {
		"resourcemanager.projects.get",
		"secretmanager.locations.get",
		"secretmanager.locations.list",
		"secretmanager.secrets.get",
		"secretmanager.secrets.getIamPolicy",
		"secretmanager.versions.get",
		"secretmanager.versions.access",
		"secretmanager.versions.list",
	},
	v1.Action_SecretPut: {
		"resourcemanager.projects.get",
		"secretmanager.versions.add",
		"secretmanager.versions.enable",
		"secretmanager.versions.destroy",
		"secretmanager.versions.disable",
		"secretmanager.versions.get",
		"secretmanager.versions.access",
		"secretmanager.versions.list",
	},
}

var collectionActions []string = nil

func getCollectionActions() []string {
	if collectionActions == nil {
		collectionActions = make([]string, 0)
		collectionActions = append(collectionActions, gcpActionsMap[v1.Action_CollectionDocumentRead]...)
		collectionActions = append(collectionActions, gcpActionsMap[v1.Action_CollectionDocumentWrite]...)
		collectionActions = append(collectionActions, gcpActionsMap[v1.Action_CollectionDocumentDelete]...)
	}

	return collectionActions
}

func filterCollectionActions(actions []string) []string {
	filteredActions := []string{}

	for _, a := range actions {
		for _, ca := range getCollectionActions() {
			if a == ca {
				filteredActions = append(filteredActions, a)
				break
			}
		}
	}

	return filteredActions
}

func actionsToGcpActions(actions []v1.Action) []string {
	gcpActions := make([]string, 0)

	for _, a := range actions {
		gcpActions = append(gcpActions, gcpActionsMap[a]...)
	}

	return gcpActions
}

func (a *NitricGcpPulumiProvider) serviceAccountForPrincipal(resource *deploymentspb.Resource) (*serviceaccount.Account, error) {
	switch resource.Id.Type {
	case resourcespb.ResourceType_ExecUnit:
		if f, ok := a.cloudRunServices[resource.Id.Name]; ok {
			return f.ServiceAccount, nil
		}
	default:
		return nil, fmt.Errorf("could not find role for principal: %+v", resource)
	}

	return nil, fmt.Errorf("could not find role for principal: %+v", resource)
}

func (p *NitricGcpPulumiProvider) Policy(ctx *pulumi.Context, parent pulumi.Resource, name string, config *deploymentspb.Policy) error {
	actions := actionsToGcpActions(config.Actions)
	opts := []pulumi.ResourceOption{pulumi.Parent(parent)}

	rolePolicy, err := NewCustomRole(ctx, name, actions, opts...)
	if err != nil {
		return err
	}

	for _, principal := range config.Principals {
		sa, err := p.serviceAccountForPrincipal(principal)
		if err != nil {
			return err
		}

		for _, resource := range config.Resources {
			memberName := fmt.Sprintf("%s-%s", principal.Id.Name, resource.Id.Name)
			memberId := pulumi.Sprintf("serviceAccount:%s", sa.Email)

			switch resource.Id.Type {
			case v1.ResourceType_Bucket:
				b := p.buckets[resource.Id.Name]

				_, err = gcpstorage.NewBucketIAMMember(ctx, memberName, &gcpstorage.BucketIAMMemberArgs{
					Bucket: b.Name,
					Member: memberId,
					Role:   rolePolicy.Name,
				}, opts...)
				if err != nil {
					return err
				}

			case v1.ResourceType_Collection:
				collActions := filterCollectionActions(actions)

				collRole, err := NewCustomRole(ctx, memberName+"-role", collActions, opts...)
				if err != nil {
					return err
				}

				_, err = projects.NewIAMMember(ctx, memberName, &projects.IAMMemberArgs{
					Member:  memberId,
					Project: pulumi.String(p.config.ProjectId),
					Role:    collRole.Name,
				}, opts...)
				if err != nil {
					return err
				}
			case v1.ResourceType_Topic:
				t := p.topics[resource.Id.Name]

				_, err = pubsub.NewTopicIAMMember(ctx, memberName, &pubsub.TopicIAMMemberArgs{
					Topic:  t.Name,
					Member: memberId,
					Role:   rolePolicy.Name,
				}, opts...)
				if err != nil {
					return err
				}

			case v1.ResourceType_Secret:
				s := p.secrets[resource.Id.Name]

				_, err = secretmanager.NewSecretIamMember(ctx, memberName, &secretmanager.SecretIamMemberArgs{
					SecretId: s.SecretId,
					Member:   memberId,
					Role:     rolePolicy.Name,
				}, opts...)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func NewCustomRole(ctx *pulumi.Context, name string, actions []string, opts ...pulumi.ResourceOption) (*projects.IAMCustomRole, error) {
	roleId, err := random.NewRandomString(ctx, fmt.Sprintf("role-%s-id", name), &random.RandomStringArgs{
		Special: pulumi.Bool(false),
		Length:  pulumi.Int(8),
		Keepers: pulumi.ToMap(map[string]interface{}{
			"policy-name": name,
		}),
	}, opts...)
	if err != nil {
		return nil, err
	}

	return projects.NewIAMCustomRole(ctx, name, &projects.IAMCustomRoleArgs{
		Title:       pulumi.String(name),
		Permissions: pulumi.ToStringArray(actions),
		RoleId:      roleId.ID(),
	}, opts...)
}
