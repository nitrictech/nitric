variable "nitric" {
  type = object({
    name     = string
    stack_id = string
    content_path = string
    services = map(object({
      actions = list(string)
      identities = map(object({
        id   = string
        role = any
      }))
    }))
  })
}

variable "project_id" {
  description = "The project ID where resources will be created."
  type        = string
}

variable "region" {
  description = "The region where resources will be created."
  type        = string
}

variable "storage_class" {
  description = "The class of storage used to store the bucket's contents. This can be STANDARD, NEARLINE, COLDLINE, ARCHIVE, or MULTI_REGIONAL."
  type        = string
  default     = "STANDARD"
}