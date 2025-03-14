resource "aws_apigatewayv2_api" "api_gateway" {
  name          = var.name
  protocol_type = "HTTP"
  body = var.spec
  fail_on_warnings = true
  version = "v1"
  tags = {
    "x-nitric-${var.stack_id}-name" = var.name,
    "x-nitric-${var.stack_id}-type" = "api",
  }
}

resource "aws_apigatewayv2_stage" "stage" {
  api_id = aws_apigatewayv2_api.api_gateway.id
  name   = "$default"
  auto_deploy = true
}

# deploy lambda permissions for execution
resource "aws_lambda_permission" "apigw_lambda" {
  for_each      = var.target_lambda_functions
  action        = "lambda:InvokeFunction"
  function_name = each.value
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_apigatewayv2_api.api_gateway.execution_arn}/*/*/*"
}