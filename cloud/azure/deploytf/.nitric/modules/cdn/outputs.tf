# TODO - add outputs for the CDN endpoint
# output "cdn_url" {
#   description = "The URL of the CDN endpoint"
#   value = "https://${jsondecode(azapi_resource.cdn_endpoint.output).properties.hostName}"
# }