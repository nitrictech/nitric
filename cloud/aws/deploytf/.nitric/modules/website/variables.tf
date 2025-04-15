variable "name" {
  type = string
  description = "The name of the website"
}

variable "stack_id" {
  type = string
  description = "The unique ID for this stack"
}

variable "base_path" {
  type = string
  description = "The base path for the website"
}

variable "local_directory" {
  type = string
  description = "The production website output directory"
}