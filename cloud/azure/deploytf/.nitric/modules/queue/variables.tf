variable "name" {
  description = "the name of the queue"
  type = string
}

variable "storage_account_name" {
  description = "the name of the storage account"
  type = string
}

variable "tags" {
  description = "the tags to apply to the queue"
  type        = map(string)
  nullable    = true
}