output "secret_arn" {
  description = "The ARN of the secret"
  value       =  aws_secretsmanager_secret.secret.arn
}
