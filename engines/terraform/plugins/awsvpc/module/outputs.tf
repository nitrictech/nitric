output "vpc_id" {
    value = module.vpc.vpc_id
    description = "ID of the created VPC"
}

output "private_subnet" {
    value = module.vpc.private_subnets
    description = "Private subnet IDs"
}

output "public_subnets" {
    value = module.vpc.public_subnets
    description = "Public subnet IDs"
}