output "cdn_url" {
  description = "The URL of the CDN endpoint"
  value = "https://${azapi_resource.cdn_endpoint.output.properties.hostName}"
}