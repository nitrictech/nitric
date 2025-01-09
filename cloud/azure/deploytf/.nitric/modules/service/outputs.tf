output "container_app_id" {
    description = "The id of the container app"
    value       = azurerm_container_app.container_app.id
}

output "dapr_app_id" {
    description = "The id of the dapr app"
    value       = azurerm_container_app.container_app.dapr_app_id
}

output "fqdn" {
    description = "The endpoint of the container app"
    value       = azurerm_container_app.container_app.latest_revision_fqdn
}

output "event_token" {
    description = "The event token of the container app"
    value       = random_string.event_token.result
}

output "client_id" {
    description = "The client id of the container app"
    value       = azurerm_container_app.container_app.client_id
}

output "service_principal_id" {
    description = "The service principal id of the container app"
    value       = azuread_service_principal.service_identity.id
}

output "tenant_id" {
    description = "The tenant id of the container app"
    value       = azurerm_container_app.container_app.tenant_id
}

output "endpoint" {
    description = "The endpoint of the container app"
    value       = "https://${azurerm_container_app.container_app.latest_revision_fqdn}"
}