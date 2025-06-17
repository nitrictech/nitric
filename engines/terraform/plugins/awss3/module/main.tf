# AWS S3 bucket

locals {
  normalized_nitric_name = provider::corefunc::str_kebab(var.nitric.name)
}

resource "aws_s3_bucket" "bucket" {
  bucket = "${var.nitric.stack_id}-${local.normalized_nitric_name}"
  tags   = var.tags
}

locals {
  read_actions = [
    "s3:GetObject",
    "s3:ListBucket",
  ]
  write_actions = [
    "s3:PutObject",
  ]
  delete_actions = [
    "s3:DeleteObject",
  ]
  relative_content_path = "${path.root}/../../../${var.nitric.content_path}"
  content_files = var.nitric.content_path != "" ? fileset(local.relative_content_path, "**/*") : []
}

# Upload each file to S3 (only if files exist)
resource "aws_s3_object" "files" {
  for_each = toset(local.content_files)
  
  bucket = aws_s3_bucket.bucket.bucket
  key    = each.value
  source = "${local.relative_content_path}/${each.value}"
  
  etag = filemd5("${local.relative_content_path}/${each.value}")
  
  content_type = lookup({
    "html" = "text/html"
    "css"  = "text/css"
    "js"   = "application/javascript"
    "json" = "application/json"
    "png"  = "image/png"
    "jpg"  = "image/jpeg"
    "jpeg" = "image/jpeg"
    "gif"  = "image/gif"
    "svg"  = "image/svg+xml"
    "pdf"  = "application/pdf"
    "txt"  = "text/plain"
  }, reverse(split(".", each.value))[0], "application/octet-stream")
}


resource "aws_iam_role_policy" "access_policy" {
  for_each = var.nitric.services
  name     = "${local.normalized_nitric_name}-${provider::corefunc::str_kebab(each.key)}"
  role     = each.value.identities["aws:iam:role"].role.name


  # Terraform's "jsonencode" function converts a
  # Terraform expression result to valid JSON syntax.
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = distinct(concat(
          contains(each.value.actions, "read") ? local.read_actions : [],
          contains(each.value.actions, "write") ? local.write_actions : [],
          contains(each.value.actions, "delete") ? local.delete_actions : []
          )
        )
        Effect   = "Allow"
        Resource = [
          aws_s3_bucket.bucket.arn,
          "${aws_s3_bucket.bucket.arn}/*"
        ]
      },
    ]
  })
}
