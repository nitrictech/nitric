variable "schedule_name" {
  description = "The name of the schedule"
  type        = string
}

variable "target_service_url" {
  description = "The URL of the target service"
  type        = string
}

variable "service_token" {
  description = "The token to authenticate with the target service"
  type        = string
}

# TODO: ensure this is parsed in the code if it's a rate schedule.
variable "schedule_expression" {
  description = "The schedule expression"
  type        = string
}

variable "schedule_timezone" {
  description = "The timezone for the schedule"
  type        = string
}