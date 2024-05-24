
# Deploy an SQS queue
resource "aws_sqs_queue" "queue" {
  name = var.queue_name
  tags = {
    "x-nitric-${var.stack_id}-name" = var.queue_name
    "x-nitric-${var.stack_id}-type" = "queue"
  }
}