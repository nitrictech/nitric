package policy

import (
	"fmt"

	v1 "github.com/nitrictech/nitric/core/pkg/api/nitric/v1"
	"github.com/pulumi/pulumi-azure-native-sdk/authorization"
	"github.com/pulumi/pulumi-random/sdk/v4/go/random"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type StackRolesArgs struct{}

type StackRoles struct {
	pulumi.ResourceState

	Name               string
	ClientID           pulumi.StringOutput
	TenantID           pulumi.StringOutput
	ServicePrincipalId pulumi.StringOutput
	ClientSecret       pulumi.StringOutput
}

type RoleDefinition struct {
	Description      pulumi.StringInput
	Permissions      authorization.PermissionArray
	AssignableScopes pulumi.StringArray
}

var roleDefinitions = map[v1.Action]RoleDefinition{
	v1.Action_BucketFileGet: {
		Description: pulumi.String("bucket read access"),
		Permissions: authorization.PermissionArray{
			authorization.PermissionArgs{
				Actions: pulumi.StringArray{
					pulumi.String("Microsoft.Storage/storageAccounts/blobServices/containers/read"),
				},
				DataActions: pulumi.StringArray{
					pulumi.String("Microsoft.Storage/storageAccounts/blobServices/containers/blobs/read"),
				},
				NotActions: pulumi.StringArray{},
			},
		},
		AssignableScopes: pulumi.ToStringArray([]string{
			"/",
		}),
	},
	v1.Action_BucketFilePut: {
		Description: pulumi.String("bucket file write access"),
		Permissions: authorization.PermissionArray{
			authorization.PermissionArgs{
				Actions: pulumi.StringArray{},
				DataActions: pulumi.StringArray{
					pulumi.String("Microsoft.Storage/storageAccounts/blobServices/containers/blobs/write"),
				},
				NotActions: pulumi.StringArray{},
			},
		},
		AssignableScopes: pulumi.ToStringArray([]string{
			"/",
		}),
	},
	v1.Action_BucketFileList: {
		Description: pulumi.String("bucket file list access"),
		Permissions: authorization.PermissionArray{
			authorization.PermissionArgs{
				DataActions: pulumi.StringArray{
					pulumi.String("Microsoft.Storage/storageAccounts/blobServices/containers/blobs/read"),
				},
				NotActions: pulumi.StringArray{},
			},
		},
		AssignableScopes: pulumi.ToStringArray([]string{
			"/",
		}),
	},
	v1.Action_BucketFileDelete: {
		Description: pulumi.String("bucket file delete access"),
		Permissions: authorization.PermissionArray{
			authorization.PermissionArgs{
				Actions: pulumi.StringArray{},
				DataActions: pulumi.StringArray{
					pulumi.String("Microsoft.Storage/storageAccounts/blobServices/containers/blobs/delete"),
				},
				NotActions: pulumi.StringArray{},
			},
		},
		AssignableScopes: pulumi.ToStringArray([]string{
			"/",
		}),
	},
	v1.Action_QueueList: {
		Description: pulumi.String("queue list access"),
		Permissions: authorization.PermissionArray{
			authorization.PermissionArgs{
				Actions: pulumi.StringArray{
					pulumi.String("Microsoft.Storage/storageAccounts/queueServices/queues/read"),
				},
				DataActions: pulumi.StringArray{},
				NotActions:  pulumi.StringArray{},
			},
		},
		AssignableScopes: pulumi.ToStringArray([]string{
			"/",
		}),
	},
	v1.Action_QueueDetail: {
		Description: pulumi.String("queue detail access"),
		Permissions: authorization.PermissionArray{
			authorization.PermissionArgs{
				Actions: pulumi.StringArray{
					pulumi.String("Microsoft.Storage/storageAccounts/queueServices/queues/read"),
				},
				DataActions: pulumi.StringArray{},
				NotActions:  pulumi.StringArray{},
			},
		},
		AssignableScopes: pulumi.ToStringArray([]string{
			"/",
		}),
	},
	v1.Action_QueueSend: {
		Description: pulumi.String("queue send access"),
		Permissions: authorization.PermissionArray{
			authorization.PermissionArgs{
				Actions: pulumi.StringArray{},
				DataActions: pulumi.StringArray{
					pulumi.String("Microsoft.Storage/storageAccounts/queueServices/queues/messages/write"),
				},
				NotActions: pulumi.StringArray{},
			},
		},
		AssignableScopes: pulumi.ToStringArray([]string{
			"/",
		}),
	},
	v1.Action_QueueReceive: {
		Description: pulumi.String("queue receive access"),
		Permissions: authorization.PermissionArray{
			authorization.PermissionArgs{
				Actions: pulumi.StringArray{
					pulumi.String("Microsoft.Storage/storageAccounts/queueServices/queues/read"),
				},
				DataActions: pulumi.StringArray{
					pulumi.String("Microsoft.Storage/storageAccounts/queueServices/queues/messages/read"),
					pulumi.String("Microsoft.Storage/storageAccounts/queueServices/queues/messages/delete"),
					// pulumi.String("Microsoft.Storage/storageAccounts/queueServices/queues/messages/update"),
				},
				NotActions: pulumi.StringArray{},
			},
		},
		AssignableScopes: pulumi.ToStringArray([]string{
			"/",
		}),
	},
	v1.Action_TopicDetail: {
		Description: pulumi.String("topic detail access"),
		Permissions: authorization.PermissionArray{
			authorization.PermissionArgs{
				Actions: pulumi.StringArray{
					pulumi.String("Microsoft.EventGrid/topics/read"),
				},
				DataActions: pulumi.StringArray{},
				NotActions:  pulumi.StringArray{},
			},
		},
		AssignableScopes: pulumi.ToStringArray([]string{
			"/",
		}),
	},
	v1.Action_TopicList: {
		Description: pulumi.String("topic list access"),
		Permissions: authorization.PermissionArray{
			authorization.PermissionArgs{
				Actions: pulumi.StringArray{
					pulumi.String("Microsoft.EventGrid/topics/read"),
				},
				DataActions: pulumi.StringArray{},
				NotActions:  pulumi.StringArray{},
			},
		},
		AssignableScopes: pulumi.ToStringArray([]string{
			"/",
		}),
	},
	v1.Action_TopicEventPublish: {
		Description: pulumi.String("topic event publish access"),
		Permissions: authorization.PermissionArray{
			authorization.PermissionArgs{
				Actions: pulumi.StringArray{
					pulumi.String("Microsoft.EventGrid/topics/*/write"),
				},
				DataActions: pulumi.StringArray{
					pulumi.String("Microsoft.EventGrid/events/send/action"),
				},
				NotActions: pulumi.StringArray{},
			},
		},
		AssignableScopes: pulumi.ToStringArray([]string{
			"/",
		}),
	},
	v1.Action_SecretAccess: {
		Description: pulumi.String("keyvault secret read access"),
		Permissions: authorization.PermissionArray{
			authorization.PermissionArgs{
				Actions: pulumi.StringArray{},
				DataActions: pulumi.StringArray{
					pulumi.String("Microsoft.KeyVault/vaults/secrets/getSecret/action"),
				},
				NotActions: pulumi.StringArray{},
			},
		},
		AssignableScopes: pulumi.ToStringArray([]string{
			"/",
		}),
	},
	v1.Action_SecretPut: {
		Description: pulumi.String("keyvault secret write access"),
		Permissions: authorization.PermissionArray{
			authorization.PermissionArgs{
				Actions: pulumi.StringArray{
					pulumi.String("Microsoft.KeyVault/vaults/secrets/write"),
				},
				DataActions: pulumi.StringArray{
					pulumi.String("Microsoft.KeyVault/vaults/secrets/setSecret/action"),
				},
				NotActions: pulumi.StringArray{},
			},
		},
		AssignableScopes: pulumi.ToStringArray([]string{
			"/",
		}),
	},
}

type Roles struct {
	pulumi.ResourceState

	Name            string
	RoleDefinitions map[v1.Action]*authorization.RoleDefinition
}

var actionNames = map[v1.Action]string{
	v1.Action_BucketFileGet:     "BucketFileGet",
	v1.Action_BucketFilePut:     "BucketFilePut",
	v1.Action_BucketFileDelete:  "BucketFileDelete",
	v1.Action_BucketFileList:    "BucketFileList",
	v1.Action_QueueList:         "QueueList",
	v1.Action_QueueDetail:       "QueueDetail",
	v1.Action_QueueSend:         "QueueSend",
	v1.Action_QueueReceive:      "QueueReceive",
	v1.Action_TopicDetail:       "TopicDetail",
	v1.Action_TopicEventPublish: "TopicPublish",
	v1.Action_TopicList:         "TopicList",
	v1.Action_SecretAccess:      "SecretAccess",
	v1.Action_SecretPut:         "SecretPut",
}

func CreateRoles(ctx *pulumi.Context, stackId string, subscriptionId string, rgName pulumi.StringInput) (*Roles, error) {
	res := &Roles{Name: "nitric-roles", RoleDefinitions: map[v1.Action]*authorization.RoleDefinition{}}

	err := ctx.RegisterComponentResource("nitric:roles:AzureADRoles", "nitric-roles", res)
	if err != nil {
		return nil, err
	}

	for id, roleDef := range roleDefinitions {
		name := actionNames[id]

		roleGuid, err := random.NewRandomUuid(ctx, name, &random.RandomUuidArgs{
			Keepers: pulumi.ToMap(map[string]interface{}{
				"subscriptionId": subscriptionId,
			}),
		}, pulumi.Parent(res))
		if err != nil {
			return nil, err
		}

		roleName := fmt.Sprintf("%s-%s", stackId, name)

		createdRole, err := authorization.NewRoleDefinition(ctx, name, &authorization.RoleDefinitionArgs{
			RoleDefinitionId: roleGuid.Result,
			RoleName:         pulumi.String(roleName),
			Description:      roleDef.Description,
			Permissions:      roleDef.Permissions,
			Scope:            pulumi.Sprintf("/subscriptions/%s/resourceGroups/%s", subscriptionId, rgName),
			AssignableScopes: pulumi.StringArray{
				pulumi.Sprintf("/subscriptions/%s/resourceGroups/%s", subscriptionId, rgName),
			},
		}, pulumi.Parent(res))
		if err != nil {
			return nil, err
		}

		res.RoleDefinitions[id] = createdRole
	}

	return res, nil
}
