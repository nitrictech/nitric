# Variables for configuration
variable "region" {
  description = "The region where resources will be created."
  type        = string
}

variable "stack_id" {
  description = "The unique identifier for the stack."
  type        = string
}

variable "website_buckets" {
  description = "A map of website bucket configurations."
  type = map(object({
    name          = string
    index_document = string
    error_document = string
    local_directory = string
  }))
}

variable "api_gateways" {
  description = "A map of API gateway configurations."
  type = map(object({
    region       = string
    gateway_id   = string
    default_host = string
  }))
}

variable "cdn_domain" {
  description = "The CDN domain configuration."
  type = object({
    domain_name = string
    zone_name   = string
  })
}