variable "name" {
  description = "The name of the API Gateway"
  type = string
}

variable "stack_id" {
  description = "The ID of the stack"
  type = string
}

variable "project_id" {
  description = "The GCP project ID"
  type = string
}

variable "openapi_spec" {
  description = "The OpenAPI spec as a JSON string"
  type = string
}

variable "invoker_email" {
  description = "The email of the service account that will invoke the API"
  type = string
}

variable "target_services" {
  description = "The list of target services"
  type = map(object({
    name = string
    location = string
    url  = string
  }))
}