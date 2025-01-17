output "bucket_read" {
  value = azurerm_role_definition.nitric_role_bucket_file_get.role_definition_resource_id
  description = "The role ID for the Nitric bucket read role"
}

output "topic_publish" {
  value = azurerm_role_definition.nitric_role_topic_publish.role_definition_resource_id
  description = "The role ID for the Nitric topic publish role"
}

output "bucket_write" {
  value = azurerm_role_definition.nitric_role_bucket_file_put.role_definition_resource_id
  description = "The role ID for the Nitric bucket write role"
}

output "bucket_delete" {
  value = azurerm_role_definition.nitric_role_bucket_file_delete.role_definition_resource_id
  description = "The role ID for the Nitric bucket delete role"
}

output "bucket_list" {
  value = azurerm_role_definition.nitric_role_bucket_file_list.role_definition_resource_id
  description = "The role ID for the Nitric bucket list role"
}

output "secret_access" {
  value       = azurerm_role_definition.nitric_role_secret_access.role_definition_resource_id
  description = "The role ID for the Nitric secrete access role"
}

output "secret_put" {
  value       = azurerm_role_definition.nitric_role_secret_put.role_definition_resource_id
  description = "The role ID for the Nitric secrete put role"
}

output "kv_read" {
  value       = azurerm_role_definition.nitric_role_kv_read.role_definition_resource_id
  description = "The role ID for the Nitric kv read role"
}

output "kv_write" {
  value       = azurerm_role_definition.nitric_role_kv_write.role_definition_resource_id
  description = "The role ID for the Nitric kv write role"
}

output "kv_delete" {
  value       = azurerm_role_definition.nitric_role_kv_delete.role_definition_resource_id
  description = "The role ID for the Nitric kv write role"
}

output "queue_enqueue" {
  value       = azurerm_role_definition.nitric_role_queue_enqueue.role_definition_resource_id
  description = "The role ID for the Nitric queue enqueue role"
}

output "queue_dequeue" {
  value       = azurerm_role_definition.nitric_role_queue_dequeue.role_definition_resource_id
  description = "The role ID for the Nitric queue dequeue role"
}