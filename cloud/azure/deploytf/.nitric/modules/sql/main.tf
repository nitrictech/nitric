# Create a new postgresql database
resource "azurerm_postgresql_database" "db" {
  name                = var.name
  resource_group_name = var.resource_group_name
  server_name         = var.server_name
  charset             = "UTF8"
  collation           = "en_US.UTF8"
}

# Push the migration image
data "docker_image" "latest" {
  name = var.image_uri
}

locals {
  remote_image_name = "${var.image_registry_server}/${var.stack_name}-${var.name}"
}

# Tag the provided docker image with the ECR repository url
resource "docker_tag" "tag" {
  source_image = data.docker_image.latest.repo_digest
  target_image = local.remote_image_name
}

# Push the tagged image to the ECR repository
resource "docker_registry_image" "push" {
  name = local.remote_image_name
  triggers = {
    source_image_id = docker_tag.tag.source_image_id
  }
}


# Create a new azure container instances to execute the migration
resource "azurerm_container_group" "migration" {
  name                = "${var.name}-migration"
  location            = var.location
  resource_group_name = var.resource_group_name
  os_type             = "Linux"
  dns_name_label      = "${var.name}-migration"
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
      DB_URL = azurerm_postgresql_server.db_flexible_connection_strings[0].value
      NITRIC_DB_NAME = var.name
    }
  }
}
