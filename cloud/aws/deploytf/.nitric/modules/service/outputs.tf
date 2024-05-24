output "lambda_arn" { 
  value = aws_lambda_function.function.arn 
}

output "lambda_function_name" {
  value = aws_lambda_function.function.function_name
}

output "invoke_arn" {
  value = aws_lambda_function.function.invoke_arn
}

output "role_arn" {
  value = aws_iam_role.role.arn
}
