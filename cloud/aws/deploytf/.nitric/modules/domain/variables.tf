variable "domain_name" {
  description = "The name of the domain. This must be globally unique."
  type        = string
}

variable "stack_id" {
  description = "The ID of the Nitric stack"
  type        = string
}

variable "api_id" {
  description = "The ID of the API"
  type        = string
}

variable "api_stage_name" {
  description = "The stage for the API"
  type        = string
}
