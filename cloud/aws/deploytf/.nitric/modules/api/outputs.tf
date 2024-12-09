output "endpoint" {
  value = aws_apigatewayv2_api.api_gateway.api_endpoint
}

output "arn" {
  value = aws_apigatewayv2_api.api_gateway.arn
}
