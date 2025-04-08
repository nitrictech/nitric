# Create a virtual network for the stack
resource "azurerm_virtual_network" "vnet" {
  count = (var.enable_database || var.enable_storage_private_endpoints) ? 1 : 0

  name                = "${var.stack_name}-vnet"
  resource_group_name = azurerm_resource_group.resource_group.name
  location            = azurerm_resource_group.resource_group.location
  address_space       = ["10.0.0.0/16"]
  tags = merge(var.tags, {
    "x-nitric-${local.stack_name}-name" = var.stack_name
    "x-nitric-${local.stack_name}-type" = "stack"
  })
}

# Create a subnet for infrastructure (Container Apps)
resource "azurerm_subnet" "infrastructure_subnet" {
  count = (var.enable_database || var.enable_storage_private_endpoints) ? 1 : 0

  name                 = "${var.stack_name}-infrastructure-subnet"
  resource_group_name  = azurerm_resource_group.resource_group.name
  virtual_network_name = azurerm_virtual_network.vnet[0].name
  address_prefixes     = ["10.0.0.0/24"]
}

# Create a subnet for private endpoints
resource "azurerm_subnet" "private_endpoints_subnet" {
  count = (var.enable_database || var.enable_storage_private_endpoints) ? 1 : 0

  name                 = "${var.stack_name}-private-endpoints-subnet"
  resource_group_name  = azurerm_resource_group.resource_group.name
  virtual_network_name = azurerm_virtual_network.vnet[0].name
  address_prefixes     = ["10.0.1.0/24"]
}

# Create a private DNS zone for storage
resource "azurerm_private_dns_zone" "storage_dns_zone" {
  count = var.enable_storage_private_endpoints ? 1 : 0

  name                = "privatelink.blob.core.windows.net"
  resource_group_name = azurerm_resource_group.resource_group.name
  tags = merge(var.tags, {
    "x-nitric-${local.stack_name}-name" = var.stack_name
    "x-nitric-${local.stack_name}-type" = "stack"
  })
}

# Link the private DNS zone to the VNet
resource "azurerm_private_dns_zone_virtual_network_link" "storage_dns_zone_link" {
  count = var.enable_storage_private_endpoints ? 1 : 0

  name                  = "${var.stack_name}-storage-dns-link"
  resource_group_name   = azurerm_resource_group.resource_group.name
  private_dns_zone_name = azurerm_private_dns_zone.storage_dns_zone[0].name
  virtual_network_id    = azurerm_virtual_network.vnet[0].id
  tags = merge(var.tags, {
    "x-nitric-${local.stack_name}-name" = var.stack_name
    "x-nitric-${local.stack_name}-type" = "stack"
  })
}

# Create a private endpoint for storage
resource "azurerm_private_endpoint" "storage_private_endpoint" {
  count = var.enable_storage_private_endpoints ? 1 : 0

  name                = "${var.stack_name}-storage-private-endpoint"
  resource_group_name = azurerm_resource_group.resource_group.name
  location            = azurerm_resource_group.resource_group.location
  subnet_id           = azurerm_subnet.private_endpoints_subnet[0].id

  private_service_connection {
    name                           = "${var.stack_name}-storage-private-connection"
    private_connection_resource_id = azurerm_storage_account.storage[0].id
    subresource_names             = ["blob"]
    is_manual_connection          = false
  }

  private_dns_zone_group {
    name                 = "storage-dns-zone-group"
    private_dns_zone_ids = [azurerm_private_dns_zone.storage_dns_zone[0].id]
  }

  tags = merge(var.tags, {
    "x-nitric-${local.stack_name}-name" = var.stack_name
    "x-nitric-${local.stack_name}-type" = "stack"
  })
}

# Create a private DNS zone for the database
resource "azurerm_private_dns_zone" "database_dns_zone" {
  count = var.enable_database ? 1 : 0

  name                = "privatelink.postgres.database.azure.com"
  resource_group_name = azurerm_resource_group.resource_group.name
  tags = merge(var.tags, {
    "x-nitric-${local.stack_name}-name" = var.stack_name
    "x-nitric-${local.stack_name}-type" = "stack"
  })
}

# Link the private DNS zone to the VNet
resource "azurerm_private_dns_zone_virtual_network_link" "database_dns_zone_link" {
  count = var.enable_database ? 1 : 0

  name                  = "${var.stack_name}-database-dns-link"
  resource_group_name   = azurerm_resource_group.resource_group.name
  private_dns_zone_name = azurerm_private_dns_zone.database_dns_zone[0].name
  virtual_network_id    = azurerm_virtual_network.vnet[0].id
  tags = merge(var.tags, {
    "x-nitric-${local.stack_name}-name" = var.stack_name
    "x-nitric-${local.stack_name}-type" = "stack"
  })
}

# Create a private endpoint for the database
resource "azurerm_private_endpoint" "database_private_endpoint" {
  count = var.enable_database ? 1 : 0

  name                = "${var.stack_name}-database-private-endpoint"
  resource_group_name = azurerm_resource_group.resource_group.name
  location            = azurerm_resource_group.resource_group.location
  subnet_id           = azurerm_subnet.private_endpoints_subnet[0].id

  private_service_connection {
    name                           = "${var.stack_name}-database-private-connection"
    private_connection_resource_id = azurerm_postgresql_flexible_server.database[0].id
    subresource_names             = ["postgresqlServer"]
    is_manual_connection          = false
  }

  private_dns_zone_group {
    name                 = "database-dns-zone-group"
    private_dns_zone_ids = [azurerm_private_dns_zone.database_dns_zone[0].id]
  }

  tags = merge(var.tags, {
    "x-nitric-${local.stack_name}-name" = var.stack_name
    "x-nitric-${local.stack_name}-type" = "stack"
  })
}