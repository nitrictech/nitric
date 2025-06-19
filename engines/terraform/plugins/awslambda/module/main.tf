# Convert standard cron expressions to AWS CloudWatch format
locals {
  split_cron_expression = {
    for key, schedule in var.nitric.schedules : key => split(" ", schedule.cron_expression)
  }

  # Apply the either-or operator rule: if DOM is *, set DOW to ?
  transformed_cron_expression = {
    for key, fields in local.split_cron_expression : key => [
      for i, field in fields : 
        (i == 4 && fields[2] == "*" && field == "*") ? "?" : field
    ]
  }

  # Convert standard cron to AWS CloudWatch cron format
  # AWS requires 6 fields: Minutes Hours Day-of-month Month Day-of-week Year
  # Either Day-of-month or Day-of-week must be ? (either-or operator)
  # Day-of-week is 1-7 (Sunday-Saturday) instead of 0-6
  convert_cron_to_aws = {
    for key, schedule in var.nitric.schedules : key => {
      cron_expression = schedule.cron_expression
      path           = schedule.path
      # Convert the standard cron expression to an AWS CloudWatch cron expression
      # Apply the either-or operator rule and add year field
      aws_cron       = "cron(${join(" ", local.transformed_cron_expression[key])} *)"
    }
  }
}

# Create an ECR repository
resource "aws_ecr_repository" "repo" {
  name = var.nitric.name
}

data "aws_ecr_authorization_token" "ecr_auth" {
}

data "docker_image" "latest" {
  name = var.nitric.image_id
}

locals {
  lambda_name = "${var.nitric.stack_id}-${var.nitric.name}"
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

resource "aws_iam_role_policy_attachment" "basic-execution" {
  role       = var.nitric.identities["aws:iam:role"].role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

# Create a lambda function using the pushed image
resource "aws_lambda_function" "function" {
  function_name = local.lambda_name
  role          = var.nitric.identities["aws:iam:role"].role.arn
  image_uri     = "${aws_ecr_repository.repo.repository_url}@${docker_registry_image.push.sha256_digest}"
  package_type  = "Image"
  timeout       = var.timeout
  memory_size   = var.memory
  ephemeral_storage {
    size = var.ephemeral_storage
  }
  environment {
    variables = merge(var.environment, var.nitric.env, {
      NITRIC_STACK_ID = var.nitric.stack_id
    })
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

  invoke_mode = "RESPONSE_STREAM"

  # cors {
  #   allow_credentials = true
  #   allow_origins     = ["*"]
  #   allow_methods     = ["*"]
  #   allow_headers     = ["date", "keep-alive"]
  #   expose_headers    = ["keep-alive", "date"]
  #   max_age           = 86400
  # }
}

# Deploy schedule events per schedule for this service
resource "aws_cloudwatch_event_rule" "schedule" {
  for_each = var.nitric.schedules

  name                = "${var.nitric.name}-${each.key}-schedule"
  description         = "Schedule for ${var.nitric.name} - ${each.key}"
  schedule_expression = local.convert_cron_to_aws[each.key].aws_cron
}

resource "aws_cloudwatch_event_target" "schedule_target" {
  for_each = var.nitric.schedules

  rule      = aws_cloudwatch_event_rule.schedule[each.key].name
  target_id = "${var.nitric.name}-${each.key}-target"
  arn       = aws_lambda_function.function.arn

  input = jsonencode({
    # Events will be forwarded to the lambda function with the path as the event payload
    path = each.value.path
  })
}

resource "aws_lambda_permission" "allow_cloudwatch" {
  for_each = var.nitric.schedules

  statement_id  = "AllowExecutionFromCloudWatch-${each.key}"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.function.function_name
  principal     = "events.amazonaws.com"
  source_arn    = aws_cloudwatch_event_rule.schedule[each.key].arn
}

