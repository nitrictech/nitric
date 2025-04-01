locals {
  domain_parts   = split(".", var.domain_name)
  parent_domain  = join(".", slice(local.domain_parts, 1, length(local.domain_parts)))
  base_name      = local.domain_parts[0]
}

# Try to find hosted zone by full domain first
data "aws_route53_zone" "primary" {
  name         = var.domain_name
  private_zone = false
}

# Fallback to parent domain if the full domain hosted zone isn't found
data "aws_route53_zone" "parent" {
  name         = local.parent_domain
  private_zone = false
}

locals {
  zone_id = try(data.aws_route53_zone.primary.id, data.aws_route53_zone.parent.id)
}

resource "aws_acm_certificate" "website-cert" {
  domain_name       = var.domain_name
  validation_method = "DNS"
}

resource "aws_route53_record" "cert-validation-dns" {
  for_each = { for dvo in aws_acm_certificate.website-cert.domain_validation_options : dvo.resource_record_name => dvo }

  zone_id = local.zone_id
  name    = each.value.resource_record_name
  type    = each.value.resource_record_type
  records = [each.value.resource_record_value]
  ttl     = "600"
}

resource "aws_acm_certificate_validation" "cert-validation" {
  certificate_arn         = aws_acm_certificate.website-cert.arn
  validation_record_fqdns = [for r in aws_route53_record.cert-validation-dns : r.fqdn]
}

resource "aws_apigatewayv2_domain_name" "api_domain_name" {
  domain_name = var.domain_name

  domain_name_configuration {
    certificate_arn = aws_acm_certificate.website-cert.arn
    endpoint_type   = "REGIONAL"
    security_policy = "TLS_1_2"
  }
}

resource "aws_apigatewayv2_api_mapping" "api_mapping" {
  api_id      = var.api_id
  stage       = var.api_stage_name
  domain_name = aws_apigatewayv2_domain_name.api_domain_name.domain_name
}

resource "aws_route53_record" "api-dnsrecord" {
  zone_id = local.zone_id
  name    = local.base_name
  type    = "A"

  alias {
    name = aws_apigatewayv2_domain_name.api_domain_name.domain_name_configuration[0].target_domain_name
    zone_id = aws_apigatewayv2_domain_name.api_domain_name.domain_name_configuration[0].hosted_zone_id
    evaluate_target_health = false
  }
}