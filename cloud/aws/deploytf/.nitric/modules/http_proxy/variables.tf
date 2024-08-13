variable "name" {
  description = "The name of the HTTP proxy gateway"
  type = string
}

variable "stack_id" {
  description = "The ID of the stack"
  type = string
}

variable "target_lambda_function" {
  description = "The name or arn of the lambda function being proxied"
  type = string
}
