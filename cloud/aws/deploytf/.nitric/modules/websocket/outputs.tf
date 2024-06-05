output "websocket_arn" {
  description = "The ARN of the deployed websocket API"
  value       =  aws_apigatewayv2_api.websocket.arn
}

output "websocket_exec_arn" {
  description = "The Execution ARN of the deployed websocket API"
  value       =  aws_apigatewayv2_api.websocket.execution_arn
}
