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
  value = azurerm_resource_group.resource_group.name
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

output "storage_account_blob_endpoint" {
  value = one(azurerm_storage_account.storage) != null ? one(azurerm_storage_account.storage).primary_blob_endpoint : null
}

output "storage_account_queue_endpoint" {
  value = one(azurerm_storage_account.storage) != null ? one(azurerm_storage_account.storage).primary_queue_endpoint : null
}

output "storage_account_web_host" {
  value = one(azurerm_storage_account.storage) != null ? one(azurerm_storage_account.storage).primary_web_host : null
}

output "registry_username" {
  value = azurerm_container_registry.container_registry.admin_username
}

output "registry_password" {
  value = azurerm_container_registry.container_registry.admin_password
}

output "infrastructure_subnet_id" {
  value = one(azurerm_subnet.database_infrastructure_subnet) != null ? one(azurerm_subnet.database_infrastructure_subnet).id : null
}

output "container_app_subnet_id" {
  value = one(azurerm_subnet.database_client_subnet) != null ? one(azurerm_subnet.database_client_subnet).id : null
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
