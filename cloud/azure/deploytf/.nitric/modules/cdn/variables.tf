variable "stack_name" {
  description = "The name of the stack"
  type        = string
}

variable "storage_account_primary_web_host" {
  description = "The primary web host of the storage account to use for the cdn"
  type        = string
}

variable "resource_group_name" {
  description = "The name of the resource group to use for the cdn"
  type        = string
  sensitive   = true
}

# Variable to hold content paths to purge
variable "cdn_purge_paths" {
  description = "Map of content paths to purge from the CDN"
  type        = map(string)
  default     = {}
}

variable "enable_api_rewrites" {
  description = "Enable API rewrites"
  type        = bool
  default     = false
}