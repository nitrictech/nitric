# Create a new cloud scheduler job
resource "google_cloud_scheduler_job" "schedule" {
  name = var.schedule_name
  time_zone = var.schedule_timezone
  schedule = var.schedule_expression

  http_target {
    uri = "${var.target_service_url}/x-nitric-schedule/${var.schedule_name}?token=${var.service_token}"
    http_method = "POST"
    headers = {
      "Content-Type" = "application/json"
    }
    body = base64encode(jsonencode({
      "schedule": var.schedule_name
    }))
    oidc_token {
      service_account_email = var.target_service_invoker_email
    }
  }
}