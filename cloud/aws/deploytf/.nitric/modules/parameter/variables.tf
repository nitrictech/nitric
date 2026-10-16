variable "parameter_name" {
  description = "The name of the parameter"
  type        = string
}

variable "access_role_names" {
  description = "The names of the roles that can access the parameter"
  type        = set(string)
}

variable "parameter_value" {
  description = "The text value of the parameter"
  type        = string
}

variable "parameter_tier" {
  description = "The tier of the SSM parameter (Standard or Advanced). Standard supports up to 4,096 characters; Advanced supports up to 8,192 characters."
  type        = string
  default     = "Standard"

  validation {
    condition     = contains(["Standard", "Advanced"], var.parameter_tier)
    error_message = "parameter_tier must be one of: Standard, Advanced."
  }
}
