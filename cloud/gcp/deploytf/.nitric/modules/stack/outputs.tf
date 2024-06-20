output "stack_id" {
  value = random_id.stack_id.hex
  description = "A unique id for this deployment"
}

output "base_compute_role" {
  value = google_project_iam_custom_role.base_role.id
  description = "The base compute role to use for the service"
}

output "iam_roles" {
  value = module.iam_roles
}