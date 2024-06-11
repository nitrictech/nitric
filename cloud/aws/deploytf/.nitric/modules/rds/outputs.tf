output "cluster_endpoint" {
    description = "The endpoint of the RDS cluster"
    value       = aws_rds_cluster.rds_cluster.endpoint
}

output "cluster_username" {
    description = "The username to connect to the RDS cluster"
    value       = aws_rds_cluster.rds_cluster.master_username
}

output "cluster_password" {
    description = "The password to connect to the RDS cluster"
    value       = aws_rds_cluster.rds_cluster.master_password
}

output "codebuild_role_arn" {
    description = "The arn of the codebuild role"
    value       = aws_iam_role.codebuild_role.arn
}

output "security_group_id" {
    description = "The security group id for the RDS cluster"
    value       = aws_security_group.rds_security_group.id
}

output "create_database_project_name" {
    description = "The name of the create database codebuild project"
    value       = aws_codebuild_project.create_database.name
}