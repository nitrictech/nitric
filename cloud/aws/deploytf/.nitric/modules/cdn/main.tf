locals {
  s3_origin_id = "publicOrigin"
}

resource "aws_cloudfront_origin_access_identity" "oai" {
  comment = "OAI for accessing S3 bucket"
}

data "aws_iam_policy_document" "allow_access_from_another_account" {
  version = "2012-10-17"
  statement {
    principals {
      type        = "AWS"
      identifiers = [
        aws_cloudfront_origin_access_identity.oai.iam_arn
      ]
    }

    effect = "Allow"

    actions = [
      "s3:GetObject",
    ]

    resources = [
      "${var.website_bucket_arn}/*",
    ]
  }
}

resource "aws_s3_bucket_policy" "allow_access_from_another_account" {
  bucket = var.website_bucket_id
  policy = data.aws_iam_policy_document.allow_access_from_another_account.json
}

resource "aws_cloudfront_function" "api-url-rewrite-function" {
  name    = "api-url-rewrite-function"
  runtime = "cloudfront-js-1.0"
  comment = "Rewrite API URLs routed to Nitric services"
  publish = true
  code    = file("${path.module}/scripts/api-url-rewrite.js")
}

resource "aws_cloudfront_function" "url-rewrite-function" {
  name    = "url-rewrite-function"
  runtime = "cloudfront-js-1.0"
  comment = "Rewrite URLs to default index document"
  publish = true
  code    = file("${path.module}/scripts/url-rewrite.js")
}

resource "aws_cloudfront_distribution" "s3_distribution" {
  default_root_object = var.website_index_document
  enabled = true
  
  origin {
    domain_name              = var.website_bucket_domain_name
    origin_id                = local.s3_origin_id

    s3_origin_config {
      origin_access_identity = aws_cloudfront_origin_access_identity.oai.cloudfront_access_identity_path
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

  default_cache_behavior {
    allowed_methods  = ["GET", "HEAD", "OPTIONS"]
    cached_methods   = ["GET", "HEAD", "OPTIONS"]
    target_origin_id = local.s3_origin_id

    forwarded_values {
      query_string = false

      cookies {
        forward = "none"
      }
    }

    viewer_protocol_policy = "redirect-to-https"
    min_ttl                = 0
    default_ttl            = 3600
    max_ttl                = 86400

    function_association {
      event_type = "viewer-request"
      function_arn = aws_cloudfront_function.url-rewrite-function.arn
    }
  }

  restrictions {
    geo_restriction {
      restriction_type = "none"
    }
  }

  viewer_certificate {
    cloudfront_default_certificate = true
  }

  custom_error_response {
    error_code = 404
    response_code = 200
    response_page_path = "/${var.website_error_document}"
  }

  custom_error_response {
    error_code = 403
    response_code = 200
    response_page_path = "/${var.website_error_document}"
  }

  tags = {
    "x-nitric-${var.stack_name}-name" = var.stack_name
    "x-nitric-${var.stack_name}-type" = "stack"
  }
}