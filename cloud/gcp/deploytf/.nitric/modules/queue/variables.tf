variable "queue_name" {
  description = "The name of the queue"
  type        = string
}

variable "kms_key" {
  description = "The KMS key to use for encryption"
  type        = string
  default     = ""
}

variable "stack_id" {
  description = "The ID of the Nitric stack"
  type        = string
}
