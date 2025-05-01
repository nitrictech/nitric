
# Create a new AWS secret manager secret
resource "aws_secretsmanager_secret" "secret" {
  # Only create a new secret if we're not reusing an existing one
  count = var.existing_secret_arn == "" ? 1 : 0

  name = var.secret_name
  tags = {
    "x-nitric-${var.stack_id}-name" = var.secret_name
    "x-nitric-${var.stack_id}-type" = "secret"
  }
}