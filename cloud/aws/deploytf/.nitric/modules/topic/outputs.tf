output "topic_arn" {
  description = "The ARN of the deployed topic"
  value       =  aws_sns_topic.topic.arn
}

output "sfn_arn" {
  description = "The ARN of the deployed step function"
  value       = aws_sfn_state_machine.publish_to_topic.arn
}
