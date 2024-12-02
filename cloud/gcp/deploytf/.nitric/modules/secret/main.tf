# Create a GCP Secret Manager secret
resource "google_secret_manager_secret" "secret" {
  # project = var.project_id
  secret_id = "${var.stack_id}-${var.secret_name}"
  labels = {
    "x-nitric-${var.stack_id}-name" = var.secret_name
    "x-nitric-${var.stack_id}-type" = "secret"
  }

  replication {
    user_managed {
      replicas {
        location = var.location
        dynamic "customer_managed_encryption" {
          for_each = var.cmek_key != "" ? [1] : []
          content {
            kms_key_name = var.cmek_key
          }
        }
      }
    }
  }
}
