output "nitric" {
  value = {
    role = google_service_account.service_account.name
    id   = google_service_account.service_account.id
  }
}
