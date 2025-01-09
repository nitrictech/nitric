output "stack_name" {
  value = local.stack_name
}

output "resource_group_name" {
  value = azurerm_resource_group.resource_group.name
}

output "app_identity" {
  value = azurerm_user_assigned_identity.app_identity.name
}

output "container_app_environment_id" {
  value = azurerm_container_app_environment.environment.id
}

output "registry_login_server" {
  value = azurerm_container_registry.container_registry.login_server
}

output "storage_account_name" {
  value = azurerm_storage_account.storage_account.name
}

output "storage_account_id" {
  value = azurerm_storage_account.storage_account.id
}

output "registry_username" {
  value = azurerm_container_registry.container_registry.admin_username
}

output "registry_password" {
  value = azurerm_container_registry.container_registry.admin_password
}