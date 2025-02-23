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

variable "location" {
    description = "The location/region where the resources will be created"
    type        = string
}