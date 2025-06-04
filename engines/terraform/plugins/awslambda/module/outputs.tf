output "nitric" {
  value = {
    id            = aws_lambda_function.function.arn
    http_endpoint = aws_lambda_function_url.endpoint.function_url
  }
}

