output "container_app_id" {
    description = "The id of the container app"
    value       = azurerm_container_app.container_app.id
}

output "fqdn" {
    description = "The endpoint of the container app"
    value       = azurerm_container_app.container_app.latest_revision_fqdn
}