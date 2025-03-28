variable "stack_id" {
  description = "The ID of the stack."
  type        = string
}

variable "website_name" {
    description = "The name of the website."
    type        = string
}

variable "region" {
  description = "The region where the bucket will be created."
  type        = string
}

variable "base_path" {
  description = "The base path for the website files."
  type        = string
}

variable "local_directory" {
  description = "The local directory containing website files."
  type        = string
}

variable "index_document" {
  description = "The index document for the website."
  type        = string
  default     = "index.html"
}

variable "error_document" {
  description = "The error document for the website."
  type        = string
  default     = "404.html"
}