# Create a role policy attachment for each provided principal

resource "aws_iam_role_policy" "policy" {
  for_each = var.principals
  role       = each.value
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = {
        Effect   = "Allow"
        Action   = var.actions
        Resource = var.resources
    }
  })
}