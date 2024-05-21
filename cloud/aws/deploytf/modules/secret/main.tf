
# Create a new AWS secret manager secret
resource "aws_secretsmanager_secret" "secret" {
  name = var.secret_name
  tags = {
    "x-nitric-${var.var.stack_id}-name" = var.secret_name
  }
}