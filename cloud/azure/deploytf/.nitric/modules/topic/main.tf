# Create an event grid topic
resource "azurerm_eventgrid_topic" "topic" {
  name                = var.name
  resource_group_name = var.resource_group_name
  location            = var.location
  tags = var.tags
}

# Create an event subscription per listener
resource "azurerm_eventgrid_event_subscription" "subscription" {
  for_each = var.listeners

  name  = each.key
  scope = azurerm_eventgrid_topic.topic.id

  retry_policy {
    max_delivery_attempts = 10
    event_time_to_live    = 5
  }
  webhook_endpoint {
    max_events_per_batch           = 1
    active_directory_app_id_or_uri = each.value.active_directory_app_id_or_uri
    active_directory_tenant_id     = each.value.active_directory_tenant_id
    url                            = "${each.value.url}/${each.value.event_token}/x-nitric-topic/${var.name}"
  }
}
