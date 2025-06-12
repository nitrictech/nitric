variable "nitric" {
  type = object({
    name       = string
    stack_id   = string
    image_id   = string
    env        = map(string)
    identities = map(any)
  })
}

variable "container_port" {
  type    = number
}

variable "alb_arn" {
  type    = string
}

variable "cpu" {
  type    = number
  default = 1024
}

variable "memory" {
  type    = number
  default = 1024
}

variable "environment" {
  type    = map(string)
  default = {}
}

variable "vpc_id" {
  type    = string
  default = null
}

variable "subnets" {
  type    = list(string)
  default = []
}

variable "security_groups" {
  type    = list(string)
  default = []
}
