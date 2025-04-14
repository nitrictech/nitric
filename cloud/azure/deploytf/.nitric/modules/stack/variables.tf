variable "stack_name" {
  description = "The name of the stack"
  type        = string
}

variable "enable_storage" {
  description = "Enable the creation of a storage account"
  type        = bool
  default     = false
}

variable "enable_keyvault" {
  description = "Enable the creation of a keyvault"
  type        = bool
  default     = false
}

variable "enable_database" {
  description = "Enable the creation of a database"
  type        = bool
  default     = false
}

variable "resource_group_name" {
  description = "The name of the resource group to reuse"
  type        = string
  nullable    = true
  default     = null
}

variable "location" {
  description = "The location/region where the resources will be created"
  type        = string
}

variable "vnet_name" {
  description = "The name of the vnet to deploy the infrastructure resources"
  type        = string
  nullable    = true
}

variable "subnet_id" {
  description = "The id of the subnet to deploy the infrastructure resources"
  type        = string
  nullable    = true
}

variable "create_dns_zones" {
  description = "Whether to create private DNS zones for private endpoints"
  type        = bool
  default     = false
}

variable "private_endpoints" {
  description = "Deploy compatible services with private endpoints"
  type        = bool
  nullable    = true
}

variable "tags" {
  description = "The tags to apply to the stack"
  type        = map(string)
  nullable    = true
}
