output "name" {
  description = "The name of the deployed secret."
  value       = google_secret_manager_secret.secret.name
}