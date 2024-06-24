
resource "google_api_gateway_api" "proxy_api" {
  provider = google-beta
  api_id   = replace(var.name, "_", "-")
  labels = {
    "x-nitric-${var.stack_id}-name" = var.name
    "x-nitric-${var.stack_id}-type" = "http-proxy"
  }
}

resource "google_api_gateway_api_config" "api_config" {
  provider      = google-beta
  api           = google_api_gateway_api.proxy_api.api_id
  api_config_id = "${replace(var.name, "_", "-")}-config"

  openapi_documents {
    document {
      path = "openapi.json"
      contents = base64encode(templatefile("${path.module}/openapi_template.json", {
        name               = var.name
        target_service_url = var.target_service_url
      }))
    }
  }
  gateway_config {
    backend_config {
      google_service_account = var.invoker_email
    }
  }

  labels = {
    "x-nitric-${var.stack_id}-name" = var.name
    "x-nitric-${var.stack_id}-type" = "http-proxy"
  }
}

resource "google_api_gateway_gateway" "gateway" {
  provider     = google-beta
  gateway_id   = "${replace(var.name, "_", "-")}-gateway"
  api_config   = google_api_gateway_api_config.api_config.id

  labels = {
    "x-nitric-${var.stack_id}-name" = var.name
    "x-nitric-${var.stack_id}-type" = "http-proxy"
  }
}
