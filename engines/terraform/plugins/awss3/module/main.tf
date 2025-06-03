
# AWS S3 bucket
resource "aws_s3_bucket" "bucket" {
  bucket = "${var.nitric.stack_id}-${var.nitric.name}"
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
  name     = "${var.nitric.name}-${each.value.name}"
  role     = each.value.identities["aws:iam"].id


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
        Resource = aws_s3_bucket.bucket.arn
      },
    ]
  })
}
