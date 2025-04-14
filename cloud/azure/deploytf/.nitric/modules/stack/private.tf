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
} 