#  Deploy a websocket API gateway

resource "aws_apigatewayv2_api" "websocket" {
  name          = var.websocket_name
  protocol_type = "WEBSOCKET"
  route_selection_expression = "$request.body.action"
  tags = {
    "x-nitric-${var.stack_id}-name" = var.websocket_name
    "x-nitric-${var.stack_id}-type" = "websocket"
  }
}

resource "aws_apigatewayv2_integration" "default" {
  api_id           = aws_apigatewayv2_api.websocket.id
  integration_type = "AWS_PROXY"
  integration_uri  = var.lambda_message_target
}

# Create an integration for the connect route
resource "aws_apigatewayv2_integration" "disconnect" {
  api_id           = aws_apigatewayv2_api.websocket.id
  integration_type = "AWS_PROXY"
  integration_uri  = var.lambda_disconnect_target
}

# Create an integration for the connect route
resource "aws_apigatewayv2_integration" "connect" {
  api_id           = aws_apigatewayv2_api.websocket.id
  integration_type = "AWS_PROXY"
  integration_uri  = var.lambda_connect_target
}

# Create the default route for the websocket
resource "aws_apigatewayv2_route" "default" {
  api_id    = aws_apigatewayv2_api.websocket.id
  route_key = "$default"
  target    = "integrations/${aws_apigatewayv2_integration.default.id}"
}

# Create the connect route for the websocket
resource "aws_apigatewayv2_route" "connect" {
  api_id    = aws_apigatewayv2_api.websocket.id
  route_key = "$connect"
  target    = "integrations/${aws_apigatewayv2_integration.connect.id}"
}

# Create the disconnect route for the websocket
resource "aws_apigatewayv2_route" "disconnect" {
  api_id    = aws_apigatewayv2_api.websocket.id
  route_key = "$disconnect"
  target    = "integrations/${aws_apigatewayv2_integration.disconnect.id}"
}

# Create execution lambda permissions for the websocket
resource "aws_lambda_permission" "websocket-message" {
  action        = "lambda:InvokeFunction"
  function_name = var.lambda_message_target
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_apigatewayv2_api.websocket.execution_arn}/*/*"
}

resource "aws_lambda_permission" "websocket-connect" {
  action        = "lambda:InvokeFunction"
  function_name = var.lambda_connect_target
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_apigatewayv2_api.websocket.execution_arn}/*/*"
}

resource "aws_lambda_permission" "websocket-disconnect" {
  action        = "lambda:InvokeFunction"
  function_name = var.lambda_disconnect_target
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_apigatewayv2_api.websocket.execution_arn}/*/*"
}

# create a stage for the api gateway
resource "aws_apigatewayv2_stage" "stage" {
  api_id      = aws_apigatewayv2_api.websocket.id
  name        = "ws"
  auto_deploy = true

  tags = {
    "x-nitric-${var.stack_id}-name" = "${var.websocket_name}DefaultStage"
    "x-nitric-${var.stack_id}-type" = "websocket"
  }
}
