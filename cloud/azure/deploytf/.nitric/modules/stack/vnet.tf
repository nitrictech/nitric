# Create a virtual network for the database server
resource "azurerm_virtual_network" "stack_network" {
  count               = var.vnet_name == null && var.enable_database ? 1 : 0
  name                = "nitric-database-vnet"
  resource_group_name = local.resource_group_name
  location            = var.location
  address_space       = ["10.0.0.0/16"]

  flow_timeout_in_minutes = 10
}

locals {
  vnet_name = var.vnet_name == null ? one(azurerm_virtual_network.stack_network).name : var.vnet_name
  subnet_id = var.subnet_id == null ? one(azurerm_subnet.database_infrastructure_subnet).id : var.subnet_id
}

# Create an infrastructure subnet for the database server
resource "azurerm_subnet" "database_infrastructure_subnet" {
  count                = var.subnet_id == null && var.enable_database ? 1 : 0
  name                 = "nitric-database-infrastructure-subnet"
  resource_group_name  = local.resource_group_name
  virtual_network_name = local.vnet_name
  address_prefixes     = ["10.0.0.0/16"]

  service_endpoints = [
    "Microsoft.Storage", 
    "Microsoft.Sql", 
    "Microsoft.KeyVault",
    "Microsoft.ContainerRegistry"
  ]

  depends_on = [azurerm_subnet.database_subnet]

  delegation {
    name = "container-instance-delegation"
    service_delegation {
      name    = "Microsoft.ContainerInstance/containerGroups"
      actions = ["Microsoft.Network/virtualNetworks/subnets/action"]
    }
  }

  delegation {
    name = "db-delegation"
    service_delegation {
      name    = "Microsoft.DBforPostgreSQL/flexibleServers"
      actions = ["Microsoft.Network/virtualNetworks/subnets/join/action"]
    }
  }
}

# Create a private zone for the database server
resource "azurerm_private_dns_zone" "database_dns_zone" {
  count               = var.enable_database ? 1 : 0
  name                = "db-private-dns.postgres.database.azure.com"
  resource_group_name = local.resource_group_name
}

# Create a private link service for the database server
resource "azurerm_private_dns_zone_virtual_network_link" "database_link_service" {
  count                 = var.enable_database ? 1 : 0
  name                  = "nitric-database-link-service"
  private_dns_zone_name = azurerm_private_dns_zone.database_dns_zone[0].name
  resource_group_name   = local.resource_group_name
  virtual_network_id    = local.vnet_name
  registration_enabled  = false

  tags = {
    "x-nitric-${var.stack_name}-name" = var.stack_name
    "x-nitric-${var.stack_name}-type" = "stack"
  }
}