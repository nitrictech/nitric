resource "random_string" "unique_id" {
  length = 4
  special = false
  upper   = false
}

# Create a new pubsub topic
resource "google_pubsub_topic" "topic" {
  name = "${var.topic_name}-${random_string.unique_id.result}"

  kms_key_name = var.kms_key != "" ? var.kms_key : null

  labels = {
    "x-nitric-${var.stack_id}-name" = var.topic_name
    "x-nitric-${var.stack_id}-type" = "topic"
  }
}

# Create all CloudRun service subscriptions
resource "google_pubsub_subscription" "topic_subscriptions" {
  count = length(var.subscriber_services)
  name = "${var.subscriber_services[count.index].name}"
  topic = google_pubsub_topic.topic.name
  ack_deadline_seconds = 300

  retry_policy {
    minimum_backoff = "15s"
    maximum_backoff = "600s"
  }

  push_config {
    push_endpoint = "${var.subscriber_services[count.index].url}/x-nitric-topic/${var.topic_name}?token=${var.subscriber_services[count.index].event_token}"
    oidc_token {
      service_account_email = var.subscriber_services[count.index].invoker_service_account_email
    }
  }

  expiration_policy {
    ttl = ""
  }
}