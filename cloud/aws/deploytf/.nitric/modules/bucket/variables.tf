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

variable "cors_rules" {
  description = "CORS rules for the bucket"
  type = list(object({
    allowed_origins = list(string)
    allowed_methods = list(string)
    allowed_headers = list(string)
    expose_headers  = list(string)
    max_age_seconds = number
  }))
  default = []
}
