output "roles" {
  value = {
    "KeyValueStoreRead" = azurerm_role_definition.nitric_role_kv_read.id
    "KeyValueStoreWrite" = azurerm_role_definition.nitric_role_kv_write.id
    "KeyValueStoreDelete" = azurerm_role_definition.nitric_role_kv_delete.id
    "QueueEnqueue" = azurerm_role_definition.nitric_role_queue_enqueue.id
    "QueueDequeue" = azurerm_role_definition.nitric_role_queue_dequeue.id
    "BucketFileGet" = azurerm_role_definition.nitric_role_bucket_file_get.id
    "BucketFilePut" = azurerm_role_definition.nitric_role_bucket_file_put.id
    "BucketFileDelete" = azurerm_role_definition.nitric_role_bucket_file_delete.id
    "BucketFileList" = azurerm_role_definition.nitric_role_bucket_file_list.id
    "TopicPublish" = azurerm_role_definition.nitric_role_topic_publish.id
    "SecretAccess" = azurerm_role_definition.nitric_role_secret_access.id
    "SecretPut" = azurerm_role_definition.nitric_role_secret_put.id
  }
}