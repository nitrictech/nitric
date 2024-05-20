
# Create an ECR repository
resource "aws_ecr_repository" "foo" {
  name                 = var.service_name
}

# Tag the provided docker image with the ECR repository url
resource "docker_image" "tag" {
  source_image = var.image
  target_image = aws_ecr_repository.repo.repository_url
}

# Push the tagged image to the ECR repository
resource "docker_registry_image" "push" {
  name          = docker_image.tag.name
  repository    = aws_ecr_repository.repo.repository_url
  tag           = docker_image.tag.tag
  build_context = docker_image.tag.build_context
}

# Create a role for the lambda function
resource "aws_iam_role" "role" {
  name = var.service_name
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Principal = {
          Service = "lambda.amazonaws.com"
        }
        Action = "sts:AssumeRole"
      }
    ]
  })
}

# Create a lambda function using the pushed image
resource "aws_lambda_function" "function" {
  function_name = var.service_name
  role          = aws_iam_role.role.arn
  image_uri     = docker_registry_image.push.image_url
  environment {
    variables = var.environment
  }
}