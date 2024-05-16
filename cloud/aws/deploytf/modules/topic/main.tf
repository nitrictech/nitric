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

# Subscriptions
resource "aws_sns_topic_subscription" "subscription_1" {
  topic_arn = aws_sns_topic.topic.arn
  protocol  = "email"
  endpoint  = "example1@example.com"
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