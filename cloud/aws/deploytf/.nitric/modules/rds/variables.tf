variable "vpc_id" {
  type        = string
  description = "the VPC to assign to the RDS cluster"
}

variable "private_subnet_ids" {
  type        = list(string)
  description = "private subnets to assign to the RDS cluster"
}

variable "min_capacity" {
  type        = number
  description = "the minimum capacity of the RDS cluster"
}

variable "max_capacity" {
  type        = number
  description = "the maximum capacity of the RDS cluster"
}

variable "stack_id" {
  type        = string
  description = "The nitric stack ID"
}
