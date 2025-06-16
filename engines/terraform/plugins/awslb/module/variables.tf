variable "load_balancer_type" {
  type    = string
  default = "application"
}

variable "name" {
  type    = string
}

variable "internal" {
  type    = bool
  default = false
}

variable "security_groups" {
  type    = list(string)
}

variable "subnets" {
  type    = list(string)
}
