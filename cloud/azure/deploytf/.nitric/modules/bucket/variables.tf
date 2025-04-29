variable "name" {
  description = "The name of the bucket"
  type        = string
}

variable "storage_account_id" {
  description = "The id of the storage account"
  type        = string
}

variable "listeners" {
  description = "The list of listeners to notify"
  type = map(object({
    url                            = string
    active_directory_app_id_or_uri = string
    active_directory_tenant_id     = string
    event_token                    = string
    event_type                     = list(string)
  }))
}

variable "tags" {
  description = "The tags to apply to the bucket"
  type        = map(string)
  nullable    = true
}
