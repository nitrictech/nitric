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

locals {
  service_account_name = "${substr("nitric-${provider::corefunc::str_kebab(var.nitric.name)}", 0, 20)}-${var.nitric.stack_id}"
}

# Create a service account for the google cloud run instance
resource "google_service_account" "service_account" {
  account_id   = local.service_account_name
  project      = var.project_id
  display_name = "${var.nitric.name} service account"
  description  = "Service account which runs the ${var.nitric.name}"

  depends_on = [ google_project_service.required_services ]
}