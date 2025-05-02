variable "secret_name" {
  description = "The name of the secret."
  type        = string
}

variable "existing_secret_arn" {
  description = "The ARN of the existing secret to import"
  type        = string
  default     = ""
}

variable "stack_id" {
  description = "The ID of the Nitric stack"
  type        = string
}
