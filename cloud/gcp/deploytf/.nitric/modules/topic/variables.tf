variable "topic_name" {
  description = "The name of the topic"
  type        = string
}

variable "stack_id" {
  description = "The ID of the Nitric stack"
  type        = string
}

variable "kms_key" {
  description = "The KMS key to use for encryption"
  type        = string
  default     = ""
}

variable "subscriber_services" {
  description = "The services to create subscriptions for"
  type = list(object({
    name                          = string
    url                           = string
    invoker_service_account_email = string
    event_token                   = string
  }))
}
