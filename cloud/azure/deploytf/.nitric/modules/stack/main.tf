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
  count = var.resource_group_name == null ? 1 : 0

  name     = "${var.stack_name}-rg${random_string.stack_id.result}"
  location = var.location
  tags = merge(var.tags, {
    "x-nitric-${local.stack_name}-name" = var.stack_name
    "x-nitric-${local.stack_name}-type" = "stack"
  })
}

locals {
  resource_group_name = var.resource_group_name == null ? one(azurerm_resource_group.resource_group).name : var.resource_group_name
}

# Create an azure storage account
resource "azurerm_storage_account" "storage" {
  count               = var.enable_storage ? 1 : 0
  name                = "${var.stack_name}sa${random_string.stack_id.result}"
  resource_group_name = local.resource_group_name
  location            = var.location
  account_tier        = "Standard"
  access_tier         = "Hot"
  # TODO: Make configurable  
  account_replication_type = "LRS"
  account_kind             = "StorageV2"
  dynamic "network_rules" {
    for_each = var.private_endpoints ? [1] : []
    content {
      default_action = "Deny"
      virtual_network_subnet_ids = [
        local.subnet_id
      ]
    }
  }

  public_network_access_enabled = !var.private_endpoints
  
  tags = merge(var.tags, {
    "x-nitric-${local.stack_name}-name" = var.stack_name
    "x-nitric-${local.stack_name}-type" = "stack"
  })
}

# Create a keyvault if secrets are enabled
resource "azurerm_key_vault" "keyvault" {
  count = var.enable_keyvault ? 1 : 0

  name                       = "${var.stack_name}kv${random_string.stack_id.result}"
  resource_group_name        = local.resource_group_name
  location                   = var.location
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
  resource_group_name = local.resource_group_name
  location            = var.location
}

# Create a container registry for storing images
resource "azurerm_container_registry" "container_registry" {
  name                = "${var.stack_name}cr${random_string.stack_id.result}"
  resource_group_name = local.resource_group_name
  location            = var.location
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
  resource_group_name = local.resource_group_name
  location            = var.location
  sku                 = "PerGB2018"
  retention_in_days   = 30
  tags = merge(var.tags, {
    "x-nitric-${local.stack_name}-name" = var.stack_name
    "x-nitric-${local.stack_name}-type" = "stack"
  })
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
  resource_group_name    = local.resource_group_name
  location               = var.location
  version                = "14"
  administrator_login    = "nitric"
  administrator_password = random_password.database_master_password[0].result

  zone = "1"

  public_network_access_enabled = false

  delegated_subnet_id = local.subnet_id
  private_dns_zone_id = one(azurerm_private_dns_zone.database_dns_zone[0]) != null ? one(azurerm_private_dns_zone.database_dns_zone[0]).id : null

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
    local.subnet_id,
    azurerm_private_dns_zone.database_dns_zone,
    azurerm_private_dns_zone_virtual_network_link.database_link_service
  ]
}

# Create a new container app managed environment
resource "azurerm_container_app_environment" "environment" {
  name                       = "${var.stack_name}kube${random_string.stack_id.result}"
  resource_group_name        = local.resource_group_name
  location                   = var.location
  log_analytics_workspace_id = azurerm_log_analytics_workspace.log_analytics.id
  tags = {
    "x-nitric-${local.stack_name}-name" = var.stack_name
    "x-nitric-${local.stack_name}-type" = "stack"
  }

  infrastructure_subnet_id = var.enable_database ? local.subnet_id: null
}
