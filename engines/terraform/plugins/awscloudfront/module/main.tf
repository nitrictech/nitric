locals {
  s3_origin_id = "publicOrigin"
  default_origin = {
    for k, v in var.nitric.origins : k => v
    if v.path == "/"
  }
  s3_bucket_origins = {
    for k, v in var.nitric.origins : k => v
    if contains(keys(v.resources), "aws_s3_bucket")
  }
  lambda_origins = {
    for k, v in var.nitric.origins : k => v
    if contains(keys(v.resources), "aws_lambda_function")
  }
  non_vpc_origins = {
    for k, v in var.nitric.origins : k => v
    if !contains(keys(v.resources), "aws_lb")
  }
  vpc_origins = {
    for k, v in var.nitric.origins : k => v
    if contains(keys(v.resources), "aws_lb")
  }
}

resource "aws_cloudfront_vpc_origin" "vpc_origin" {
  for_each = local.vpc_origins

  vpc_origin_endpoint_config {
    name = each.key
    arn = each.value.resources["aws_lb"]
    http_port = each.value.resources["aws_lb:target_port"]
    https_port = each.value.resources["aws_lb:target_port"]
    origin_protocol_policy = "https-only"

    origin_ssl_protocols {
      items    = ["TLSv1.2"]
      quantity = 1
    }
  }
}

resource "aws_cloudfront_origin_access_control" "lambda_oac" {
  count = length(local.lambda_origins) > 0 ? 1 : 0

  name                              = "lambda-oac"
  origin_access_control_origin_type = "lambda"
  signing_behavior                  = "always"
  signing_protocol                  = "sigv4"
}

resource "aws_cloudfront_origin_access_control" "s3_oac" {
  count = length(local.s3_bucket_origins) > 0 ? 1 : 0

  name                              = "s3-oac"
  origin_access_control_origin_type = "s3"
  signing_behavior                  = "always"
  signing_protocol                  = "sigv4"
}

# Allow cloudfront to execute the function urls of any provided AWS lambda functions
resource "aws_lambda_permission" "allow_cloudfront_to_execute_lambda" {
  for_each = local.lambda_origins

  function_name = each.value.resources["aws_lambda_function"]
  principal = "cloudfront.amazonaws.com"
  action = "lambda:InvokeFunctionUrl"
  source_arn = aws_cloudfront_distribution.distribution.arn
}

resource "aws_s3_bucket_policy" "allow_bucket_access" {
  for_each = local.s3_bucket_origins

  bucket = replace(each.value.id, "arn:aws:s3:::", "")

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Principal = {
          Service = "cloudfront.amazonaws.com"
        }
        Action = "s3:GetObject"
        Resource = "${each.value.id}/*"
        Condition = {
          StringEquals = {
            "AWS:SourceArn" = aws_cloudfront_distribution.distribution.arn
          }
        }
      }
    ]
  })
}

resource "aws_cloudfront_function" "api-url-rewrite-function" {
  name    = "api-url-rewrite-function"
  runtime = "cloudfront-js-1.0"
  comment = "Rewrite API URLs routed to Nitric services"
  publish = true
  code    = templatefile("${path.module}/scripts/url-rewrite.js", {
    base_paths = join(",", [for k, v in var.nitric.origins : v.path])
  })
}

resource "aws_cloudfront_distribution" "distribution" {
  enabled = true

  dynamic "origin" {
    for_each = local.non_vpc_origins

    content {
      # TODO: Only have services return their domain name instead? 
      domain_name = origin.value.domain_name
      origin_id = "${origin.key}"
      origin_access_control_id = contains(keys(origin.value.resources), "aws_lambda_function") ? aws_cloudfront_origin_access_control.lambda_oac[0].id : contains(keys(origin.value.resources), "aws_s3_bucket") ? aws_cloudfront_origin_access_control.s3_oac[0].id : null

      dynamic "custom_origin_config" {
        for_each = !contains(keys(origin.value.resources), "aws_s3_bucket") ? [1] : []

        content {
          origin_read_timeout = 30
          origin_protocol_policy = "https-only"
          origin_ssl_protocols = ["TLSv1.2", "SSLv3"]
          http_port = 80
          https_port = 443
        }
      }
    }
  }

  dynamic "origin" {
    for_each = local.vpc_origins

    content {
      domain_name = origin.value.domain_name
      origin_id = "${origin.key}"
      vpc_origin_config {
        vpc_origin_id = aws_cloudfront_vpc_origin.vpc_origin[origin.key].id
      }
    }
  }

  dynamic "ordered_cache_behavior" {
    for_each = {
      for k, v in var.nitric.origins : k => v
      if v.path != "/"
    }

    content {
      path_pattern = "${ordered_cache_behavior.value.path}*"

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
