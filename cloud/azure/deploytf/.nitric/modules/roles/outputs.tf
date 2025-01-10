output "bucket_read" {
  value = azurerm_role_definition.nitric_role_bucket_file_get.id
  description = "The role ID for the Nitric bucket read role"
}

output "topic_publish" {
  value = azurerm_role_definition.nitric_role_topic_publish.id
  description = "The role ID for the Nitric topic publish role"
}

output "bucket_write" {
  value = azurerm_role_definition.nitric_role_bucket_file_put.id
  description = "The role ID for the Nitric bucket write role"
}

output "bucket_delete" {
  value = azurerm_role_definition.nitric_role_bucket_file_delete.id
  description = "The role ID for the Nitric bucket delete role"
}

output "bucket_list" {
  value = azurerm_role_definition.nitric_role_bucket_file_list.id
  description = "The role ID for the Nitric bucket list role"
}

output "secret_access" {
  value       = azurerm_role_definition.nitric_role_secret_access.id
  description = "The role ID for the Nitric secrete access role"
}

output "secret_put" {
  value       = azurerm_role_definition.nitric_role_secret_put.id
  description = "The role ID for the Nitric secrete put role"
}

output "kv_read" {
  value       = azurerm_role_definition.nitric_role_kv_read.id
  description = "The role ID for the Nitric kv read role"
}

output "kv_write" {
  value       = azurerm_role_definition.nitric_role_kv_write.id
  description = "The role ID for the Nitric kv write role"
}

output "kv_delete" {
  value       = azurerm_role_definition.nitric_role_kv_delete.id
  description = "The role ID for the Nitric kv write role"
}

output "queue_enqueue" {
  value       = azurerm_role_definition.nitric_role_queue_enqueue.id
  description = "The role ID for the Nitric queue enqueue role"
}

output "queue_dequeue" {
  value       = azurerm_role_definition.nitric_role_queue_dequeue.id
  description = "The role ID for the Nitric queue dequeue role"
}