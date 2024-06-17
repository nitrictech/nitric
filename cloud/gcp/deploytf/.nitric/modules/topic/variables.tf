variable "topic_name" {
  description = "The name of the bucket. This must be globally unique."
  type        = string
}

variable "stack_id" {
  description = "The ID of the Nitric stack"
  type        = string
}

variable "subscriber_services" {
  description = "The services to create subscriptions for"
  type = list(object({
    name                  = string
    url                   = string
    service_account_email = string
    event_token           = string
  }))
}
