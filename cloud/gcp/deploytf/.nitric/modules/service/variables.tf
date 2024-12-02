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

variable "project_id" {
  description = "The ID of the Google Cloud project where the service is created."
  type        = string
}

variable "stack_id" {
  description = "The ID of the Nitric stack"
  type        = string
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

variable "region" {
  description = "The region the service is being deployed to"
  type        = string
}

variable "base_compute_role" {
  description = "The base compute role to use for the service"
  type        = string
}

variable "artifact_registry_repository" {
  description = "The base URI for the artifact registry repository the push this services image to"
  type        = string
}

variable "kms_key" {
  description = "The name of the KMS key to use"
  type        = string
  default     = ""
}

variable "vpc" {
  description = "cloud run service directvpc config"
  nullable = true
  type = object({
    network = string
    subnet = string
    network_tags = list(string)
  })
}
