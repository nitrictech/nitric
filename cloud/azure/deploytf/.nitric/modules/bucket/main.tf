# Generate a random id for the azure blob container
resource "random_id" "storage_id" {
  byte_length = 8

  keepers = {
    # Generate a new id each time we switch to a new AMI id
    storage_name = var.name
  }
}

# Create a new azure blob container
resource "azurerm_storage_container" "container" {
  # TODO: implement random suffix (requires runtime lookup)
  # name                  = "${var.name}st${random_id.storage_id.hex}"
  name                  = var.name
  storage_account_id    = var.storage_account_id
  container_access_type = "private"
  
  # metadata = var.tags
}

# Create an event subscription
resource "azurerm_eventgrid_event_subscription" "subscription" {
  for_each = var.listeners

  name                  = replace(each.key, "_", "-")
  scope                 = var.storage_account_id
  event_delivery_schema = "EventGridSchema"
  included_event_types  = each.value.event_type

  retry_policy {
    max_delivery_attempts = 10
    event_time_to_live    = 5
  }
  webhook_endpoint {
    max_events_per_batch           = 1
    active_directory_app_id_or_uri = each.value.active_directory_app_id_or_uri
    active_directory_tenant_id     = each.value.active_directory_tenant_id
    url                            = "${each.value.url}/${each.value.event_token}/x-nitric-notification/bucket/${var.name}"
  }
}
