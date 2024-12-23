# Create a new azure storage table
resource "azurerm_storage_table" "table" {
  # Normalize the name by removing hyphens
  name                 = replace(var.name, "-", "")
  storage_account_name = var.storage_account_name
}
