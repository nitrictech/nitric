locals {
  domain_list = tolist(var.domain_names)
  base_names = { for domain in local.domain_list : domain => split(".", domain)[0] }
}

resource "aws_acm_certificate" "website-cert" {
  for_each = var.domain_names

  domain_name       = each.value
  validation_method = "DNS"
}

locals {
  domain_validation_options = { 
    for domain in var.domain_names : 
    domain => one(aws_acm_certificate.website-cert[domain].domain_validation_options)
  }
}

resource "aws_route53_record" "cert-validation-dns" {
  for_each = var.domain_names

  zone_id = var.zone_ids[each.value]
  ttl     = "600"
  name    = local.domain_validation_options[each.value].resource_record_name
  type    = local.domain_validation_options[each.value].resource_record_type
  records = [local.domain_validation_options[each.value].resource_record_value]

  depends_on = [aws_acm_certificate.website-cert]
}

resource "aws_acm_certificate_validation" "cert-validation" {
  for_each = var.domain_names

  certificate_arn         = aws_acm_certificate.website-cert[each.key].arn
  validation_record_fqdns = [aws_route53_record.cert-validation-dns[each.key].fqdn]
}

resource "aws_apigatewayv2_domain_name" "api_domain_name" {
  for_each = var.domain_names

  domain_name = each.value

  domain_name_configuration {
    certificate_arn = aws_acm_certificate_validation.cert-validation[each.key].certificate_arn
    endpoint_type   = "REGIONAL"
    security_policy = "TLS_1_2"
  }
}

resource "aws_apigatewayv2_api_mapping" "api_mapping" {
  for_each = var.domain_names
  
  api_id      = aws_apigatewayv2_api.api_gateway.id
  stage       = aws_apigatewayv2_stage.stage.id
  domain_name = aws_apigatewayv2_domain_name.api_domain_name[each.key].domain_name
}

resource "aws_route53_record" "api-dnsrecord" {
  for_each = var.domain_names

  zone_id = var.zone_ids[each.key]
  name    = local.base_names[each.key]
  type    = "A"

  alias {
    name = aws_apigatewayv2_domain_name.api_domain_name[each.key].domain_name_configuration[0].target_domain_name
    zone_id = aws_apigatewayv2_domain_name.api_domain_name[each.key].domain_name_configuration[0].hosted_zone_id
    evaluate_target_health = false
  }
}