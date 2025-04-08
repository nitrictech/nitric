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
  tags = merge(var.tags, {
    "x-nitric-${local.stack_name}-name" = var.stack_name
    "x-nitric-${local.stack_name}-type" = "stack"
  })
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

  network_rules {
    default_action = var.enable_storage_private_endpoints ? "Deny" : "Allow"
    bypass         = ["AzureServices"]
  }

  tags = merge(var.tags, {
    "x-nitric-${local.stack_name}-name" = var.stack_name
    "x-nitric-${local.stack_name}-type" = "stack"
  })
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

  tags = merge(var.tags, {
    "x-nitric-${local.stack_name}-name" = var.stack_name
    "x-nitric-${local.stack_name}-type" = "stack"
  })
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
  tags = merge(var.tags, {
    "x-nitric-${local.stack_name}-name" = var.stack_name
    "x-nitric-${local.stack_name}-type" = "stack"
  })
}

# Create an operational insights workspace
resource "azurerm_log_analytics_workspace" "log_analytics" {
  name                = "${var.stack_name}log${random_string.stack_id.result}"
  resource_group_name = azurerm_resource_group.resource_group.name
  location            = azurerm_resource_group.resource_group.location
  sku                 = "PerGB2018"
  retention_in_days   = 30
  tags = merge(var.tags, {
    "x-nitric-${local.stack_name}-name" = var.stack_name
    "x-nitric-${local.stack_name}-type" = "stack"
  })
}

# Create a random master password for the database server
resource "random_password" "database_master_password" {
  length  = 32
  special = true
  upper   = true
  lower   = true
  numeric  = true
}

# Create a database server if enabled
resource "azurerm_postgresql_flexible_server" "database" {
  count = var.enable_database ? 1 : 0

  name                = "${var.stack_name}db${random_string.stack_id.result}"
  resource_group_name = azurerm_resource_group.resource_group.name
  location            = azurerm_resource_group.resource_group.location
  version             = "13"
  administrator_login = "postgres"
  administrator_password = random_password.database_master_password.result
  zone                = "1"
  storage_mb          = 32768
  sku_name            = "B_Standard_B1ms"

  tags = merge(var.tags, {
    "x-nitric-${local.stack_name}-name" = var.stack_name
    "x-nitric-${local.stack_name}-type" = "stack"
  })
}

# Create a database if enabled
resource "azurerm_postgresql_flexible_server_database" "database" {
  count = var.enable_database ? 1 : 0

  name      = "nitric"
  server_id = azurerm_postgresql_flexible_server.database[0].id
  collation = "en_US.utf8"
  charset   = "utf8"
}

# Create a container app environment
resource "azurerm_container_app_environment" "environment" {
  name                       = "${var.stack_name}kube${random_string.stack_id.result}"
  resource_group_name        = azurerm_resource_group.resource_group.name
  location                   = azurerm_resource_group.resource_group.location
  log_analytics_workspace_id = azurerm_log_analytics_workspace.log_analytics.id
  tags = {
    "x-nitric-${local.stack_name}-name" = var.stack_name
    "x-nitric-${local.stack_name}-type" = "stack"
  }

  infrastructure_subnet_id = (var.enable_database || var.enable_storage_private_endpoints) ? azurerm_subnet.infrastructure_subnet[0].id : null
}
