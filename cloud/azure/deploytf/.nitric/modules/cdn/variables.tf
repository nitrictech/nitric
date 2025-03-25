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

variable "zone_resource_group_name" {
  description = "The name of the resource group for the CDN zone"
  type        = string
}

variable "zone_name" {
  description = "The name of the CDN zone"
  type        = string
}

variable "enable_custom_domain" {
  description = "Enable custom domains"
  type        = bool
  default     = false
}

variable "domain_name" {
  description = "The domain name for the CDN"
  type        = string
}

variable "custom_domain_host_name" {
  description = "Custom domain host name"
  type        = string
}

variable "is_apex_domain" {
  description = "Is the custom domain an apex domain"
  type        = bool
  default     = false
}