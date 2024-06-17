variable "secret_name" {
  description = "The name of the secret."
  type        = string
}

variable "stack_id" {
  description = "The ID of the Nitric stack"
  type        = string
}

variable "stack_name" {
  description = "The name of the Nitric stack, used to uniquely name the secret."
  type        = string
}
