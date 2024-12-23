resource "aws_apigatewayv2_api" "api_gateway" {
  name          = var.name
  protocol_type = "HTTP"
  body = var.spec
  fail_on_warnings = true
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

# look up existing certificate for domains
data "aws_acm_certificate" "cert" {
  for_each = var.domains
  domain = each.value
}

# deploy custom domain names
resource "aws_apigatewayv2_domain_name" "domain" {
  for_each = var.domains
  domain_name = each.value
  domain_name_configuration {
    certificate_arn = data.aws_acm_certificate.cert[each.key].arn
    endpoint_type = "REGIONAL"
    security_policy = "TLS_1_2"
  }
}