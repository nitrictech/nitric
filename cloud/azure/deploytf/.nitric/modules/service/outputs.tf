output "container_app_id" {
  description = "The id of the container app"
  value       = azurerm_container_app.container_app.id
}

output "dapr_app_id" {
  description = "The id of the dapr app"
  value       = local.container_app_name
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
  value       = azuread_service_principal.service_identity.client_id
}

output "client_secret" {
  description = "The client secret of the container app"
  value       = azuread_service_principal_password.service_identity.value
}
output "service_principal_id" {
  description = "The service principal id of the container app"
  value       = azuread_service_principal.service_identity.object_id
}

output "tenant_id" {
  description = "The tenant id of the container app"
  value       = data.azurerm_client_config.current.tenant_id
}

output "endpoint" {
  description = "The endpoint of the container app"
  # Gets the persistent endpoint of the container app, not the latest revision
  # See: https://github.com/hashicorp/terraform-provider-azurerm/issues/26559
  value = "https://${azurerm_container_app.container_app.ingress[0].fqdn}"
}

output "poll_url" {
  description = "The URL polling task to ensure the lastest revision is available"
  value       = null_resource.poll_url
}
