variable "name" {
  description = "The name of the website"
  type        = string
}

variable "base_path" {
  description = "The base path for the website"
  type        = string
}

variable "local_directory" {
  description = "The local directory to deploy the website from"
  type        = string
}

variable "storage_account_name" {
  description = "The name of the storage account"
  type        = string
}

variable "storage_account_connection_string" {
  description = "The connection string for the storage account"
  type        = string
}
