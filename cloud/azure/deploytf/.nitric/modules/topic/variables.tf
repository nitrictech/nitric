variable "name" {
  description = "The name of the topic"
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
    tenant_id = string
    client_secret = string
    client_id = string
    url                            = string
    active_directory_app_id_or_uri = string
    active_directory_tenant_id     = string
    event_token                    = string
  }))
}

variable "tags" {
  description = "The tags to apply to the topic"
  type        = map(string)
  nullable    = true
}

