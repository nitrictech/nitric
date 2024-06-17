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
    url = string
    event_token = string
    prefix = string
    events = list(string)
  }))
}

variable "bucket_location" {
  description = "The location where the bucket and its contents are stored."
  type        = string
  # default     = "US"
}

variable "project_id" {
  description = "The ID of the Google Cloud project where the bucket is created."
  type        = string
}

variable "storage_class" {
  description = "The class of storage used to store the bucket's contents. This can be STANDARD, NEARLINE, COLDLINE, ARCHIVE, or MULTI_REGIONAL."
  type        = string
  default     = "STANDARD"
}