# Create a random string for the api management service name
resource "random_string" "api_id" {
  length  = 8
  special = false
  upper   = false
}

locals {
  api_name = "${var.name}-${random_string.api_id.result}"
}

# Deploy a new azure api management service
resource "azurerm_api_management" "api" {
  name                = local.api_name
  resource_group_name = var.resource_group_name
  location            = var.location
  publisher_name      = var.publisher_name
  publisher_email     = var.publisher_email
  sku_name            = "Consumption_0"
  identity {
    type         = "UserAssigned"
    identity_ids = [var.app_identity]
  }
}


# Deploy a new azure api management api
resource "azurerm_api_management_api" "api" {
  name                = local.api_name
  resource_group_name = var.resource_group_name
  api_management_name = azurerm_api_management.api.name
  # TODO: This may need to increment if the api changes   
  revision     = "1"
  display_name = "${var.name}-api"
  protocols   = ["https"]
  description = var.description
  subscription_required = false
  import {
    content_format = "openapi+json"
    content_value  = var.openapi_spec
  }
}

# Create api operation policies
resource "azurerm_api_management_api_operation_policy" "api" {
  for_each = var.operation_policy_templates

  api_name            = azurerm_api_management_api.api.name
  api_management_name = azurerm_api_management.api.name
  resource_group_name = var.resource_group_name
  operation_id        = each.key
  xml_content         = each.value
}
