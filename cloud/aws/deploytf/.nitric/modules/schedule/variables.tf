variable "schedule_name" {
  description = "The name of the schedule"
  type        = string
}

variable "target_lambda_arn" {
  description = "The ARN of the target lambda function"
  type        = string
}

variable "schedule_expression" {
  description = "The schedule expression"
  type        = string
}

variable "schedule_timezone" {
  description = "The timezone for the schedule"
  type        = string
}

variable "stack_id" {
  description = "The ID of the Nitric stack"
  type        = string
}
