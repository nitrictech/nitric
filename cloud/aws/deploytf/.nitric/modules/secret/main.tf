

resource "random_id" "secret_id" {
  byte_length = 4

  prefix = "${var.secret_name}-"
  keepers = {
    # Generate a new id each time we switch to a new AMI id
    secret_name = var.secret_name
  }
}

# Create a new AWS secret manager secret
resource "aws_secretsmanager_secret" "secret" {
  name = random_id.secret_id.hex
  tags = {
    "x-nitric-${var.stack_id}-name" = var.secret_name
    "x-nitric-${var.stack_id}-type" = "secret"
  }
}