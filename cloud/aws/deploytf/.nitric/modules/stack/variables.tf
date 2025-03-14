variable "enable_website" {
  description = "Enable the creation of a website"
  type        = bool
  default     = false
}

variable "website_root_index_document" {
  description = "The root index document for the website"
  type        = string
  default     = "index.html"
}

variable "website_root_error_document" {
  description = "The root error document for the website"
  type        = string
  default     = "404.html"
}
