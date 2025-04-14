# Create private DNS zone for storage services if required
resource "azurerm_private_dns_zone" "storage" {
  count               = var.enable_storage && var.private_endpoints && var.create_dns_zones ? 1 : 0
  name                = "privatelink.blob.core.windows.net"
  resource_group_name = local.resource_group_name
}

# Link private DNS zone to VNet if required
resource "azurerm_private_dns_zone_virtual_network_link" "storage_link" {
  count                 = var.enable_storage && var.private_endpoints && var.create_dns_zones ? 1 : 0
  name                  = "${local.stack_name}-storage-link"
  resource_group_name   = local.resource_group_name
  private_dns_zone_name = azurerm_private_dns_zone.storage[0].name
  virtual_network_id    = var.subnet_id != null ? join("/", slice(split("/", var.subnet_id), 0, 9)) : azurerm_virtual_network.stack_network[0].id
  registration_enabled  = false
}

# Create private endpoint for all storage services
resource "azurerm_private_endpoint" "storage" {
  count               = var.enable_storage && var.private_endpoints ? 1 : 0
  name                = "${local.stack_name}-storage-pe"
  location            = var.location
  resource_group_name = local.resource_group_name
  subnet_id           = var.subnet_id

  private_service_connection {
    name                           = "${local.stack_name}-storage-connection"
    private_connection_resource_id = azurerm_storage_account.storage[0].id
    subresource_names             = ["blob", "queue", "table"]
    is_manual_connection          = false
  }

  dynamic "private_dns_zone_group" {
    for_each = var.create_dns_zones ? [1] : []
    content {
      name                 = "default"
      private_dns_zone_ids = [azurerm_private_dns_zone.storage[0].id]
    }
  }
}



