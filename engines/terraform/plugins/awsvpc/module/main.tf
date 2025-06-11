
# Get Azs via data
data "aws_availability_zones" "availability_zones" {}

locals {
  # Automatically get AZs if not provided
  azs = var.azs != null ? var.azs : data.aws_availability_zones.availability_zones.names
  vpc_name = "${var.nitric.stack_id}-${var.nitric.name}"
}

module "vpc" {
  source = "terraform-aws-modules/vpc/aws"

  name = local.vpc_name
  cidr = var.networking.cidr_block

  azs             = local.azs
  private_subnets = var.networking.private_subnets
  public_subnets  = var.networking.public_subnets

  single_nat_gateway = var.single_nat_gateway
  enable_nat_gateway = var.enable_nat_gateway
  enable_vpn_gateway = var.enable_vpn_gateway

  tags = var.tags
}