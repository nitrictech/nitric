locals {
  s3_origin_id = "publicOrigin"
  default_origin = {
    for k, v in var.nitric.origins : k => v
    if v.path == "/"
  }
  s3_bucket_origins = {
    for k, v in var.nitric.origins : k => v
    if lookup(v, "raw.aws_s3_bucket", null) != null
  }
  lambda_origins = {
    for k, v in var.nitric.origins : k => v
    if lookup(v, "raw.aws_lambda_function", null) != null
  }
}

resource "aws_cloudfront_origin_access_identity" "oai" {
  comment = "OAI for accessing S3 bucket"
}

# Allow cloudfront to execute the function urls of any provided AWS lambda functions
resource "aws_lambda_permission" "allow_cloudfront_to_execute_lambda" {
  for_each = local.lambda_origins

  function_name = each.value.raw["aws_lambda_function"]
  principal = "cloudfront.amazonaws.com"
  action = "lambda:InvokeFunction"
  source_arn = aws_cloudfront_distribution.distribution.arn
}

resource "aws_s3_bucket_policy" "allow_bucket_access" {
  for_each = local.s3_bucket_origins

  bucket = each.value.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Principal = {
          AWS = aws_cloudfront_origin_access_identity.oai.iam_arn
        }
        Action = "s3:GetObject"
        Resource = "${each.value.id}/*"
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
    for_each = var.nitric.origins

    content {
      # TODO: Only have services return their domain name instead? 
      domain_name = origin.value.domain_name
      origin_id = "${origin.key}"

      dynamic "s3_origin_config" {
        for_each = lookup(origin.value, "raw.aws_s3_bucket", null) != null ? [1] : []

        content {
          origin_access_identity = aws_cloudfront_origin_access_identity.oai.iam_arn
        }
      }

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
