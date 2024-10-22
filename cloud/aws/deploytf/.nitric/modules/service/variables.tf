variable "service_name" {
  type        = string
  description = "The name of the service"
}

variable "image" {
  type        = string
  description = "The docker image to deploy"
}

# environment variables
variable "environment" {
  type        = map(string)
  description = "Environment variables to set on the lambda function"
}

variable "stack_id" {
  description = "The ID of the Nitric stack"
  type        = string
}

variable "memory" {
  description = "The amount of memory to allocate to the lambda function"
  type        = number
  default     = 128
}

variable "ephemeral_storage" {
  description = "The amount of ephemeral storage to allocate to the lambda function"
  type        = number
  default     = 512
}

variable "timeout" {
  description = "The timeout for the lambda function"
  type        = number
  default     = 15
}

variable "subnet_ids" {
  description = "The subnet ids to use for the aws lambda function"
  type        = list(string)
  default     = []
}

variable "security_group_ids" {
  description = "The security group ids to use for the aws lambda function"
  type        = list(string)
  default     = []

}


