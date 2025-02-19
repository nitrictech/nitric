variable "name" {
    description = "The name of the service"
    type        = string
}

variable "stack_name" {
    description = "The name of the stack"
    type        = string
}

variable "resource_group_name" {
    description = "The name of the resource group"
    type        = string
  
}

variable "container_app_environment_id" {
    description = "The id of the container app environment"
    type        = string
}

variable "image_uri" {
    description = "The image uri for the container"
    type        = string
}

variable "cpu" {
    description = "The cpu limit for the container"
    type        = number
}

variable "registry_login_server" {
    description = "The login server for the container registry"
    type        = string
}

variable "registry_username" {
    description = "The username for the container registry"
    type        = string
}

variable "registry_password" {
    description = "The password for the container registry"
    type        = string
}

variable "memory" {
    description = "The memory limit for the container"
    type        = string
}

variable "env" {
    description = "The environment variables to set"
    type        = map(string)
}

variable "is_schedule" { 
    description = "Does this service contain a schedule"
    type        = bool
}

variable "min_replicas" {
  description = "Minimum number of replicas for the service"
  type        = number
}