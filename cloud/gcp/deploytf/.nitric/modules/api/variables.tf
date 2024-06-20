variable "name" {
  description = "The name of the API Gateway"
  type = string
}

variable "stack_id" {
  description = "The ID of the stack"
  type = string
}

variable "openapi_spec" {
  description = "The OpenAPI spec as a JSON string"
  type = string
}

variable "target_services" {
  description = "The map of target service names"
  type = map(string)
}