
# Generate a random id for the bucket
resource "random_id" "bucket_id" {
  byte_length = 8

  keepers = {
    # Generate a new id each time we switch to a new AMI id
    bucket_name = var.bucket_name
  }
}

# AWS S3 bucket
resource "aws_s3_bucket" "bucket" {
  bucket = "${var.bucket_name}-${random_id.bucket_id.hex}"

  tags = {
    "x-nitric-${var.stack_id}-name" = var.bucket_name
    "x-nitric-${var.stack_id}-type" = "bucket"
  }
}

# Deploy bucket lambda invocation permissions
resource "aws_lambda_permission" "allow_bucket" {
  for_each = var.notification_targets
  action        = "lambda:InvokeFunction"
  function_name = each.value.arn
  principal     = "s3.amazonaws.com"
  source_arn    = aws_s3_bucket.bucket.arn
}

# Deploy lambda notifications
resource "aws_s3_bucket_notification" "bucket_notification" {
  count = length(var.notification_targets) > 0 ? 1 : 0
  
  bucket = aws_s3_bucket.bucket.id

  // make dynamic blocks for lambda function
  dynamic "lambda_function" {
    for_each = var.notification_targets
    content {
      lambda_function_arn = lambda_function.value.arn
      events              = lambda_function.value.events
      filter_prefix       = lambda_function.value.prefix
    }
  }
}