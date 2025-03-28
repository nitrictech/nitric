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
