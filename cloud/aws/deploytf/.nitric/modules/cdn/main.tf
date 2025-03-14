locals {
  s3_origin_id = "publicOrigin"
  website_files = join("", [for key, value in var.websites : join("", value.changed_files)])
}

resource "aws_cloudfront_origin_access_identity" "oai" {
  comment = "OAI for accessing S3 bucket"
}

resource "aws_s3_bucket_policy" "allow_access_from_another_account" {
  for_each = var.websites

  bucket = each.value.bucket_id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Principal = {
          AWS = aws_cloudfront_origin_access_identity.oai.iam_arn
        }
        Action = "s3:GetObject"
        Resource = "${each.value.bucket_arn}/*"
      }
    ]
  })
}

resource "aws_cloudfront_function" "api-url-rewrite-function" {
  name    = "api-url-rewrite-function"
  runtime = "cloudfront-js-1.0"
  comment = "Rewrite API URLs routed to Nitric services"
  publish = true
  code    = file("${path.module}/scripts/api-url-rewrite.js")
}

resource "aws_cloudfront_function" "url-rewrite-function" {
  for_each = {
    for name, website in var.websites : name => website
  }

  name    = "url-rewrite-function-${each.key}"
  runtime = "cloudfront-js-1.0"
  comment = "Rewrite URLs to default index document"
  publish = true
  code    = templatefile("${path.module}/scripts/url-rewrite.tmpl.js", {
    base_path = each.value.base_path
  })
}

resource "aws_cloudfront_distribution" "s3_distribution" {
  aliases = var.domain_name != "" ? [var.domain_name] : []

  default_root_object = var.root_website.index_document
  enabled = true

  dynamic "origin" {
    for_each = var.websites

    content {
      domain_name = origin.value.bucket_domain_name
      origin_id = origin.key

      s3_origin_config {
        origin_access_identity = aws_cloudfront_origin_access_identity.oai.cloudfront_access_identity_path
      }
    }
  }

  dynamic "origin" {
    for_each = var.apis

    content {
      domain_name = replace(origin.value.gateway_url, "https://", "")
      origin_id = origin.key

      custom_origin_config {
        origin_read_timeout = 30
        origin_protocol_policy = "https-only"
        origin_ssl_protocols = ["TLSv1.2", "SSLv3"]
        http_port = 80
        https_port = 443
      }
    }
  }

  dynamic "ordered_cache_behavior" {
    for_each = var.apis

    content {
      path_pattern = "api/${ordered_cache_behavior.key}/*"

      function_association {
        event_type = "viewer-request"
        function_arn = aws_cloudfront_function.api-url-rewrite-function.arn
      }

      allowed_methods = ["GET","HEAD","OPTIONS","PUT","POST","PATCH","DELETE"]
      cached_methods = ["GET","HEAD","OPTIONS"]
      target_origin_id = ordered_cache_behavior.key

      forwarded_values {
        query_string = true
        cookies {
          forward = "all"
        }
      }

      viewer_protocol_policy = "https-only"
    }
  }

  dynamic "ordered_cache_behavior" {
    for_each = {
      for name, website in var.websites : name => website
      if website.base_path != "/"
    }

    content {
      path_pattern = "${trimprefix(ordered_cache_behavior.value.base_path, "/")}*"

      function_association {
        event_type = "viewer-request"
        function_arn = aws_cloudfront_function.url-rewrite-function[ordered_cache_behavior.key].arn
      }

      allowed_methods  = ["GET", "HEAD", "OPTIONS"]
      cached_methods   = ["GET", "HEAD", "OPTIONS"]
      target_origin_id = ordered_cache_behavior.key

      forwarded_values {
        query_string = false
        cookies {
          forward = "none"
        }
      }

      viewer_protocol_policy = "redirect-to-https"
    }
  }

  default_cache_behavior {
    allowed_methods  = ["GET", "HEAD", "OPTIONS"]
    cached_methods   = ["GET", "HEAD", "OPTIONS"]
    target_origin_id = var.root_website.name
    viewer_protocol_policy = "redirect-to-https"
    min_ttl                = 0
    default_ttl            = 3600
    max_ttl                = 86400

    function_association {
      event_type = "viewer-request"
      function_arn = aws_cloudfront_function.url-rewrite-function[var.root_website.name].arn
    }

    forwarded_values {
      query_string = false

      cookies {
        forward = "none"
      }
    }
  }

  restrictions {
    geo_restriction {
      restriction_type = "none"
    }
  }

  viewer_certificate {
    acm_certificate_arn      = var.domain_name != "" ? var.certificate_arn : null
    cloudfront_default_certificate = var.domain_name == "" ? true : false
    ssl_support_method       = var.domain_name != "" ? "sni-only" : null
    minimum_protocol_version = "TLSv1.2_2021"
  }

  custom_error_response {
    error_code = 404
    response_code = 200
    response_page_path = "/${var.root_website.error_document}"
  }

  custom_error_response {
    error_code = 403
    response_code = 200
    response_page_path = "/${var.root_website.error_document}"
  }
}

resource "aws_route53_record" "cdn-dnsrecord" {
  count = var.domain_name != "" ? 1 : 0

  zone_id = var.zone_id
  name    = var.domain_name
  type    = "A"

  alias {
    name = aws_cloudfront_distribution.s3_distribution.domain_name
    zone_id = aws_cloudfront_distribution.s3_distribution.hosted_zone_id
    evaluate_target_health = false
  }
}

locals {
  // Check if the current OS is Windows by checking the drive letter in the path
  is_windows = can(regex("^[A-Za-z]:/", abspath(path.root)))
}

resource "terraform_data" "invalidate_cache" {
  count = var.skip_cache_invalidation ? 0 : 1

  triggers_replace = [
    local.website_files
  ]

  provisioner "local-exec" {
    command = "aws cloudfront create-invalidation --distribution-id ${aws_cloudfront_distribution.s3_distribution.id} --paths '/*'"
    interpreter = local.is_windows ? ["powershell", "-Command"] : null
  }
}