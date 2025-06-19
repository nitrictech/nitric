output "nitric" {
  value = {
    role = google_service_account.service_account.email
    id   = google_service_account.service_account.id
  }
}
