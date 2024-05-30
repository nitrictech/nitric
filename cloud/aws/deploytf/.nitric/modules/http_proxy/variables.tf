variable "name" {
  description = "The name of the API Gateway"
  type = string
}

variable "stack_id" {
  description = "The ID of the stack"
  type = string
}

variable "target_lambda_function" {
  description = "The name or arn of the target lambda functin"
  type = string
}
