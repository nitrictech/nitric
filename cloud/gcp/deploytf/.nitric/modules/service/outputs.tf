output "service_endpoint" { 
  value = google_cloud_run_service.service.status[0].url
}

output "service_account_email" {
  value = google_service_account.service_account.email
}

output "invoker_service_account_email" {
  value = google_service_account.invoker_service_account.email
}

output "event_token" {
  value = random_password.event_token.result
}

output "service_name" {
  value = google_cloud_run_service.service.name
}