variable "nitric" {
  type = object({
    name       = string
    stack_id   = string
    image_id   = string
    env        = map(string)
    identities = map(any)
  })
}

variable "environment" {
    type = map(string)
    description = "Environment variables to set on the lambda function"
    default     = {}
}
# TODO: review defaults
variable "memory_mb" {
    description = "The amount of memory to allocate to the CloudRun service in MB"
    type        = number
    default     = 512
}

variable "cpus" {
    description = "The amount of cpus to allocate to the CloudRun service"
    type        = number
    default     = 1
}

variable "gpus" {
    description = "The amount of gpus to allocate to the CloudRun service"
    type        = number
    default     = 0
}

variable "min_instances" {
    description = "The minimum number of instances to run"
    type        = number
    default     = 0
}

variable "max_instances" {
    description = "The maximum number of instances to run"
    type        = number
    default     = 10
}

variable "container_concurrency" {
    description = "The number of concurrent requests the CloudRun service can handle"
    type        = number
    default     = 80
}

variable "timeout_seconds" {
    description = "The timeout for the CloudRun service in seconds"
    type        = number
    default     = 10
}

variable "project_id" {
    description = "The project ID to deploy the CloudRun service to"
    type        = string
}

variable "region" {
    description = "The region to deploy the CloudRun service to"
    type        = string
}

variable "ingress_port" {
    description = "The port to expose the CloudRun service to"
    type        = number
    default     = 9001
}