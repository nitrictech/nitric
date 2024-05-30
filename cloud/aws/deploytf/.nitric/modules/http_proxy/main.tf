resource "aws_apigatewayv2_api" "api_gateway" {
  name          = var.name
  protocol_type = "HTTP"
  tags = {
    "x-nitric-${var.stack_id}-name" = var.name,
    "x-nitric-${var.stack_id}-type" = "http-proxy",
  }
}

resource "aws_apigatewayv2_integration" "proxy_integration"{
  api_id           = aws_apigatewayv2_api.api_gateway.id
  integration_type = "AWS_PROXY"
  integration_uri  = var.target_lambda_function
}

resource "aws_apigatewayv2_route" "example" {
  api_id    = aws_apigatewayv2_api.api_gateway.id
  route_key = "ANY /{proxy+}"

  target = "integrations/${aws_apigatewayv2_integration.proxy_integration.id}"
}

resource "aws_apigatewayv2_stage" "stage" {
  api_id = aws_apigatewayv2_api.api_gateway.id
  name   = "$default"
  auto_deploy = true
}

# deploy lambda permissions for execution
resource "aws_lambda_permission" "apigw_lambda" {
  action        = "lambda:InvokeFunction"
  function_name = var.target_lambda_function
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_apigatewayv2_api.api_gateway.execution_arn}/*/*/*"
}
