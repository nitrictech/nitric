variable "stack_name" {
  description = "The name of the nitric stack"
  type        = string
}

variable "cmek_enabled" {
  description = "Enable customer managed encryption keys"
  type        = bool
}

variable "location" {
  description = "The location to deploy the stack"
  type        = string
}