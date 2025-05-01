output "secret_arn" {
  description = "The ARN of the secret"
  value       = var.existing_secret_arn == "" ? one(aws_secretsmanager_secret.secret).arn : var.existing_secret_arn
}
