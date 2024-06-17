output "service_endpoint" { 
  value = google_cloud_run_service.service.status[0].url
}

output "service_account_email" {
  value = google_service_account.service_account.service_account_email
}

output "invoker_service_account_email" {
  value = google_service_account.invoker_service_account.service_account_email
}
