output "database_master_password" {
  value = random_password.database_master_password.result
}

output "database_server_fqdn" {
  value = azurerm_postgresql_flexible_server.database.fqdn
}

output "database_server_id" {
  value = azurerm_postgresql_flexible_server.database.id
}

output "infrastructure_subnet_id" {
  value = azurerm_subnet.database_infrastructure_subnet.id
}

output "container_app_subnet_id" {
  value = azurerm_subnet.database_client_subnet.id
}
