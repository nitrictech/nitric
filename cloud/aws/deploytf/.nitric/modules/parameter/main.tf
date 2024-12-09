resource "random_string" "random" {
  length  = 4
  special = false
}

locals {
  policy_name = "nitric-param-access-${random_string.random.result}"
}

# Create a new SSM Parameter Store parameter
resource "aws_ssm_parameter" "text_parameter" {
  name      = var.parameter_name
  type      = "String"
  value     = var.parameter_value
  data_type = "text"
}

# Create the access policy
resource "aws_iam_policy" "access_policy" {
  name = local.policy_name
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect   = "Allow"
      Action   = "ssm:GetParameter"
      Resource = aws_ssm_parameter.text_parameter.arn
    }]
  })
}

# Create the role policy attachment
resource "aws_iam_role_policy_attachment" "policy_attachment" {
  for_each   = var.access_role_names
  role       = each.value
  policy_arn = aws_iam_policy.access_policy.arn
}
