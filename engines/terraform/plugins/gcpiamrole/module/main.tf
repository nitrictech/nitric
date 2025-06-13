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

resource "google_project_iam_custom_role" "role" {
  project = var.project_id
  role_id     = var.nitric.name
  title       = "IAM Role"
  description = "Custom role for IAM role permissions"
  permissions = var.trusted_actions

  depends_on = [ google_project_service.required_services ]
}
