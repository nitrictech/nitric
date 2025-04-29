variable "stack_id" {
  description = "The id of the stack"
  type        = string
}

variable "base_path" {
  description = "The base path for the website"
  type        = string
}

variable "local_directory" {
  description = "The local directory to deploy the website from"
  type        = string
}

variable "resource_group_name" {
  description = "The name of the resource group to use for the cdn"
  type        = string
  sensitive   = true
}

variable "location" {
  description = "The location/region where the resources will be created"
  type        = string
}

variable "index_document" {
  description = "The index document for the website"
  type        = string
  default     = "index.html"
}

variable "error_document" {
  description = "The error document for the website"
  type        = string
  default     = "404.html"
}