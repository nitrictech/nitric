variable "websocket_name" {
  description = "The name of the websocket."
  type        = string
}

variable "lambda_connect_target" {
  description = "The ARN of the lambda to send websocket connection events to"
  type        = string
}

variable "lambda_disconnect_target" {
  description = "The ARN of the lambda to send websocket disconnection events to"
  type        = string
}

variable "lambda_message_target" {
  description = "The ARN of the lambda to send websocket disconnection events to"
  type        = string
}

variable "stack_id" {
  description = "The ID of the Nitric stack"
  type        = string
}
