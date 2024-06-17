variable "project_id" {
  description = "The google project id"
  type        = string
}

variable "resource_type" {
  description = "The type of the resource (Bucket, Secret, KeyValueStore, Queue)"
  type        = string
}

variable "resource_name" {
  description = "The name of the resource"
  type        = string
}

variable "service_account_email" {
  description = "The service account to apply the policy to"
  type = string
}

variable "actions" {
  description = "The actions to apply to the policy"
  type        = list(string)
}