variable "cron_expression" {
  description = "The cron expression for the schedule"
  type = string  
}

variable "name" {
  description = "The name of the schedule"
  type = string
}

variable "container_app_environment_id" {
  description = "The container app environment id"
  type = string
}

variable "target_app_id" {
  description = "The target app id for the schedule"
  type = string
}

variable "target_event_token" {
  description = "The target event token for the schedule"
  type = string
}