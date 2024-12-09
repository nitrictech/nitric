output "websocket_arn" {
  description = "The ARN of the deployed websocket API"
  value       = aws_apigatewayv2_api.websocket.arn
}

output "endpoint" {
  description = "The endpoint of the deployed websocket API"
  value       = aws_apigatewayv2_api.websocket.api_endpoint
}

output "websocket_exec_arn" {
  description = "The Execution ARN of the deployed websocket API"
  value       = aws_apigatewayv2_api.websocket.execution_arn
}
