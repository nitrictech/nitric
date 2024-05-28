# Generate a random id for the topic
resource "random_id" "topic_id" {
  byte_length = 8

  keepers = {
    # Generate a new id each time we switch to a new name
    topic_name = var.topic_name
  }
}


# AWS SNS Topic
resource "aws_sns_topic" "topic" {
  name = "${var.topic_name}-${random_id.bucket_id.hex}"

  tags = {
    "x-nitric-${var.stack_id}-name" = var.topic_name
    "x-nitric-${var.stack_id}-type" = "topic"
  }
}

# Loop over the subsribers and deploy subscriptions and permissions
resource "aws_sns_topic_subscription" "subscription" {
  for_each = var.lambda_subscribers

  topic_arn = aws_sns_topic.topic.arn
  protocol  = "lambda"
  endpoint  = each.value.lambda_arn
}

resource "aws_lambda_permission" "sns" {
  for_each = var.lambda_subscribers

  statement_id  = "AllowExecutionFromSNS"
  action        = "lambda:InvokeFunction"
  function_name = each.value
  principal     = "sns.amazonaws.com"
  source_arn    = aws_sns_topic.topic.arn
}

resource "aws_iam_role" "sns_publish_role" {
  name = "${var.topic_name}-sns-publish-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Principal = {
          Service = "states.amazonaws.com"
        }
        Action = "sts:AssumeRole"
      }
    ]
  })
}

# Attach the policy to the role inline
resource "aws_iam_role_policy" "test_policy" {
  role = aws_iam_role.test_role.id

  # Terraform's "jsonencode" function converts a
  # Terraform expression result to valid JSON syntax.
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect   = "Allow"
        Action   = "sns:Publish"
        Resource = aws_sns_topic.topic.arn
      }
    ]
  })
}

# Create a step function that publishes to the topic
resource "aws_sfn_state_machine" "publish_to_topic" {
  name     = "${var.topic_name}-publish-to-topic"
  role_arn = aws_iam_role.sns_publish_role.arn
  definition = jsonencode({
    Comment = "",
    StartAt = "Wait",
    States = {
      Wait = {
        Type = "Wait",
        SecondsPath : "$.seconds",
        Next = "Publish"
      }
      Publish = {
        Type     = "Task",
        Resource = "arn:aws:states:::sns:publish",
        Parameters = {
          TopicArn = aws_sns_topic.topic.arn,
          "Message.$" : "$.message",
        },
        End = true
      }
    }
  })
}


