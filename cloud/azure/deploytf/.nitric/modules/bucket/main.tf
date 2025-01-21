# Create a new azure blob container
resource "azurerm_storage_container" "container" {
  name                  = var.name
  storage_account_id    = var.storage_account_id
  container_access_type = "private"
}

# FIXME: This is a duplicate for the topic module
# Need to deduplicate this code by moving it into the service module and tying its output as a dependency` 
# Poll the listener URL until it is available
# resource "null_resource" "poll_url" {
#   for_each = var.listeners

#   provisioner "local-exec" {
#     command = <<EOT
#       echo "Polling URL: ${each.value.url}"
#       max_attempts=10
#       attempt=0

#       while true; do
#         echo "Polling attempt $attempt"
#         echo "Sending subscription validation request to ${each.value.url}/${each.value.event_token}/x-nitric-topic/test"
#         response=$(curl -s -w "%%{http_code}" -o /dev/null -X POST "${each.value.url}/${each.value.event_token}/x-nitric-topic/test" -H "aeg-event-type: SubscriptionValidation" -H "Content-Type: application/json" -d '{ "id": "", "data": "", "eventType": "", "subject": "", "dataVersion": "" }')
#         echo "Got $response response was expecting 200"
#         exit_code=$?
#         if [ $exit_code -eq 0 ]; then
#           echo "Service is available at ${each.value.url}"
#           break
#         fi

#         attempt=$((attempt + 1))
#         if [ $attempt -eq $max_attempts ]; then
#           echo "Service did not become available after $max_attempts attempts"
#           exit 1
#         fi

#         echo "Waiting for service to be available..."
#         sleep 10
#       done
#     EOT
#   }
# }

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
    max_events_per_batch = 1
    active_directory_app_id_or_uri = each.value.active_directory_app_id_or_uri
    active_directory_tenant_id = each.value.active_directory_tenant_id
    url = "${each.value.url}/${each.value.event_token}/x-nitric-notification/bucket/${var.name}"
  }
}
