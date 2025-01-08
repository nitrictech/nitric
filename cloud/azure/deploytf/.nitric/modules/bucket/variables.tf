variable "name" {
  description = "The name of the bucket"
  type        = string
}

variable "stack_name" {
  description = "The name of the stack"
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
  }))
}
