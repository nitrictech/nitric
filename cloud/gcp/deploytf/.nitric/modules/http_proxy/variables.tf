variable "name" {
  description = "The name of the HTTP proxy gateway"
  type = string
}

variable "stack_id" {
  description = "The ID of the stack"
  type = string
}

variable "target_service_url" {
  description = "The URL of the service being proxied"
  type = string
}

variable "invoker_email" {
  description = "The email of the service account that will invoke the API"
  type = string
}