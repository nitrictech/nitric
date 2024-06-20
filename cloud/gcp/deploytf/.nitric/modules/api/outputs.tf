output "endpoint" {
  value = "https://${google_api_gateway_gateway.gateway.default_hostname}"
}