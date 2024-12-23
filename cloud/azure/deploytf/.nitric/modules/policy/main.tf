# Create a new Azure role assignment

resource "azurerm_role_assignment" "role_assignment" {
  principal_id       = var.service_principal_id
  principal_type     = "ServicePrincipal"
  role_definition_id = var.role_definition_id
  scope              = var.scope
}