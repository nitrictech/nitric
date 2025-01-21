# Create an event grid topic
resource "azurerm_eventgrid_topic" "topic" {
  name                = var.name
  resource_group_name = var.resource_group_name
  location            = var.location
  tags = {
    "x-nitric-${var.stack_name}-name" = var.name
    "x-nitric-${var.stack_name}-type" = "topic"
  }
}

# token_response=$(curl -s -X POST -H "Content-Type: application/x-www-form-urlencoded" -d "client_id=${var.client_id}&client_secret=${var.client_secret}&scope=https://management.azure.com/.default&grant_type=client_credentials" "https://login.microsoftonline.com/${var.tenant_id}/oauth2/v2.0/token")
# access_token=$(echo $token_response | jq -r .access_token)

# Generate the authorization header
# auth_header="Authorization: Bearer $access_token"

# Poll the listener URL until it is available
resource "null_resource" "poll_url" {
  for_each = var.listeners

  provisioner "local-exec" {
    command = <<EOT
      echo "Polling URL: ${each.value.url}"
      max_attempts=10
      attempt=0

      while true; do
        echo "Polling attempt $attempt"
        echo "Sending subscription validation request to ${each.value.url}/${each.value.event_token}/x-nitric-topic/test"
        response=$(curl -s -w "%%{http_code}" -o /dev/null -X POST "${each.value.url}/${each.value.event_token}/x-nitric-topic/test" -H "aeg-event-type: SubscriptionValidation" -H "Content-Type: application/json" -d '{ "id": "", "data": "", "eventType": "", "subject": "", "dataVersion": "" }')
        if [ $exit_code -eq 0 ]; then
          echo "Service is available at ${each.value.url}"
          break
        fi
        
        echo "Got $response response was expecting 200"

        attempt=$((attempt + 1))
        if [ $attempt -eq $max_attempts ]; then
          echo "Service did not become available after $max_attempts attempts"
          exit 1
        fi

        echo "Waiting for service to be available..."
        sleep 10
      done
    EOT
  }
}

# Create an event subscription per listener
resource "azurerm_eventgrid_event_subscription" "subscription" {
  for_each = var.listeners

  name                  = each.key
  scope                 = azurerm_eventgrid_topic.topic.id

  retry_policy {
    max_delivery_attempts = 10
    event_time_to_live    = 5
  }
  webhook_endpoint {
    max_events_per_batch = 1
    active_directory_app_id_or_uri = each.value.active_directory_app_id_or_uri
    active_directory_tenant_id = each.value.active_directory_tenant_id
    url = "${each.value.url}/${each.value.event_token}/x-nitric-topic/${var.name}"
  }

  depends_on = [null_resource.poll_url]
}