output "vpc_id" {
    value = module.vpc.vpc_id
    description = "ID of the created VPC"
}

output "private_subnets" {
    value = module.vpc.private_subnets
    description = "Private subnet IDs"
}

output "public_subnets" {
    value = module.vpc.public_subnets
    description = "Public subnet IDs"
}

output "subnets" {
    value = concat(module.vpc.private_subnets, module.vpc.public_subnets)
    description = "All subnet IDs"
}

output "default_security_group_id" {
    value = module.vpc.default_security_group_id
    description = "Default security group ID"
}