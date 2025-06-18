locals {
  required_services = [
    // Enable IAM API
    "iam.googleapis.com",
  ]
}

# Enable the required services
resource "google_project_service" "required_services" {
  for_each = toset(local.required_services)

  service = each.key
  project = var.project_id
  # Leave API enabled on destroy
  disable_on_destroy = false
  disable_dependent_services = false
}

# Create a random ID for the service name, so that it confirms to regex restrictions
resource "random_string" "service_account_id" {
  length  = 30 - length(var.nitric.stack_id)
  special = false
  upper   = false
}

# Create a service account for the google cloud run instance
resource "google_service_account" "service_account" {
  account_id   = "${random_string.service_account_id.id}${var.nitric.stack_id}"
  project      = var.project_id
  display_name = "${var.nitric.name} service account"
  description  = "Service account which runs the ${var.nitric.name}"

  depends_on = [ google_project_service.required_services ]
}