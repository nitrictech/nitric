variable "name" {
  description = "The name of the API"
  type        = string
}

variable "description" {
  description = "The description of the API"
  type        = string
}

variable "location" {
  description = "The location of the API"
  type        = string
}

variable "app_identity" {
  description = "The identity of the app"
  type        = string
}

variable "openapi_spec" {
  description = "The openapi spec to deploy"
  type        = string
}

variable "publisher_name" {
  description = "The name of the publisher"
  type        = string
}

variable "publisher_email" {
  description = "The email of the publisher"
  type        = string
}

variable "resource_group_name" {
  description = "The name of the resource group"
  type        = string
}

# TODO: May be able to apply these directly in the terraform instead and
# Just supply target services to route to
variable "operation_policy_templates" {
  description = "The policy templates to apply"
  type        = map(string)
}

variable "tags" {
  description = "The tags to apply to the API"
  type        = map(string)
  nullable    = true
}

