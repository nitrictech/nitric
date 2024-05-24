# Create role and policy to allow schedule to invoke lambda
resource "aws_iam_role" "role" {
  name = "example"
  assume_role_policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Effect = "Allow",
        Principal = {
          Service = "events.amazonaws.com"
        },
        Action = "sts:AssumeRole"
      }
    ]
  })
}

resource "aws_iam_policy" "policy" {
  name = "example"
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
  name = var.schedule_name

  flexible_time_window {
    mode = "OFF"
  }

  schedule_expression = "rate(1 hours)"

  target {
    arn      = var.target_lambda_arn
    role_arn = aws_iam_role.role.arn

    input = jsonencode({
        "x-nitric-schedule": var.schedule_name
    })
  }
}