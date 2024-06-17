
resource "google_api_gateway_api" "proxy_api" {
  provider = google-beta
  api_id   = var.name
  labels = {
    "x-nitric-${var.stack_id}-name" = var.name
    "x-nitric-${var.stack_id}-type" = "http-proxy"
  }
}

data "template_file" "openapi_spec" {
  template = file("${path.module}/openapi_template.json")
  vars = {
    name         = var.name
    target_service_url = var.target_service_url
  }
}

resource "google_api_gateway_api_config" "api_config" {
  provider      = google-beta
  api           = google_api_gateway_api.proxy_api.api_id
  api_config_id = "${var.name}-config"

  openapi_documents {
    document {
      path     = "openapi.json"
      contents = base64encode(data.template_file.openapi_spec.rendered)
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
  gateway_id   = "${var.name}-gateway"
  api_config   = google_api_gateway_api_config.api_config.id

  labels = {
    "x-nitric-${var.stack_id}-name" = var.name
    "x-nitric-${var.stack_id}-type" = "http-proxy"
  }
}
