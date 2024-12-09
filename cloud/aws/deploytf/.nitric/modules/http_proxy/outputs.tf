output "arn" {
  value = aws_apigatewayv2_api.api_gateway.arn
}

output "endpoint" {
  value = aws_apigatewayv2_api.api_gateway.api_endpoint
}
