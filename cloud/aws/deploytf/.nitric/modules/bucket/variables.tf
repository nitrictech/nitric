variable "bucket_name" {
  description = "The name of the bucket. This must be globally unique."
  type        = string
}

variable "stack_id" {
  description = "The ID of the Nitric stack"
  type        = string
}

variable "notification_targets" {
  description = "The notification target configurations"
  type        = map(object({
    arn = string
    prefix = string
    events = list(string)
  }))
}
