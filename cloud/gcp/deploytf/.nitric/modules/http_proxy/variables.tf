variable "name" {
  description = "The name of the API Gateway"
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

variable "project_id" {
  description = "The GCP project ID"
  type = string
}

variable "invoker_email" {
  description = "The email of the service account that will invoke the API"
  type = string
}