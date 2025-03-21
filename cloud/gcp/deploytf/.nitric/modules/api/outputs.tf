output "endpoint" {
  value = "https://${google_api_gateway_gateway.gateway.default_hostname}"
}

output "region" {
  value = var.region
}

output gateway_id {
  value = google_api_gateway_gateway.gateway.id
}

output default_host {
  value = google_api_gateway_gateway.gateway.default_hostname
}
