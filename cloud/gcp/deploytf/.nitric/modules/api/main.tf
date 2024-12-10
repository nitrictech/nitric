resource "google_api_gateway_api" "api" {
  provider = google-beta
  api_id   = "${var.name}-${random_string.api_config_id.result}"
  labels = {
    "x-nitric-${var.stack_id}-name" = var.name
    "x-nitric-${var.stack_id}-type" = "api"
  }
}

# Create a random id that changes when the openapi_spec changes
# required to prevent name collisions on update
resource "random_string" "api_config_id" {
  length = 4
  special = false
  upper   = false

  keepers = {
    spec: var.openapi_spec
  }
}

resource "google_service_account" "service_account" {
  provider   = google-beta
  account_id = "${var.name}-api-${random_string.api_config_id.result}"
}

# Get the current service account email and ensure it has service account user role on the service account
data "google_client_openid_userinfo" "deployer" {
}

locals {
  deployer_email = data.google_client_openid_userinfo.deployer.email
  deployer_type  = endswith(local.deployer_email, "gserviceaccount.com") ? "serviceAccount" : "user"
}

# If we're impersonation a service account, we need to grant that account the service account user role on the service account
resource "google_service_account_iam_member" "service_account_iam_member" {
  for_each = var.target_services

  service_account_id = google_service_account.service_account.name
  role               = "roles/iam.serviceAccountUser"
  member             = "${local.deployer_type}:${local.deployer_email}"
}


resource "google_cloud_run_service_iam_member" "member" {
  for_each = var.target_services

  service = each.value
  role    = "roles/run.invoker"
  member  = "serviceAccount:${google_service_account.service_account.email}"
}

resource "google_api_gateway_api_config" "api_config" {
  provider      = google-beta
  api           = google_api_gateway_api.api.api_id
  # A random ID, along with 'create_before_destroy' lifecycle, is used to prevent errors when updating (replacing) the API config
  api_config_id = "${var.name}-conf-${random_string.api_config_id.result}"

  openapi_documents {
    document {
      path     = "openapi.json"
      contents = base64encode(var.openapi_spec)
    }
  }

  gateway_config {
    backend_config {
      google_service_account = google_service_account.service_account.email
    }
  }

  labels = {
    "x-nitric-${var.stack_id}-name" = var.name
    "x-nitric-${var.stack_id}-type" = "api"
  }

  lifecycle {
    create_before_destroy = true
  }
}

resource "google_api_gateway_gateway" "gateway" {
  provider   = google-beta
  gateway_id = "${var.name}-gateway-${random_string.api_config_id.result}"
  api_config = google_api_gateway_api_config.api_config.id

  labels = {
    "x-nitric-${var.stack_id}-name" = var.name
    "x-nitric-${var.stack_id}-type" = "api"
  }
}


