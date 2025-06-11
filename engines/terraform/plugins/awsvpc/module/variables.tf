variable "nitric" {
    type = object({
        name = string
        stack_id = string
    })
}

variable "networking" {
    type = object({
        cidr_block = string
        private_subnets = list(string)
        public_subnets = list(string)
    })
    default = {
        cidr_block = "10.0.0.0/16"
        private_subnets = ["10.0.1.0/24", "10.0.2.0/24", "10.0.3.0/24"]
        public_subnets = ["10.0.101.0/24", "10.0.102.0/24", "10.0.103.0/24"]
    }
}
variable "azs" {
    type = list(string)
    nullable = true
    default = null
}

variable "enable_nat_gateway" {
    type = bool
    default = false
}
variable "enable_vpn_gateway" {
    type = bool
    default = false
}

variable "single_nat_gateway" {
    type = bool
    default = false
}

variable "tags" {
    type = map(string)
}