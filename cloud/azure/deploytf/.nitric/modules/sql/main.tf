terraform {
  required_providers {
    docker = {
      source  = "kreuzwerker/docker"
      version = "3.0.2"
    }
  }
}

# Create a new postgresql database
resource "azurerm_postgresql_flexible_server_database" "db" {
  name                = var.name
  server_id           = var.server_id
  charset             = "UTF8"
  collation           = "en_US.UTF8"
}

# Push the migration image
data "docker_image" "latest" {
  count = var.migration_image != "" ? 1 : 0

  name = var.migration_image
}

locals {
  count = var.migration_image != "" ? 1 : 0

  remote_image_name = "${var.image_registry_server}/${var.stack_name}-${var.name}"
}

# Tag the provided docker image with the ECR repository url
resource "docker_tag" "tag" {
  count = var.migration_image != "" ? 1 : 0

  source_image = data.docker_image.latest[0].repo_digest
  target_image = local.remote_image_name
}

# Push the tagged image to the ECR repository
resource "docker_registry_image" "push" {
  count = var.migration_image != "" ? 1 : 0

  name = local.remote_image_name
  triggers = {
    source_image_id = docker_tag.tag[0].source_image_id
  }
}

# Create a new azure container instances to execute the migration
resource "azurerm_container_group" "migration" {
  count = var.migration_image != "" ? 1 : 0

  name                = "${var.name}-migration"
  location            = var.location
  resource_group_name = var.resource_group_name
  os_type             = "Linux"
  restart_policy      = "Never"
  sku = "Standard"

  subnet_ids = [ var.migration_container_subnet_id ]

  image_registry_credential {
    server = var.image_registry_server
    username = var.image_registry_username
    password = var.image_registry_password
  }
  container {
    name   = "${var.name}-migration"
    image  = local.remote_image_name
    cpu    = 1
    memory = 1

    environment_variables = {
      DB_URL = var.database_server_fqdn
      NITRIC_DB_NAME = var.name
    }
  }
}
