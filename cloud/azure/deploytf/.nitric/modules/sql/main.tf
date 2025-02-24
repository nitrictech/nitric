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
  charset             = "utf8"
  collation           = "en_US.utf8"
}

locals {
  count = var.migration_image != "" ? 1 : 0

  remote_image_name = "${var.image_registry_server}/${var.stack_name}-${var.name}:latest"

  db_url = "postgres://nitric:${var.database_master_password}@${var.database_server_fqdn}:5432/${var.name}"
}

# Tag the provided docker image with the ECR repository url
resource "docker_tag" "tag" {
  count = var.migration_image != "" ? 1 : 0

  source_image = var.migration_image
  target_image = local.remote_image_name
}

data "docker_image" "latest" {
  name = var.migration_image
}

# Push the tagged image to the ECR repository
resource "docker_registry_image" "push" {
  count = var.migration_image != "" ? 1 : 0

  name = local.remote_image_name

  triggers = {
    source_image_id = data.docker_image.latest.id
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

  ip_address_type = "Private"

  container {
    name   = "${var.name}-migration"
    # point to the pushed image sha256 digest to ensure container is updated when image changes
    image  = "${var.image_registry_server}/${var.stack_name}-${var.name}@${docker_registry_image.push[count.index].sha256_digest}"
    cpu    = 1
    memory = 1
    

    ports {
      port     = 80
      protocol = "TCP"
    }

    environment_variables = {
      DB_URL = local.db_url
      NITRIC_DB_NAME = var.name
    }
    
  }

  depends_on = [ docker_registry_image.push ]
}
