
# Create a random string for the stack id
resource "random_string" "stack_id" {
  length  = 8
  special = false
  lower   = true
  upper   = false
}

locals {
  stack_name = "${var.stack_name}-${random_string.stack_id.result}"
}

data "azurerm_client_config" "current" {}

# Create an azure resource group
resource "azurerm_resource_group" "resource_group" {
  name     = "${var.stack_name}-rg${random_string.stack_id.result}"
  location = var.location
  tags = {
    "x-nitric-${local.stack_name}-name" = var.stack_name
    "x-nitric-${local.stack_name}-type" = "stack"
  }
}

# Create an azure storage account
resource "azurerm_storage_account" "storage" {
  count               = var.enable_storage ? 1 : 0
  name                = "${var.stack_name}sa${random_string.stack_id.result}"
  resource_group_name = azurerm_resource_group.resource_group.name
  location            = azurerm_resource_group.resource_group.location
  account_tier        = "Standard"
  access_tier         = "Hot"
  # TODO: Make configurable  
  account_replication_type = "LRS"
  account_kind             = "StorageV2"
  tags = {
    "x-nitric-${local.stack_name}-name" = var.stack_name
    "x-nitric-${local.stack_name}-type" = "stack"
  }
}

# Create a keyvault if secrets are enabled
resource "azurerm_key_vault" "keyvault" {
  count = var.enable_keyvault ? 1 : 0

  name                       = "${var.stack_name}kv${random_string.stack_id.result}"
  resource_group_name        = azurerm_resource_group.resource_group.name
  location                   = azurerm_resource_group.resource_group.location
  sku_name                   = "standard"
  soft_delete_retention_days = 7
  tenant_id                  = data.azurerm_client_config.current.tenant_id
  enable_rbac_authorization  = true

  tags = {
    "x-nitric-${local.stack_name}-name" = var.stack_name
    "x-nitric-${local.stack_name}-type" = "stack"
  }
}

# Create a User assigned managed identity
resource "azurerm_user_assigned_identity" "managed_identity" {
  name                = "managed-identity-${local.stack_name}"
  resource_group_name = azurerm_resource_group.resource_group.name
  location            = azurerm_resource_group.resource_group.location
}

# Create a container registry for storing images
resource "azurerm_container_registry" "container_registry" {
  name                = "${var.stack_name}cr${random_string.stack_id.result}"
  resource_group_name = azurerm_resource_group.resource_group.name
  location            = azurerm_resource_group.resource_group.location
  sku                 = "Basic"
  admin_enabled       = true
  tags = {
    "x-nitric-${local.stack_name}-name" = var.stack_name
    "x-nitric-${local.stack_name}-type" = "stack"
  }
}

# Create an operational insights workspace
resource "azurerm_log_analytics_workspace" "log_analytics" {
  name                = "${var.stack_name}log${random_string.stack_id.result}"
  resource_group_name = azurerm_resource_group.resource_group.name
  location            = azurerm_resource_group.resource_group.location
  sku                 = "PerGB2018"
  retention_in_days   = 30
  tags = {
    "x-nitric-${local.stack_name}-name" = var.stack_name
    "x-nitric-${local.stack_name}-type" = "stack"
  }
}










# Create a virtual network for the database server
resource "azurerm_virtual_network" "database_network" {
  count               = var.enable_database ? 1 : 0
  name                = "nitric-database-vnet"
  resource_group_name = azurerm_resource_group.resource_group.name
  location            = var.location
  address_space       = ["10.0.0.0/16"]

  flow_timeout_in_minutes = 10
}

# Create a subnet for the database server
resource "azurerm_subnet" "database_subnet" {
  count                = var.enable_database ? 1 : 0
  name                 = "nitric-database-subnet"
  resource_group_name  = azurerm_resource_group.resource_group.name
  virtual_network_name = azurerm_virtual_network.database_network[0].name
  address_prefixes     = ["10.0.0.0/18"]

  service_endpoints = ["Microsoft.Storage"]

  delegation {
    name = "db-delegation"
    service_delegation {
      name    = "Microsoft.DBforPostgreSQL/flexibleServers"
      actions = ["Microsoft.Network/virtualNetworks/subnets/join/action"]
    }
  }
}

