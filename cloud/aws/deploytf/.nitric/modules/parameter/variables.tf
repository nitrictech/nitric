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
