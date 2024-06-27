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
    name = string
    url = string
    event_token = string
    invoker_service_account_email = string
    prefix = string
    events = list(string)
  }))
}

variable "storage_class" {
  description = "The class of storage used to store the bucket's contents. This can be STANDARD, NEARLINE, COLDLINE, ARCHIVE, or MULTI_REGIONAL."
  type        = string
  default     = "STANDARD"
}