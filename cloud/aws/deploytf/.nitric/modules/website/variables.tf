variable "name" {
  type = string
  description = "The name of the website"
}

variable "base_path" {
  type = string
  description = "The base path for the website"
}

variable "local_directory" {
  type = string
  description = "The production website output directory"
}

variable "website_bucket_id" {
  type = string
  description = "The id of the bucket to deploy the website files to"
}