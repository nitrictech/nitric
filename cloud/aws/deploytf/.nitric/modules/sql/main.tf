# tag the local image uri
terraform {
  required_providers {
    docker = {
      source  = "kreuzwerker/docker"
      version = "3.0.2"
    }
  }
}

# Create an ECR repository
# If there is a migration image uri tag it and push it to ECR
resource "aws_ecr_repository" "migrate_database" {
  name = "nitric-migrate-database-${var.db_name}"
}

data "aws_ecr_authorization_token" "ecr_auth" {
}

provider "docker" {
  registry_auth {
    address  = data.aws_ecr_authorization_token.ecr_auth.proxy_endpoint
    username = data.aws_ecr_authorization_token.ecr_auth.user_name
    password = data.aws_ecr_authorization_token.ecr_auth.password
  }
}

# Tag the provided docker image with the ECR repository url
resource "docker_tag" "tag" {
  source_image = var.image_uri
  target_image = aws_ecr_repository.migrate_database.repository_url
}

# Push the tagged image to the ECR repository
resource "docker_registry_image" "push" {
  name = aws_ecr_repository.migrate_database.repository_url
  triggers = {
    source_image_id = docker_tag.tag.source_image_id
  }
}

locals {
  db_url = "postgres://${var.rds_cluster_username}:${var.rds_cluster_password}@${var.rds_cluster_endpoint}/${var.db_name}"
}

# Create an AWS codebuild project for running migrations
resource "aws_codebuild_project" "migrate_database" {
  name         = "nitric-migrate-database-${var.db_name}"
  description  = "Migrate the database on the RDS cluster"
  service_role = var.codebuild_role_arn

  artifacts {
    type = "NO_ARTIFACTS"
  }

  environment {
    compute_type = "BUILD_GENERAL1_SMALL"
    image        = "${aws_ecr_repository.migrate_database.repository_url}@${docker_registry_image.push.sha256_digest}"
    type         = "LINUX_CONTAINER"
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
            "cd ${var.work_dir}",
            var.migrate_command
          ]
        }
      }
    })
  }
}

# Execute the create database codebuild job and wait for it to complete
resource "null_resource" "execute_create_database" {
  provisioner "local-exec" {
    command = <<EOF
      BUILD_ID=$(aws codebuild start-build --project-name ${var.create_database_project_name} --query 'build.id' --output text)
      STATUS="IN_PROGRESS"
      while [[ $STATUS == "IN_PROGRESS" ]]; do
        sleep 5
        STATUS=$(aws codebuild batch-get-builds --ids $BUILD_ID | jq -r '.builds[0].buildStatus')
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
  provisioner "local-exec" {
    command = <<EOF
      BUILD_ID=$(aws codebuild start-build --project-name ${aws_codebuild_project.migrate_database.name} --query 'build.id' --output text)
      STATUS="IN_PROGRESS"
      while [[ $STATUS == "IN_PROGRESS" ]]; do
        sleep 5
        STATUS=$(aws codebuild batch-get-builds --ids $BUILD_ID | jq -r '.builds[0].buildStatus')
      done
      if [[ $STATUS != "SUCCEEDED" ]]; then
        echo "Build failed with status $STATUS"
        exit 1
      fi
    EOF
  }
  # Create the database first   
  depends_on = [null_resource.execute_create_database]
}
