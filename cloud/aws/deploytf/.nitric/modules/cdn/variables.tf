variable "stack_name" {
  description = "The name of the stack"
  type        = string
}

variable "root_website" {
  description = "information about the root website for default behaviour"
  type = object({
    name = string
    index_document = optional(string, "index.html")
    error_document = optional(string, "404.html")
  })
}

variable "websites" {
  description = "Map of websites and their storage information"
  type = map(object({
    bucket_domain_name = string
    bucket_arn = string
    bucket_id = string
    base_path = string
    changed_files = list(string)
  }))
}

variable "apis" {
  description = "Map of APIs and their gateway information"
  type = map(object({
    gateway_url = string
    # Add any other API properties you might need
  }))
  default = {}
}

variable "domain_name" {
  description = "Custom domain for distribution"
  type = string
  default = ""
}

variable "certificate_arn" {
  description = "Certificate ARN for us-east-1 specific certificate"
  type = string
  default = ""
}

variable "zone_id" {
  description = "The ID of the hosted zone to store route53 records"
  type = string
  default = ""
}

variable "skip_cache_invalidation" {
  description = "Skip invalidating the cache. Defaults to false."
  type = bool
  default = false
}