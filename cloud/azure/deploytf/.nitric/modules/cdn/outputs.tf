output "cdn_url" {
  description = "The URL of the endpoint"
  value = "https://${azurerm_cdn_frontdoor_endpoint.cdn_endpoint.host_name}"
}