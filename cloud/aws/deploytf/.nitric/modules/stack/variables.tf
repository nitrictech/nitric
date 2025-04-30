variable "enable_website" {
  description = "Enable the creation of a website"
  type        = bool
  default     = false
}

variable "project_name" {
  description = "The name of the project"
  type        = string
}

variable "stack_name" {
  description = "The name of the stack"
  type        = string
}
