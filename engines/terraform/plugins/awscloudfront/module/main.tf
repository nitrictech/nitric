locals {
  s3_origin_id = "publicOrigin"
  default_origin = {
    for k, v in var.nitric.origins : k => v
    if v.path == "/"
  }
}

resource "aws_cloudfront_function" "api-url-rewrite-function" {
  name    = "api-url-rewrite-function"
  runtime = "cloudfront-js-1.0"
  comment = "Rewrite API URLs routed to Nitric services"
  publish = true
  code    = file("${path.module}/scripts/api-url-rewrite.js")
}

resource "aws_cloudfront_distribution" "distribution" {
  enabled = true

  dynamic "origin" {
    for_each = var.nitric.origins

    content {
      # TODO: Only have services return their domain name instead? 
      domain_name = replace(origin.value.http_endpoint, "https://", "")
      origin_id = "${origin.key}"

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
    for_each = {
      for k, v in var.nitric.origins : k => v
      if v.path != "/"
    }

    content {
      path_pattern = "${each.value.path}*"

      function_association {
        event_type = "viewer-request"
        function_arn = aws_cloudfront_function.api-url-rewrite-function.arn
      }

      allowed_methods = ["GET","HEAD","OPTIONS","PUT","POST","PATCH","DELETE"]
      cached_methods = ["GET","HEAD","OPTIONS"]
      target_origin_id = "${ordered_cache_behavior.key}"

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
    allowed_methods  = ["GET", "HEAD", "OPTIONS", "PUT", "POST", "PATCH", "DELETE"]
    cached_methods   = ["GET", "HEAD", "OPTIONS"]
    target_origin_id = "${keys(local.default_origin)[0]}"
    viewer_protocol_policy = "redirect-to-https"
    min_ttl                = 0
    default_ttl            = 3600
    max_ttl                = 86400

    function_association {
      event_type = "viewer-request"
      function_arn = aws_cloudfront_function.api-url-rewrite-function.arn
    }

    forwarded_values {
      query_string = true
      cookies {
        forward = "all"
      }
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
}
