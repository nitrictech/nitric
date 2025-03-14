resource "aws_route53_zone" "main" {
  name = var.domain_name
}

resource "aws_acm_certificate" "resume-app-cert" {
  domain_name       = var.domain_name
  validation_method = "DNS"

  subject_alternative_names = [var.wildcard_domain_name]

  tags = {
    Name = var.aws_certificate_name
  }
}

resource "aws_route53_record" "dev-ns" {
  zone_id = aws_route53_zone.main.zone_id
  name    = "dev.example.com"
  type    = "NS"
  ttl     = "30"
  records = aws_route53_zone.dev.name_servers
}