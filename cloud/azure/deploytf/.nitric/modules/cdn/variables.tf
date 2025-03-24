variable "stack_name" {
  description = "The name of the stack"
  type        = string
}

variable "resource_group_name" {
  description = "The name of the resource group to use for the cdn"
  type        = string
  sensitive   = true
}

variable "uploaded_files" {
  description = "Map of uploaded files with their MD5 hashes"
  type        = map(string)
  default     = {}
}

variable "primary_web_host" {
  description = "The primary host for the CDN"
  type = string
}

variable "enable_api_rewrites" {
  description = "Enable API rewrites"
  type        = bool
  default     = false
}