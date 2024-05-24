terraform {
  required_providers {
    docker = {
      source  = "kreuzwerker/docker"
      version = "3.0.2"
    }
  }
}

# Create an ECR repository
resource "aws_ecr_repository" "repo" {
  name = var.service_name
}

data "aws_ecr_authorization_token" "ecr_auth" {
  depends_on = [aws_ecr_repository.repo]
}

provider "docker" {
  registry_auth {
    address  = data.aws_ecr_authorization_token.ecr_auth.proxy_endpoint
    username = data.aws_ecr_authorization_token.ecr_auth.user_name
    password = data.aws_ecr_authorization_token.ecr_auth.password
  }
}

# Tag the provided docker image with the ECR repository url
resource "docker_tag" "tag" {
  source_image = var.image
  target_image = aws_ecr_repository.repo.repository_url
}

# Push the tagged image to the ECR repository
resource "docker_registry_image" "push" {
  name = aws_ecr_repository.repo.repository_url
  triggers = {
    source_image_id = docker_tag.tag.source_image_id
  }
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

# TODO Make a common policy and attach separately
# as a base common compute policy
resource "aws_iam_role_policy" "resource-list-access" {
  name = "resource-list-access"
  role = aws_iam_role.role.name

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "sns:ListTopics",
          "sqs:ListQueues",
          "dynamodb:ListTables",
          "s3:ListAllMyBuckets",
          "tag:GetResources",
          "apigateway:GET",
        ]
        Resource = "*"
      }
    ]
  })
}


resource "aws_iam_role_policy_attachment" "basic-execution" {
  role       = aws_iam_role.role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

# Create a lambda function using the pushed image
resource "aws_lambda_function" "function" {
  function_name = "${var.service_name}-${var.stack_id}"
  role          = aws_iam_role.role.arn
  image_uri     = "${aws_ecr_repository.repo.repository_url}:latest"
  package_type  = "Image"
  # TODO: Make configurable
  timeout = 30
  environment {
    variables = var.environment
  }

  depends_on = [docker_registry_image.push]

  tags = {
    "x-nitric-${var.stack_id}-name" = var.service_name,
    "x-nitric-${var.stack_id}-type" = "service",
  }
}
