output "nitric" {
  value = {
    id            = aws_lambda_function.function.arn
    domain_name = split("/", aws_lambda_function_url.endpoint.function_url)[2]
    raw_type = "aws_lambda_function"
    raw = aws_lambda_function.function
  }
}

