locals {
  s3_origin_id = "publicOrigin"
}

resource "aws_cloudfront_origin_access_identity" "oai" {
  comment = "OAI for accessing S3 bucket"
}

resource "aws_s3_bucket_policy" "allow_access_from_another_account" {
  bucket = var.website_bucket_id
  policy = data.aws_iam_policy_document.allow_access_from_another_account.json
}

data "aws_iam_policy_document" "allow_access_from_another_account" {
  statement {
    principals {
      type        = "AWS"
      identifiers = [
        aws_cloudfront_origin_access_identity.oai.cloudfront_access_identity_path
      ]
    }

    actions = [
      "s3:GetObject",
    ]

    resources = [
      "${var.website_bucket_arn}/*",
    ]
  }
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
  origin {
    domain_name              = var.website_bucket_domain_name
    origin_access_control_id = aws_cloudfront_origin_access_identity.oai.id
    origin_id                = local.s3_origin_id
  }
  
  origin {
    domain_name = var.api_endpoint
    origin_id = var.website_name

    custom_origin_config {
      origin_read_timeout = 30
      origin_protocol_policy = "https-only"
      origin_ssl_protocols = ["TLSv1.2", "SSLv3"]
      http_port = 80
      https_port = 443
    }
  }

  enabled = true

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

  # Cache behavior with precedence 0
  ordered_cache_behavior {
    path_pattern = "api/${var.website_name}/*"

    function_association {
      event_type = "viewer-request"
      function_arn = aws_cloudfront_function.api-url-rewrite-function.arn
    }

    allowed_methods = ["GET","HEAD","OPTIONS","PUT","POST","PATCH","DELETE"]
    cached_methods = ["GET","HEAD","OPTIONS"]
    target_origin_id = var.website_name

    forwarded_values {
      query_string = true
      cookies {
        forward = "all"
      }
    }

    viewer_protocol_policy = "https-only"
  }

  default_root_object = var.website_index_document

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
    "x-nitric-${local.stack_name}-name" = var.stack_name
    "x-nitric-${local.stack_name}-type" = "stack"
  }
}