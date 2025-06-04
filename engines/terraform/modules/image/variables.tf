variable "tag" {
  type        = string
  description = "The tag to use for the build"
}

variable "dockerfile" {
  type        = string
  description = "The dockerfile to use for the build"
  default     = "Dockerfile"
}

variable "image_id" {
  type        = string
  description = "An existing image id to use for the build"
  default     = null
  nullable    = true
}

variable "build_context" {
  type        = string
  description = "The context for the build"
  default     = null
  nullable    = true
}

variable "platform" {
  type    = string
  default = "linux/amd64"
}

variable "args" {
  type        = map(string)
  description = "The arguments to pass to the build"
  default     = {}
}
