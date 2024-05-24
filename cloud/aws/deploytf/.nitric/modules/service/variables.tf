variable "service_name" {
    type = string
    description = "The name of the service"
}

variable "image" {
    type = string
    description = "The docker image to deploy"
}

# environment variables
variable "environment" {
    type = map(string)
    description = "Environment variables to set on the lambda function"
}

variable "stack_id" {
  description = "The ID of the Nitric stack"
  type        = string
}
