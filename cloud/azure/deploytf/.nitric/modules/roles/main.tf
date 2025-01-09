data "azurerm_subscription" "current" {}

resource "azurerm_role_definition" "nitric_role_kv_read" {
  description = "keyvalue read access"
  name        = "KeyValueStoreRead"
  scope       = "/subscriptions/${data.azurerm_subscription.current.subscription_id}/resourceGroups/${var.resource_group_name}"

  permissions {
    actions = []
    data_actions = [
      "Microsoft.Storage/storageAccounts/tableServices/tables/entities/read"
    ]
    not_actions = []
  }

  assignable_scopes = "/subscriptions/${data.azurerm_subscription.current.subscription_id}/resourceGroups/${var.resource_group_name}"
}

resource "azurerm_role_definition" "nitric_role_kv_write" {
  description = "nitric keyvalue write access"
  name        = "KeyValueStoreWrite"
  scope       = "/subscriptions/${data.azurerm_subscription.current.subscription_id}/resourceGroups/${var.resource_group_name}"

  permissions {
    actions = []
    data_actions = [
      "Microsoft.Storage/storageAccounts/tableServices/tables/entities/write",
      "Microsoft.Storage/storageAccounts/tableServices/tables/entities/delete"
    ]
    not_actions = []
  }

  assignable_scopes = "/subscriptions/${data.azurerm_subscription.current.subscription_id}/resourceGroups/${var.resource_group_name}"
}

resource "azurerm_role_definition" "nitric_role_kv_delete" {
  description = "nitric keyvalue delete access"
  name        = "KeyValueStoreDelete"
  scope       = "/subscriptions/${data.azurerm_subscription.current.subscription_id}/resourceGroups/${var.resource_group_name}"

  permissions {
    actions = []
    data_actions = [
      "Microsoft.Storage/storageAccounts/tableServices/tables/entities/delete"
    ]
    not_actions = []
  }

  assignable_scopes = "/subscriptions/${data.azurerm_subscription.current.subscription_id}/resourceGroups/${var.resource_group_name}"
}

resource "azurerm_role_definition" "nitric_role_queue_enqueue" {
  description = "nitric queue enqueue access"
  name        = "QueueEnqueue"
  scope       = "/subscriptions/${data.azurerm_subscription.current.subscription_id}/resourceGroups/${var.resource_group_name}"

  permissions {
    actions = [
      "Microsoft.Storage/storageAccounts/queueServices/queues/read"
    ]
    data_actions = [
      "Microsoft.Storage/storageAccounts/queueServices/queues/messages/write"
    ]
    not_actions = []
  }

  assignable_scopes = "/subscriptions/${data.azurerm_subscription.current.subscription_id}/resourceGroups/${var.resource_group_name}"
}
resource "azurerm_role_definition" "nitric_role_queue_dequeue" {
  description = "nitric queue dequeue access"
  name        = "QueueDequeue"
  scope       = "/subscriptions/${data.azurerm_subscription.current.subscription_id}/resourceGroups/${var.resource_group_name}"

  permissions {
    actions = [
      "Microsoft.Storage/storageAccounts/queueServices/queues/read"
    ]
    data_actions = [
      "Microsoft.Storage/storageAccounts/queueServices/queues/messages/read",
      "Microsoft.Storage/storageAccounts/queueServices/queues/messages/delete"
    ]
    not_actions = []
  }

  assignable_scopes = "/subscriptions/${data.azurerm_subscription.current.subscription_id}/resourceGroups/${var.resource_group_name}"
}

resource "azurerm_role_definition" "nitric_role_bucket_file_get" {
  description = "nitric bucket file get access"
  name        = "BucketFileGet"
  scope       = "/subscriptions/${data.azurerm_subscription.current.subscription_id}/resourceGroups/${var.resource_group_name}"

  permissions {
    actions = [
      "Microsoft.Storage/storageAccounts/blobServices/containers/read"
    ]
    data_actions = [
      "Microsoft.Storage/storageAccounts/blobServices/containers/blobs/read"
    ]
    not_actions = []
  }

  assignable_scopes = "/subscriptions/${data.azurerm_subscription.current.subscription_id}/resourceGroups/${var.resource_group_name}"
}

