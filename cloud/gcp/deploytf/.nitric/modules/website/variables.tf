variable "region" {
  description = "The region where the bucket will be created."
  type        = string
}

variable "bucket_name" {
  description = "The name of the bucket."
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