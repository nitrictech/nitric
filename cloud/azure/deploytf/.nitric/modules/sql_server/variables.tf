variable "stack_name" {
    description = "The name of the stack"
    type        = string
}

variable "stack_id" {
    description = "The random id generated for this stack"
    type = string
}

variable "location" {
    description = "The location/region where the resources will be created"
    type        = string
}

variable "resource_group_name" {
  description = "The name of the resource group"
  type        = string
}