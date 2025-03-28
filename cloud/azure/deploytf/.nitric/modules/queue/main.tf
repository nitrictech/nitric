# Create a new azure storage queue
resource "azurerm_storage_queue" "queue" {
  name                  = var.name
  storage_account_name  = var.storage_account_name

  metadata = var.tags
}