# Create an infrastructure subnet for the database server
resource "azurerm_subnet" "database_infrastructure_subnet" {
  count                = var.enable_database ? 1 : 0
  name                 = "nitric-database-infrastructure-subnet"
  resource_group_name  = azurerm_resource_group.resource_group.name
  virtual_network_name = azurerm_virtual_network.database_network[0].name
  address_prefixes     = ["10.0.64.0/18"]

  depends_on = [azurerm_subnet.database_subnet]
}

# Create a subnet for containers to connect to the database
resource "azurerm_subnet" "database_client_subnet" {
  count                = var.enable_database ? 1 : 0
  name                 = "nitric-database-client-subnet"
  resource_group_name  = azurerm_resource_group.resource_group.name
  virtual_network_name = azurerm_virtual_network.database_network[0].name
  address_prefixes     = ["10.0.192.0/18"]

  delegation {
    name = "container-instance-delegation"
    service_delegation {
      name    = "Microsoft.ContainerInstance/containerGroups"
      actions = ["Microsoft.Network/virtualNetworks/subnets/action"]
    }
  }

  depends_on = [azurerm_subnet.database_infrastructure_subnet]
}

# Create a private zone for the database server
resource "azurerm_private_dns_zone" "database_dns_zone" {
  count               = var.enable_database ? 1 : 0
  name                = "db-private-dns.postgres.database.azure.com"
  resource_group_name = azurerm_resource_group.resource_group.name
}

# Create a private link service for the database server
resource "azurerm_private_dns_zone_virtual_network_link" "database_link_service" {
  count                 = var.enable_database ? 1 : 0
  name                  = "nitric-database-link-service"
  private_dns_zone_name = azurerm_private_dns_zone.database_dns_zone[0].name
  resource_group_name   = azurerm_resource_group.resource_group.name
  virtual_network_id    = azurerm_virtual_network.database_network[0].id
  registration_enabled  = false

  tags = {
    "x-nitric-${var.stack_name}-name" = var.stack_name
    "x-nitric-${var.stack_name}-type" = "stack"
  }
}

# Create a random master password for the database server
resource "random_password" "database_master_password" {
  count   = var.enable_database ? 1 : 0
  length  = 16
  special = false
}

# Create a database service if required
resource "azurerm_postgresql_flexible_server" "database" {
  count                  = var.enable_database ? 1 : 0
  name                   = "nitric-db-${random_string.stack_id.result}"
  resource_group_name    = azurerm_resource_group.resource_group.name
  location               = var.location
  version                = "14"
  administrator_login    = "nitric"
  administrator_password = random_password.database_master_password[0].result

  zone = "1"

  public_network_access_enabled = false

  delegated_subnet_id = azurerm_subnet.database_subnet[0].id
  private_dns_zone_id = azurerm_private_dns_zone.database_dns_zone[0].id

  # default to 32Gb storage
  # TODO: Make configurable   
  storage_mb = 32768

  # TODO: Make configurable  
  sku_name = "B_Standard_B1ms"

  tags = {
    "x-nitric-${var.stack_name}-name" = var.stack_name
    "x-nitric-${var.stack_name}-type" = "stack"
  }

  depends_on = [
    azurerm_subnet.database_subnet,
    azurerm_private_dns_zone.database_dns_zone,
    azurerm_private_dns_zone_virtual_network_link.database_link_service
  ]
}















# Create a new container app managed environment
resource "azurerm_container_app_environment" "environment" {
  name                       = "${var.stack_name}kube${random_string.stack_id.result}"
  resource_group_name        = azurerm_resource_group.resource_group.name
  location                   = azurerm_resource_group.resource_group.location
  log_analytics_workspace_id = azurerm_log_analytics_workspace.log_analytics.id
  tags = {
    "x-nitric-${local.stack_name}-name" = var.stack_name
    "x-nitric-${local.stack_name}-type" = "stack"
  }

  infrastructure_subnet_id = azurerm_subnet.database_infrastructure_subnet[0].id
}
