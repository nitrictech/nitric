output "queue_arn" {
  description = "The ARN of the deployed queue"
  value       =  aws_sqs_queue.queue.arn
}
