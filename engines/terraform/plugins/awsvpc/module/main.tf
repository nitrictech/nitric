
# Get Azs via data
data "aws_availability_zones" "availability_zones" {}

locals {
  # Automatically get AZs if not provided
  azs = var.azs != null ? var.azs : data.aws_availability_zones.availability_zones.names
}

module "vpc" {
  source = "terraform-aws-modules/vpc/aws"

  name = var.name
  cidr = var.networking.cidr_block

  azs             = local.azs
  private_subnets = var.networking.private_subnets
  public_subnets  = var.networking.public_subnets

  single_nat_gateway = var.single_nat_gateway
  enable_nat_gateway = var.enable_nat_gateway
  enable_vpn_gateway = var.enable_vpn_gateway

  # DO NOT USE THESE, USE THE INGRESS AND EGRESS RULES BELOW
  # Otherwise the two can clobber each other
  # default_security_group_ingress = var.default_security_group_ingress
  # default_security_group_egress = var.default_security_group_egress

  tags = var.tags
}

# Setup ingress on the container port for the security groups
resource "aws_security_group_rule" "ingress" {
  count = length(var.default_security_group_ingress)

  security_group_id = module.vpc.default_security_group_id
  type = "egress"

  self             =  lookup(each.value, "self", null)
  cidr_blocks      = compact(split(",", lookup(each.value, "cidr_blocks", "")))
  ipv6_cidr_blocks = compact(split(",", lookup(each.value, "ipv6_cidr_blocks", "")))
  prefix_list_ids  = compact(split(",", lookup(each.value, "prefix_list_ids", "")))
  description      = lookup(each.value, "description", null)
  from_port        = lookup(each.value, "from_port", 0)
  to_port          = lookup(each.value, "to_port", 0)
  protocol         = lookup(each.value, "protocol", "-1")
}

# Setup ingress on the container port for the security groups
resource "aws_security_group_rule" "egress" {
  count = length(var.default_security_group_egress)

  security_group_id = module.vpc.default_security_group_id
  type = "egress"

  self             =  lookup(each.value, "self", null)
  cidr_blocks      = compact(split(",", lookup(each.value, "cidr_blocks", "")))
  ipv6_cidr_blocks = compact(split(",", lookup(each.value, "ipv6_cidr_blocks", "")))
  prefix_list_ids  = compact(split(",", lookup(each.value, "prefix_list_ids", "")))
  description      = lookup(each.value, "description", null)
  from_port        = lookup(each.value, "from_port", 0)
  to_port          = lookup(each.value, "to_port", 0)
  protocol         = lookup(each.value, "protocol", "-1")
}
