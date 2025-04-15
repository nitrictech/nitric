output "endpoint" {
  value = aws_apigatewayv2_api.api_gateway.api_endpoint
}

output "arn" {
  value = aws_apigatewayv2_api.api_gateway.arn
}

output "id" {
  value = aws_apigatewayv2_api.api_gateway.id
}

output "stage_name" {
  value = aws_apigatewayv2_stage.stage.name
}