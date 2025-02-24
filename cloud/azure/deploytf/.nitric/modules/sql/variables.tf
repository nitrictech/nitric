variable "name" {
  description = "The name of the database"
  type        = string
}

variable "stack_name" {
  description = "The name of the stack"
  type        = string
}

variable "location" {
  description = "The location/region the migration container should be deployed"
  type        = string
}

variable "resource_group_name" {
  description = "The name of the resource group"
  type        = string
}

variable "server_id" {
  description = "The id of the postgresql flexible server"
  type        = string
}

variable "migration_image" {
  description = "The migration image to use"
  type        = string
}

variable "migration_container_subnet_id" {
  description = "The subnet id to deploy the migration container"
  type        = string
}

variable "image_registry_server" {
  description = "The image registry server"
  type        = string
}

variable "image_registry_username" {
  description = "The image registry username"
  type        = string
}

variable "image_registry_password" {
  description = "The image registry password"
  type        = string
}

variable "database_server_fqdn" {
  description = "The database server fully qualified domain name"
  type = string
}

variable "database_master_password" {
  description = "The database master password"
  type = string
}