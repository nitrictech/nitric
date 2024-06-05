# Create role and policy to allow schedule to invoke lambda
resource "aws_iam_role" "role" {
  assume_role_policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Effect = "Allow",
        Principal = {
          Service = "scheduler.amazonaws.com"
        },
        Action = "sts:AssumeRole"
      }
    ]
  })
}

resource "aws_iam_role_policy" "role_policy" {
  role = aws_iam_role.role.id
  policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Effect = "Allow",
        Action = "lambda:InvokeFunction",
        Resource = var.target_lambda_arn
      }
    ]
  })
}

# Create an AWS eventbridge schedule
resource "aws_scheduler_schedule" "schedule" {
  flexible_time_window {
    mode = "OFF"
  }

  schedule_expression_timezone = var.schedule_timezone

  schedule_expression = var.schedule_expression

  target {
    arn      = var.target_lambda_arn
    role_arn = aws_iam_role.role.arn

    input = jsonencode({
        "x-nitric-schedule": var.schedule_name
    })
  }
}