# Create private DNS zone for blob storage if required
resource "azurerm_private_dns_zone" "storage_blob" {
  count               = var.enable_storage && var.private_endpoints && var.create_dns_zones ? 1 : 0
  name                = "privatelink.blob.core.windows.net"
  resource_group_name = local.resource_group_name
}

# Create private DNS zone for queue storage if required
resource "azurerm_private_dns_zone" "storage_queue" {
  count               = var.enable_storage && var.private_endpoints && var.create_dns_zones ? 1 : 0
  name                = "privatelink.queue.core.windows.net"
  resource_group_name = local.resource_group_name
}

# Create private DNS zone for table storage if required
resource "azurerm_private_dns_zone" "storage_table" {
  count               = var.enable_storage && var.private_endpoints && var.create_dns_zones ? 1 : 0
  name                = "privatelink.table.core.windows.net"
  resource_group_name = local.resource_group_name
}

# Link private DNS zones to VNet if required
resource "azurerm_private_dns_zone_virtual_network_link" "storage_blob_link" {
  count                 = var.enable_storage && var.private_endpoints && var.create_dns_zones ? 1 : 0
  name                  = "${local.stack_name}-blob-link"
  resource_group_name   = local.resource_group_name
  private_dns_zone_name = azurerm_private_dns_zone.storage_blob[0].name
  virtual_network_id    = local.vnet_id
  registration_enabled  = false
}

resource "azurerm_private_dns_zone_virtual_network_link" "storage_queue_link" {
  count                 = var.enable_storage && var.private_endpoints && var.create_dns_zones ? 1 : 0
  name                  = "${local.stack_name}-queue-link"
  resource_group_name   = local.resource_group_name
  private_dns_zone_name = azurerm_private_dns_zone.storage_queue[0].name
  virtual_network_id    = local.vnet_id
  registration_enabled  = false
}

resource "azurerm_private_dns_zone_virtual_network_link" "storage_table_link" {
  count                 = var.enable_storage && var.private_endpoints && var.create_dns_zones ? 1 : 0
  name                  = "${local.stack_name}-table-link"
  resource_group_name   = local.resource_group_name
  private_dns_zone_name = azurerm_private_dns_zone.storage_table[0].name
  virtual_network_id    = local.vnet_id
  registration_enabled  = false
}

# Create private endpoint for blob storage
resource "azurerm_private_endpoint" "storage_blob" {
  count               = var.enable_storage && var.private_endpoints ? 1 : 0
  name                = "${local.stack_name}-blob-pe"
  location            = var.location
  resource_group_name = local.resource_group_name
  subnet_id           = local.subnet_id

  private_service_connection {
    name                           = "${local.stack_name}-blob-connection"
    private_connection_resource_id = azurerm_storage_account.storage[0].id
    subresource_names             = ["blob"]
    is_manual_connection          = false
  }

  dynamic "private_dns_zone_group" {
    for_each = var.create_dns_zones ? [1] : []
    content {
      name                 = "default"
      private_dns_zone_ids = [azurerm_private_dns_zone.storage_blob[0].id]
    }
  }
}

# Create private endpoint for queue storage
resource "azurerm_private_endpoint" "storage_queue" {
  count               = var.enable_storage && var.private_endpoints ? 1 : 0
  name                = "${local.stack_name}-queue-pe"
  location            = var.location
  resource_group_name = local.resource_group_name
  subnet_id           = local.subnet_id

  private_service_connection {
    name                           = "${local.stack_name}-queue-connection"
    private_connection_resource_id = azurerm_storage_account.storage[0].id
    subresource_names             = ["queue"]
    is_manual_connection          = false
  }

  dynamic "private_dns_zone_group" {
    for_each = var.create_dns_zones ? [1] : []
    content {
      name                 = "default"
      private_dns_zone_ids = [azurerm_private_dns_zone.storage_queue[0].id]
    }
  }
}

# Create private endpoint for table storage
resource "azurerm_private_endpoint" "storage_table" {
  count               = var.enable_storage && var.private_endpoints ? 1 : 0
  name                = "${local.stack_name}-table-pe"
  location            = var.location
  resource_group_name = local.resource_group_name
  subnet_id           = local.subnet_id

  private_service_connection {
    name                           = "${local.stack_name}-table-connection"
    private_connection_resource_id = azurerm_storage_account.storage[0].id
    subresource_names             = ["table"]
    is_manual_connection          = false
  }

  dynamic "private_dns_zone_group" {
    for_each = var.create_dns_zones ? [1] : []
    content {
      name                 = "default"
      private_dns_zone_ids = [azurerm_private_dns_zone.storage_table[0].id]
    }
  }
}



