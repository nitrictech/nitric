output "nitric" {
  value = {
    role = google_project_iam_custom_role.role
    id   = google_project_iam_custom_role.role.name
  }
}
