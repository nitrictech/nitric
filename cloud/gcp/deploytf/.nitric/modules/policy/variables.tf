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
  type        = string
}

variable "actions" {
  description = "The actions to apply to the policy"
  type        = list(string)
}

variable "iam_roles" {
  description = "The IAM roles available to the policy"
  type = object({
    base_compute_role = string
    bucket_delete     = string
    bucket_read       = string
    bucket_write      = string
    kv_delete         = string
    kv_read           = string
    kv_write          = string
    queue_dequeue     = string
    queue_enqueue     = string
    secret_access     = string
    secret_put        = string
    topic_publish     = string
  })
}
