# Generate a random id for the bucket
resource "random_id" "bucket_id" {
  byte_length = 8

  keepers = {
    # Generate a new id each time we switch to a new AMI id
    bucket_name = var.bucket_name
  }
}

# Get the location from the provider
data "google_client_config" "this" {
}

# Google Stora bucket
resource "google_storage_bucket" "bucket" {
  name          = "${var.bucket_name}-${random_id.bucket_id.hex}"
  location      = data.google_client_config.this.region
  storage_class = var.storage_class
  labels = {
    "x-nitric-${var.stack_id}-name" = var.bucket_name
    "x-nitric-${var.stack_id}-type" = "bucket"
  }
}

locals {
  has_notification_targets = length(var.notification_targets) > 0 ? 1 : 0
}

# Create a pubsub topic here for storage notifications
resource "google_pubsub_topic" "bucket_notification_topic" {
  count = local.has_notification_targets
  name  = "${var.bucket_name}-${random_id.bucket_id.hex}"
}


# Create a topic Iam binding for the storage notification topic
resource "google_pubsub_topic_iam_binding" "bucket_notification_topic_iam_binding" {
  count = local.has_notification_targets
  topic = google_pubsub_topic.bucket_notification_topic[0].id
  role  = "roles/pubsub.publisher"
  members = [
    "serviceAccount:${data.google_storage_project_service_account.storage_service_account.email_address}"
  ]
}

# Create a gcs storage notification that publishes events to the topic
resource "google_storage_notification" "bucket_notification" {
  for_each           = var.notification_targets
  bucket             = google_storage_bucket.bucket.name
  topic              = google_pubsub_topic.bucket_notification_topic[0].id
  event_types        = each.value.events
  payload_format     = "JSON_API_V1"
  object_name_prefix = each.value.prefix
  custom_attributes = {
    # The target is the notification target name
    "target" = each.value.name
  }
  depends_on         = [google_pubsub_topic_iam_binding.bucket_notification_topic_iam_binding]
}

# For each notification target create a pubsub subscription
resource "google_pubsub_subscription" "bucket_notification_subscription" {
  for_each             = var.notification_targets
  name                 = "${var.bucket_name}-${random_id.bucket_id.hex}"
  topic                = google_pubsub_topic.bucket_notification_topic[0].name
  ack_deadline_seconds = 300

  retry_policy {
    minimum_backoff = "15s"
    maximum_backoff = "600s"
  }

  # Filter will only allow messages that match the target attribute for itself
  filter = "attributes.target = \"${each.value.name}\""

  push_config {
    push_endpoint = "${each.value.url}/x-nitric-notification/bucket/${var.bucket_name}?token=${each.value.event_token}"
    oidc_token {
      service_account_email = each.value.invoker_service_account_email
    }
  }

  expiration_policy {
    ttl = ""
  }
}

data "google_storage_project_service_account" "storage_service_account" {
}

