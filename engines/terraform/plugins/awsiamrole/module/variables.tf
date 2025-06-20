variable "nitric" {
  type = object({
    name = string
    stack_id = string
  })
}

variable "trusted_services" {
  type    = list(string)
  default = ["lambda.amazonaws.com", "ec2.amazonaws.com", "ecs-tasks.amazonaws.com"]
}

variable "trusted_actions" {
  type    = list(string)
  default = ["sts:AssumeRole"]
}
