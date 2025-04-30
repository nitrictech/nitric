variable "name" {
  description = "The name of the API Gateway"
  type = string
}

variable "stack_id" {
  description = "The ID of the stack"
  type = string
}

variable "spec" {
  description = "Open API spec"
  type = string
}

variable "target_lambda_functions" {
  description = "The names of the target lambda functions"
  type = map(string)
}

variable "domain_names" {
  description = "A set of each domain name."
  type        = set(string)
}

variable "zone_ids" {
  description = "The id of the hosted zone mapped to the domain name."
  type = map(string)
}