variable "name" {
  description = "The name of the topic"
  type        = string
}

variable "stack_name" {
  description = "The name of the stack"
  type        = string
}

variable "resource_group_name" {
  description = "The name of the resource group"
  type        = string
}

variable "location" {
  description = "The location of the topic"
  type        = string
}

variable "listeners" {
  description = "The list of listeners to notify"
  type = map(object({
    url                            = string
    active_directory_app_id_or_uri = string
    active_directory_tenant_id     = string
    event_token                    = string
  }))
}

