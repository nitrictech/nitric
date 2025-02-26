# Create a new dapr component for a cron binding
resource "azurerm_container_app_environment_dapr_component" "schedule" {
  name                         = var.name
  container_app_environment_id = var.container_app_environment_id
  component_type               = "bindings.cron"
  version                      = "v1"

  metadata {
    name = "schedule"
    value = var.cron_expression
  }

  metadata {
    name = "route"
    value = "${var.target_event_token}/x-nitric-schedule/${var.name}"
  }

  scopes = [ var.target_app_id ]
}
