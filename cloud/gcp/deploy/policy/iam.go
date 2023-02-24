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

package policy

import (
	"fmt"

	"github.com/nitrictech/nitric/cloud/gcp/deploy/bucket"
	"github.com/nitrictech/nitric/cloud/gcp/deploy/events"
	"github.com/nitrictech/nitric/cloud/gcp/deploy/queue"
	"github.com/nitrictech/nitric/cloud/gcp/deploy/secret"
	deploy "github.com/nitrictech/nitric/core/pkg/api/nitric/deploy/v1"
	v1 "github.com/nitrictech/nitric/core/pkg/api/nitric/v1"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/projects"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/pubsub"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/secretmanager"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/serviceaccount"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/storage"
	"github.com/pulumi/pulumi-random/sdk/v4/go/random"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Policy struct {
	pulumi.ResourceState

	Name         string
	RolePolicies []*projects.IAMMember
}

type StackResources struct {
	Topics        map[string]*events.PubSubTopic
	Queues        map[string]*queue.PubSubTopic
	Subscriptions map[string]*pubsub.Subscription
	Buckets       map[string]*bucket.CloudStorageBucket
	Secrets       map[string]*secret.SecretManagerSecret
}

type PrincipalMap = map[v1.ResourceType]map[string]*serviceaccount.Account

type PolicyArgs struct {
	Policy *deploy.Policy
	// Resources in the stack that must be protected
	Resources *StackResources
	// Resources in the stack that may act as actors
	Principals PrincipalMap

	ProjectID pulumi.StringInput

	StackID pulumi.StringInput
}

var gcpListActions []string = []string{
	"pubsub.topics.list",
	"pubsub.snapshots.list",
	"resourcemanager.projects.get",
	"secretmanager.secrets.list",
	"apigateway.gateways.list",
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
	v1.Action_QueueSend: {
		"pubsub.topics.get",
		"pubsub.topics.publish",
	},
	v1.Action_QueueReceive: {
		"pubsub.topics.get",
		"pubsub.topics.attachSubscription",
		"pubsub.snapshots.seek",
		"pubsub.subscriptions.consume",
	},
	v1.Action_QueueDetail: {
		"pubsub.topics.get",
	},
	v1.Action_QueueList: {}, // see above in gcpListActions
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

func NewIAMPolicy(ctx *pulumi.Context, name string, args *PolicyArgs, opts ...pulumi.ResourceOption) (*Policy, error) {
	res := &Policy{Name: name, RolePolicies: make([]*projects.IAMMember, 0)}

	err := ctx.RegisterComponentResource("nitric:func:GCPPolicy", name, res, opts...)
	if err != nil {
		return nil, err
	}

	actions := actionsToGcpActions(args.Policy.Actions)

	rolePolicy, err := newCustomRole(ctx, name, actions, pulumi.Parent(res))
	if err != nil {
		return nil, err
	}

	// for project level listings
	listRolePolicy, err := newCustomRole(ctx, name+"-list", gcpListActions, pulumi.Parent(res))
	if err != nil {
		return nil, err
	}

	for _, principal := range args.Policy.Principals {
		sa := args.Principals[v1.ResourceType_Function][principal.Name]

		for _, resource := range args.Policy.Resources {
			memberName := fmt.Sprintf("%s-%s", principal.Name, resource.Name)
			memberId := pulumi.Sprintf("serviceAccount:%s", sa.Email)

			// for project level listings
			listRolePolicy.Title.ToStringOutput().ApplyT(func(id string) (string, error) {
				_, err = projects.NewIAMMember(ctx, id+"-member", &projects.IAMMemberArgs{
					Member:  memberId,
					Project: args.ProjectID,
					Role:    listRolePolicy.Name,
				}, pulumi.Parent(res))

				return "", err
			})

			switch resource.Type {
			case v1.ResourceType_Bucket:
				b := args.Resources.Buckets[resource.Name]

				_, err = storage.NewBucketIAMMember(ctx, memberName, &storage.BucketIAMMemberArgs{
					Bucket: b.CloudStorage.Name,
					Member: memberId,
					Role:   rolePolicy.Name,
				}, pulumi.Parent(res))
				if err != nil {
					return nil, err
				}

			case v1.ResourceType_Collection:
				collActions := filterCollectionActions(actions)

				collRole, err := newCustomRole(ctx, memberName+"-role", collActions, pulumi.Parent(res))
				if err != nil {
					return nil, err
				}

				_, err = projects.NewIAMMember(ctx, memberName, &projects.IAMMemberArgs{
					Member:  memberId,
					Project: args.ProjectID,
					Role:    collRole.Name,
				}, pulumi.Parent(res))
				if err != nil {
					return nil, err
				}

			case v1.ResourceType_Queue:
				q := args.Resources.Queues[resource.Name]
				s := args.Resources.Subscriptions[resource.Name] // the subscription and topic have the same name

				_, err = pubsub.NewTopicIAMMember(ctx, memberName, &pubsub.TopicIAMMemberArgs{
					Topic:  q.PubSub.Name,
					Member: memberId,
					Role:   rolePolicy.Name,
				}, pulumi.Parent(res))
				if err != nil {
					return nil, err
				}

				needSubConsume := false

				for _, act := range args.Policy.Actions {
					if act == v1.Action_QueueReceive {
						needSubConsume = true
						break
					}
				}

				if needSubConsume {
					subRolePolicy, err := newCustomRole(ctx, name+"subscription", []string{"pubsub.subscriptions.consume"}, pulumi.Parent(res))
					if err != nil {
						return nil, err
					}

					_, err = pubsub.NewSubscriptionIAMMember(ctx, memberName, &pubsub.SubscriptionIAMMemberArgs{
						Subscription: s.Name,
						Member:       memberId,
						Role:         subRolePolicy.Name,
					}, pulumi.Parent(res))
					if err != nil {
						return nil, err
					}
				}

			case v1.ResourceType_Topic:
				t := args.Resources.Topics[resource.Name]

				_, err = pubsub.NewTopicIAMMember(ctx, memberName, &pubsub.TopicIAMMemberArgs{
					Topic:  t.PubSub.Name,
					Member: memberId,
					Role:   rolePolicy.Name,
				}, pulumi.Parent(res))
				if err != nil {
					return nil, err
				}

			case v1.ResourceType_Secret:
				s := args.Resources.Secrets[resource.Name]

				_, err = secretmanager.NewSecretIamMember(ctx, memberName, &secretmanager.SecretIamMemberArgs{
					SecretId: s.Secret.SecretId,
					Member:   memberId,
					Role:     rolePolicy.Name,
				}, pulumi.Parent(res))
				if err != nil {
					return nil, err
				}
			}
		}
	}

	return res, nil
}

func newCustomRole(ctx *pulumi.Context, name string, actions []string, opts ...pulumi.ResourceOption) (*projects.IAMCustomRole, error) {
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
