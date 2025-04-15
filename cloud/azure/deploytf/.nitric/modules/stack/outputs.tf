output "stack_name" {
  value = local.stack_name
}

output "stack_id" {
  value = random_string.stack_id.result
}

output "subscription_id" {
  value = data.azurerm_client_config.current.subscription_id
}

output "keyvault_name" {
  value = one(azurerm_key_vault.keyvault) != null ? one(azurerm_key_vault.keyvault).name : null
}

output "resource_group_name" {
  value = local.resource_group_name
}

output "app_identity" {
  value = azurerm_user_assigned_identity.managed_identity.id
}

output "app_identity_client_id" {
  value = azurerm_user_assigned_identity.managed_identity.client_id
}

output "container_app_environment_id" {
  value = azurerm_container_app_environment.environment.id
}

output "registry_login_server" {
  value = azurerm_container_registry.container_registry.login_server
}

output "storage_account_name" {
  value = one(azurerm_storage_account.storage) != null ? one(azurerm_storage_account.storage).name : null
}

output "storage_account_id" {
  value = one(azurerm_storage_account.storage) != null ? one(azurerm_storage_account.storage).id : null
}

output "storage_account_connection_string" {
  value = one(azurerm_storage_account.storage) != null ? one(azurerm_storage_account.storage).primary_connection_string : null
}

output "storage_account_blob_endpoint" {
  value = one(azurerm_storage_account.storage) != null ? one(azurerm_storage_account.storage).primary_blob_endpoint : null
}

output "storage_account_queue_endpoint" {
  value = one(azurerm_storage_account.storage) != null ? one(azurerm_storage_account.storage).primary_queue_endpoint : null
}

output "storage_account_table_endpoint" {
  value = one(azurerm_storage_account.storage) != null ? one(azurerm_storage_account.storage).primary_table_endpoint : null
}

output "registry_username" {
  value = azurerm_container_registry.container_registry.admin_username
}

output "registry_password" {
  value = azurerm_container_registry.container_registry.admin_password
}

output "subnet_id" {
  value = local.subnet_id
}

output "database_master_password" {
  value = one(random_password.database_master_password) != null ? one(random_password.database_master_password).result : null
}

output "database_server_id" {
  value = one(azurerm_postgresql_flexible_server.database) != null ? one(azurerm_postgresql_flexible_server.database).id : null
}

output "database_server_fqdn" {
  value = one(azurerm_postgresql_flexible_server.database) != null ? one(azurerm_postgresql_flexible_server.database).fqdn : null
}

output "database_server_name" {
  value = one(azurerm_postgresql_flexible_server.database) != null ? one(azurerm_postgresql_flexible_server.database).name : null
}
