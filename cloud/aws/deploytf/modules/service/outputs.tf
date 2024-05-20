output "lambda_arn" { 
  value = aws_lambda_function.lambda.arn 
}

output "role_arn" {
  value = aws_iam_role.role.arn
}
