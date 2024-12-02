variable "secret_name" {
  description = "The name of the secret."
  type        = string
}

variable "stack_id" {
  description = "The ID of the Nitric stack"
  type        = string
}

variable "location" {
  description = "location of the secret"
  type        = string
}

variable "cmek_key" {
  description = "The KMS key to use for encryption"
  type        = string
  default     = ""
}
