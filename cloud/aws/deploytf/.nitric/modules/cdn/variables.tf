variable "stack_name" {
  description = "The name of the stack"
  type        = string
}

variable "website_bucket_id" {
  description = "The ID for the website bucket"
  type = string
}

variable "website_bucket_arn" {
  description = "The ARN for the website bucket"
  type = string
}

variable "website_bucket_domain_name" {
  description = "The domain name for the website bucket"
  type        = string
}

variable "website_index_document" {
  description = "The website index document"
  type = string
  default = "index.html"
}

variable "website_error_document" {
  description = "The website error document"
  type = string
  default = "404.html"
}

variable "apis" {
  description = "Map of APIs and their gateway information"
  type = map(object({
    gateway_url = string
    # Add any other API properties you might need
  }))
  default = {}
}