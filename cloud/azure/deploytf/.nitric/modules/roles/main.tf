
locals {
  role_definitions = {
    "KeyValueStoreRead" = {
        description = "keyvalue read access"
        permissions = {
            actions = []
            data_actions = [
                "Microsoft.Storage/storageAccounts/tableServices/tables/entities/read"
            ]
            not_actions = []
        }
    },
    "KeyValueStoreWrite" = {
        description = "keyvalue write access"
        permissions = {
            actions = []
            data_actions = [
                "Microsoft.Storage/storageAccounts/tableServices/tables/entities/write",
                "Microsoft.Storage/storageAccounts/tableServices/tables/entities/delete"
            ]
            not_actions = []
        }
    },
    "KeyValueStoreDelete" = {
        description = "keyvalue delete access"
        permissions = {
            actions = []
            data_actions = [
                "Microsoft.Storage/storageAccounts/tableServices/tables/entities/delete"
            ]
            not_actions = []
        }
    },
    "QueueEnqueue" = {
        description = "queue enqueue access"
        permissions = {
            actions = [
                "Microsoft.Storage/storageAccounts/queueServices/queues/read"
            ]
            data_actions = [
                "Microsoft.Storage/storageAccounts/queueServices/queues/messages/write"
            ]
            not_actions = []
        }
    },
    "QueueDequeue" = {
        description = "queue dequeue access"
        permissions = {
            actions = [
                "Microsoft.Storage/storageAccounts/queueServices/queues/read"
            ]
            data_actions = [
                "Microsoft.Storage/storageAccounts/queueServices/queues/messages/read",
                "Microsoft.Storage/storageAccounts/queueServices/queues/messages/delete"
            ]
            not_actions = []
        },
    },
    "BucketFileGet" = {
        description = "bucket file get"
        permissions = {
            actions = [
                "Microsoft.Storage/storageAccounts/blobServices/containers/read"
            ]
            data_actions = [
                "Microsoft.Storage/storageAccounts/blobServices/containers/blobs/read"
            ]
            not_actions = []
        }
    },
    "BucketFilePut" = {
        description = "bucket file put"
        permissions = {
            actions = []
            data_actions = [
                "Microsoft.Storage/storageAccounts/blobServices/containers/blobs/write"
            ]
            not_actions = []
        }
    },
    "BucketFileDelete" = {
        description = "bucket file delete"
        permissions = {
            actions = []
            data_actions = [
                "Microsoft.Storage/storageAccounts/blobServices/containers/blobs/delete"
            ]
            not_actions = []
        }
    },
    "BucketFileList" = {
        description = "bucket file list"
        permissions = {
            actions = []
            data_actions = [
                "Microsoft.Storage/storageAccounts/blobServices/containers/blobs/read"
            ]
            not_actions = []
        }
    },
    "TopicPublish" = {
        description = "topic publish"
        permissions = {
            actions = [
                "Microsoft.EventGrid/topics/read",
                "Microsoft.EventGrid/topics/*/write"
            ]
            data_actions = [
                "Microsoft.EventGrid/events/send/action"
            ]
            not_actions = []
        }
    },
    "SecretAccess" = {
        description = "secret access"
        permissions = {
            actions = []
            data_actions = [
                "Microsoft.KeyVault/vaults/secrets/getSecret/action"
            ]
            not_actions = []
        }
    },
    "SecretPut" = {
        description = "secret put"
        permissions = {
            actions = []
            data_actions = [
                "Microsoft.KeyVault/vaults/secrets/setSecret/action"
            ]
            not_actions = []
        }
    },
  }
}

data "azurerm_subscription" "current" {}

resource "azurerm_role_definition" "nitric" {
  for_each = { for role in local.role_definitions : role.description => role }

  description = each.value.description
  name               = each.key
  scope              = "/subscriptions/${data.azurerm_subscription.current.subscription_id}/resourceGroups/${var.resource_group_name}"

  permissions {
    actions     = each.value.permissions.actions
    data_actions = each.value.permissions.data_actions
    not_actions = each.value.permissions.not_actions
  }

  assignable_scopes = "/subscriptions/${data.azurerm_subscription.current.subscription_id}/resourceGroups/${var.resource_group_name}"
}