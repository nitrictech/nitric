output "kv_arn" {
  description = "The ARN of the deployed dynamodb table."
  value       =  aws_dynamodb_table.table.arn
}
