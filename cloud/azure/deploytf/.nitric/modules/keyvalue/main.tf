# Create a new azure storage table
resource "azurerm_storage_table" "table" {
  name                 = var.name
  storage_account_name = var.storage_account_name
}