resource "azurerm_role_definition" "nitric_role_bucket_file_put" {
  description = "nitric bucket file put access"
  name        = "BucketFilePut"
  scope       = "/subscriptions/${data.azurerm_subscription.current.subscription_id}/resourceGroups/${var.resource_group_name}"

  permissions {
    actions = []
    data_actions = [
      "Microsoft.Storage/storageAccounts/blobServices/containers/blobs/write"
    ]
    not_actions = []
  }

  assignable_scopes = "/subscriptions/${data.azurerm_subscription.current.subscription_id}/resourceGroups/${var.resource_group_name}"
}

resource "azurerm_role_definition" "nitric_role_bucket_file_delete" {
  description = "nitric bucket file delete access"
  name        = "BucketFileDelete"
  scope       = "/subscriptions/${data.azurerm_subscription.current.subscription_id}/resourceGroups/${var.resource_group_name}"

  permissions {
    actions = []
    data_actions = [
      "Microsoft.Storage/storageAccounts/blobServices/containers/blobs/delete"
    ]
    not_actions = []
  }

  assignable_scopes = "/subscriptions/${data.azurerm_subscription.current.subscription_id}/resourceGroups/${var.resource_group_name}"
}

resource "azurerm_role_definition" "nitric_role_bucket_file_list" {
  description = "nitric bucket file list access"
  name        = "BucketFileList"
  scope       = "/subscriptions/${data.azurerm_subscription.current.subscription_id}/resourceGroups/${var.resource_group_name}"

  permissions {
    actions = []
    data_actions = [
      "Microsoft.Storage/storageAccounts/blobServices/containers/blobs/read"
    ]
    not_actions = []
  }

  assignable_scopes = "/subscriptions/${data.azurerm_subscription.current.subscription_id}/resourceGroups/${var.resource_group_name}"
}

resource "azurerm_role_definition" "nitric_role_topic_publish" {
  description = "nitric topic publish access"
  name        = "TopicPublish"
  scope       = "/subscriptions/${data.azurerm_subscription.current.subscription_id}/resourceGroups/${var.resource_group_name}"

  permissions {
    actions = [
      "Microsoft.EventGrid/topics/read",
      "Microsoft.EventGrid/topics/*/write"
    ]
    data_actions = [
      "Microsoft.EventGrid/events/send/action"
    ]
    not_actions = []
  }

  assignable_scopes = "/subscriptions/${data.azurerm_subscription.current.subscription_id}/resourceGroups/${var.resource_group_name}"
}

resource "azurerm_role_definition" "nitric_role_secret_access" {
  description = "nitric secret access access"
  name        = "SecretAccess"
  scope       = "/subscriptions/${data.azurerm_subscription.current.subscription_id}/resourceGroups/${var.resource_group_name}"

  permissions {
    actions = []
    data_actions = [
      "Microsoft.KeyVault/vaults/secrets/getSecret/action"
    ]
    not_actions = []
  }

  assignable_scopes = "/subscriptions/${data.azurerm_subscription.current.subscription_id}/resourceGroups/${var.resource_group_name}"
}

resource "azurerm_role_definition" "nitric_role_secret_put" {
  description = "nitric secret put access"
  name        = "SecretPut"
  scope       = "/subscriptions/${data.azurerm_subscription.current.subscription_id}/resourceGroups/${var.resource_group_name}"

  permissions {
    actions = []
    data_actions = [
      "Microsoft.KeyVault/vaults/secrets/setSecret/action"
    ]
    not_actions = []
  }

  assignable_scopes = "/subscriptions/${data.azurerm_subscription.current.subscription_id}/resourceGroups/${var.resource_group_name}"
}
