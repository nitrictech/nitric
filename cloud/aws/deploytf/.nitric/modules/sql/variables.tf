variable "db_name" {
  description = "The name of the database to create"
  type        = string
}

variable "image_uri" {
  description = "The URI of the docker image to use for the codebuild project"
  type        = string
}

variable "rds_cluster_endpoint" {
  description = "The endpoint of the RDS cluster to connect to"
  type        = string
}

variable "rds_cluster_username" {
  description = "The username to connect to the RDS cluster"
  type        = string
}

variable "rds_cluster_password" {
  description = "The password to connect to the RDS cluster"
  type        = string
}

variable "security_group_ids" {
  description = "The security group ids to use for the codebuild project"
  type        = list(string)
}

variable "subnet_ids" {
  description = "The subnet ids to use for the codebuild project"
  type        = list(string)
}

variable "vpc_id" {
  description = "The vpc id to use for the codebuild project"
  type        = string
}

variable "work_dir" {
  description = "The working directory for the codebuild project"
  type        = string
}

variable "codebuild_role_arn" {
  description = "The arn of the codebuild role"
  type        = string
}

variable "codebuild_region" {
  description = "The region of the codebuild project"
  type        = string
}

variable "migrate_command" {
  description = "The command to run to migrate the database"
  type        = string
}

variable "create_database_project_name" {
  description = "The name of the create database codebuild project"
  type        = string
}
