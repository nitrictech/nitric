output "cluster_endpoint" {
    description = "The endpoint of the RDS cluster"
    value       = aws_rds_cluster.cluster.endpoint
}

output "cluster_username" {
    description = "The username to connect to the RDS cluster"
    value       = aws_rds_cluster.cluster.master_username
}

output "cluster_password" {
    description = "The password to connect to the RDS cluster"
    value       = aws_rds_cluster.cluster.master_password
}

output "security_group_id" {
    description = "The security group id for the RDS cluster"
    value       = aws_security_group.rds_security_group.id
}

output "create_database_project_name" {
    description = "The name of the create database codebuild project"
    value       = aws_codebuild_project.create_database.name
}