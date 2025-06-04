# Create an ECR repository
resource "aws_ecr_repository" "repo" {
  name = var.nitric.name
}

data "aws_ecr_authorization_token" "ecr_auth" {
}

data "docker_image" "latest" {
  name = var.nitric.image_id
}

# Tag the provided docker image with the ECR repository url
resource "docker_tag" "tag" {
  source_image = length(data.docker_image.latest.repo_digest) > 0 ? data.docker_image.latest.repo_digest : data.docker_image.latest.id
  target_image = aws_ecr_repository.repo.repository_url
}

# Push the tagged image to the ECR repository
resource "docker_registry_image" "push" {
  name = aws_ecr_repository.repo.repository_url
  auth_config {
    address  = data.aws_ecr_authorization_token.ecr_auth.proxy_endpoint
    username = data.aws_ecr_authorization_token.ecr_auth.user_name
    password = data.aws_ecr_authorization_token.ecr_auth.password
  }
  triggers = {
    source_image_id = docker_tag.tag.source_image_id
  }
}

locals {
  role_name = var.nitric.identities["aws:iam:role"].role.arn
}

resource "aws_iam_role_policy_attachment" "basic-execution" {
  role       = local.role_name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

# Create a lambda function using the pushed image
resource "aws_lambda_function" "function" {
  function_name = var.nitric.name
  role          = local.role_name
  image_uri     = "${aws_ecr_repository.repo.repository_url}@${docker_registry_image.push.sha256_digest}"
  package_type  = "Image"
  timeout       = var.timeout
  memory_size   = var.memory
  ephemeral_storage {
    size = var.ephemeral_storage
  }
  environment {
    variables = merge(var.environment, var.nitric.env)
  }

  architectures = [var.architecture]
  dynamic "vpc_config" {
    for_each = length(var.subnet_ids) > 0 ? ["1"] : []
    content {
      subnet_ids         = var.subnet_ids
      security_group_ids = var.security_group_ids
    }
  }

  depends_on = [docker_registry_image.push]
}

resource "aws_lambda_function_url" "endpoint" {
  function_name = aws_lambda_function.function.function_name
  # qualifier          = "my_alias"
  authorization_type = var.function_url_auth_type

  # cors {
  #   allow_credentials = true
  #   allow_origins     = ["*"]
  #   allow_methods     = ["*"]
  #   allow_headers     = ["date", "keep-alive"]
  #   expose_headers    = ["keep-alive", "date"]
  #   max_age           = 86400
  # }
}
