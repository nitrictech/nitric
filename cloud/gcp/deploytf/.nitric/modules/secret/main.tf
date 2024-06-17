
# # Create a new AWS secret manager secret
# resource "aws_secretsmanager_secret" "secret" {
#   name = var.secret_name
#   tags = {
#     "x-nitric-${var.stack_id}-name" = var.secret_name
#     "x-nitric-${var.stack_id}-type" = "secret"
#   }
# }


# Create a GCP Secret Manager secret
resource "google_secret_manager_secret" "secret" {
  # project = var.project_id
  secret_id = "${var.stack_name}-${var.secret_name}"
  labels = {
    "x-nitric-${var.stack_id}-name" = var.secret_name
    "x-nitric-${var.stack_id}-type" = "secret"
  }

  replication {
    auto {
    }
  }
}