terraform {
  required_providers {
    docker = {
      source  = "kreuzwerker/docker"
      version = "3.0.2"
    }
  }
}

# Create an ECR repository for the migration image
# If there is a migration image uri tag it and push it to ECR
resource "aws_ecr_repository" "migrate_database" {
  # Only do this if an image uri is provided
  count = var.migrations != null ? 1 : 0
  name  = "nitric-migrate-database-${var.db_name}"
}

# Tag the provided docker image with the ECR repository url
resource "docker_tag" "migrate_tag" {
  count        = length(aws_ecr_repository.migrate_database)
  source_image = var.migrations.image_uri
  target_image = aws_ecr_repository.migrate_database[count.index].repository_url
}

# Push the tagged image to the ECR repository
resource "docker_registry_image" "migrate_image_push" {
  count = length(aws_ecr_repository.migrate_database)
  name  = aws_ecr_repository.migrate_database[count.index].repository_url
  triggers = {
    source_image_id = docker_tag.migrate_tag[count.index].source_image_id
  }
}

locals {
  db_url = "postgres://${var.rds_cluster_username}:${var.rds_cluster_password}@${var.rds_cluster_endpoint}/${var.db_name}"
}

# Create an AWS codebuild project for running migrations
resource "aws_codebuild_project" "migrate_database" {
  count        = length(aws_ecr_repository.migrate_database)
  name         = "nitric-migrate-database-${var.db_name}"
  description  = "Migrate the database on the RDS cluster"
  service_role = var.codebuild_role_arn

  artifacts {
    type = "NO_ARTIFACTS"
  }

  environment {
    compute_type                = "BUILD_GENERAL1_SMALL"
    image                       = "${aws_ecr_repository.migrate_database[count.index].repository_url}@${docker_registry_image.migrate_image_push[count.index].sha256_digest}"
    image_pull_credentials_type = "SERVICE_ROLE"
    type                        = "LINUX_CONTAINER"
    environment_variable {
      name  = "DB_NAME"
      value = var.db_name
    }
    environment_variable {
      name  = "DB_URL"
      value = local.db_url
    }
  }

  vpc_config {
    subnets            = var.subnet_ids
    security_group_ids = var.security_group_ids
    vpc_id             = var.vpc_id
  }

  source {
    type = "NO_SOURCE"
    buildspec = jsonencode({
      version = "0.2",
      phases = {
        build = {
          commands = [
            "echo 'Migrating database $DB_NAME'",
            "cd ${var.migrations.work_dir}",
            var.migrations.migrate_command
          ]
        }
      }
    })
  }
}

# Execute the create database codebuild job and wait for it to complete
resource "null_resource" "execute_create_database" {
  provisioner "local-exec" {
    interpreter = ["bash", "-c"]
    command     = <<EOF
      BUILD_ID=$(aws codebuild start-build --project-name ${var.create_database_project_name} --region ${var.codebuild_region}  --environment-variables "name=DB_NAME,value=${var.db_name}" --query 'build.id' --output text)
      STATUS="IN_PROGRESS"
      while [[ $STATUS == "IN_PROGRESS" ]]; do
        sleep 5
        STATUS=$(aws codebuild batch-get-builds --region ${var.codebuild_region} --ids $BUILD_ID --query 'builds[0].buildStatus' --output text)
      done
      if [[ $STATUS != "SUCCEEDED" ]]; then
        echo "Build failed with status $STATUS"
        exit 1
      fi
    EOF
  }
}

# Execute the codebuild job using the aws-cli and wait for it to complete
resource "null_resource" "execute_migrate_database" {
  count = length(aws_ecr_repository.migrate_database)
  provisioner "local-exec" {
    interpreter = ["bash", "-c"]
    command     = <<EOF
      # Start the database migrations
      BUILD_ID=$(aws codebuild start-build --project-name ${aws_codebuild_project.migrate_database[count.index].name} --region ${var.codebuild_region} --query 'build.id' --output text)
      STATUS="IN_PROGRESS"
      while [[ $STATUS == "IN_PROGRESS" ]]; do
        sleep 5
        # Check database migration status
        STATUS=$(aws codebuild batch-get-builds --region ${var.codebuild_region} --ids $BUILD_ID --query 'builds[0].buildStatus' --output text)
      done
      if [[ $STATUS != "SUCCEEDED" ]]; then
        echo "Database migrations failed $STATUS"
        exit 1
      fi
    EOF
  }
  triggers = {
    image_id = docker_registry_image.migrate_image_push[count.index].id
  }
  # Create the database first   
  depends_on = [docker_registry_image.migrate_image_push, aws_codebuild_project.migrate_database, null_resource.execute_create_database]
}
