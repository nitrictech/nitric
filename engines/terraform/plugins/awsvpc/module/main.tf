
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
  depends_on = [ module.vpc ]

  count = length(var.default_security_group_ingress)

  security_group_id = module.vpc.default_security_group_id
  type = "ingress"

  self             = lookup(element(var.default_security_group_ingress, count.index), "self", null)
  cidr_blocks      = compact(split(",", lookup(element(var.default_security_group_ingress, count.index), "cidr_blocks", "")))
  ipv6_cidr_blocks = compact(split(",", lookup(element(var.default_security_group_ingress, count.index), "ipv6_cidr_blocks", "")))
  prefix_list_ids  = compact(split(",", lookup(element(var.default_security_group_ingress, count.index), "prefix_list_ids", "")))
  description      = lookup(element(var.default_security_group_ingress, count.index), "description", null)
  from_port        = lookup(element(var.default_security_group_ingress, count.index), "from_port", 0)
  to_port          = lookup(element(var.default_security_group_ingress, count.index), "to_port", 0)
  protocol         = lookup(element(var.default_security_group_ingress, count.index), "protocol", "-1")
}

# Setup ingress on the container port for the security groups
resource "aws_security_group_rule" "egress" {
  depends_on = [ module.vpc ]

  count = length(var.default_security_group_egress)

  security_group_id = module.vpc.default_security_group_id
  type = "egress"

  self             = lookup(element(var.default_security_group_egress, count.index), "self", null)
  cidr_blocks      = compact(split(",", lookup(element(var.default_security_group_egress, count.index), "cidr_blocks", "")))
  ipv6_cidr_blocks = compact(split(",", lookup(element(var.default_security_group_egress, count.index), "ipv6_cidr_blocks", "")))
  prefix_list_ids  = compact(split(",", lookup(element(var.default_security_group_egress, count.index), "prefix_list_ids", "")))
  description      = lookup(element(var.default_security_group_egress, count.index), "description", null)
  from_port        = lookup(element(var.default_security_group_egress, count.index), "from_port", 0)
  to_port          = lookup(element(var.default_security_group_egress, count.index), "to_port", 0)
  protocol         = lookup(element(var.default_security_group_egress, count.index), "protocol", "-1")
}
