variable "stack_name" {
  description = "The name of the stack"
  type        = string
}

variable "storage_account_id" {
  description = "The id of the storage account to use for the cdn"
  type        = string
}

variable "storage_account_name" {
  description = "The name of the storage account to use for the cdn"
  type        = string
}

variable "storage_account_primary_web_host" {
  description = "The primary web host of the storage account to use for the cdn"
  type        = string
}

variable "resource_group_name" {
  description = "The name of the resource group to use for the cdn"
  type        = string
}

variable "apis" {
  description = "Map of APIs and their gateway information"
  type = map(object({
    gateway_url = string
  }))
  default = {}
}

variable "location" {
  description = "The location/region where the resources will be created"
  type        = string
}

variable "publisher_name" {
  description = "The name of the publisher"
  type        = string
}

variable "publisher_email" {
  description = "The email of the publisher"
  type        = string
}

# Variable to hold content paths to purge
variable "cdn_purge_paths" {
  description = "Map of content paths to purge from the CDN"
  type        = map(string)
  default     = {}
}