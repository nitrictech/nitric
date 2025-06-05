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
