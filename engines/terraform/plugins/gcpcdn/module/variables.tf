variable "nitric" {
  type = object({
    name     = string
    stack_id = string
    origins = map(object({
      path = string
      type = string
      domain_name = string
      id = string
      resources = map(string)
    }))
  })
}

variable "project_id" {
  description = "The project ID where resources will be created."
  type        = string
}

variable "region" {
  description = "The region where resources will be created."
  type        = string
}

variable "cdn_domain" {
  description = "The CDN domain configuration."
  type = object({
    domain_name = string
    zone_name   = string
    domain_ttl  = optional(number, 300)
  })
}
