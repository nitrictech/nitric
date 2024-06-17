# Generate a random id for the bucket
resource "random_id" "bucket_id" {
  byte_length = 8

  keepers = {
    # Generate a new id each time we switch to a new AMI id
    bucket_name = var.bucket_name
  }
}

# Google Stora bucket
resource "google_storage_bucket" "bucket" {
  name = "${var.bucket_name}-${random_id.bucket_id.hex}"
  location      = var.bucket_location
  project       = var.project_id
  storage_class = var.storage_class
  labels = {
    "x-nitric-${var.stack_id}-name" = var.bucket_name
    "x-nitric-${var.stack_id}-type" = "bucket"
  }
}

# Create a pubsub topic here for storage notifications
resource "google_pubsub_topic" "bucket_notification_topic" {
  count = length(var.notification_targets) > 0 ? 1 : 0
  name = "${var.bucket_name}-${random_id.bucket_id.hex}"
  project = var.project_id
}

# Create a gcs storage notification that publishes events to the topic
resource "google_storage_notification" "bucket_notification" {
  bucket = google_storage_bucket.bucket.name
  topic = google_pubsub_topic.bucket_notification_topic.name
  event_types = ["OBJECT_FINALIZE", "OBJECT_DELETE"]
  payload_format = "JSON_API_V1"
}

# For each notification target create a pubsub subscription
resource "google_pubsub_subscription" "bucket_notification_subscription" {
  for_each = var.notification_targets
  name = "${var.bucket_name}-${random_id.bucket_id.hex}"
  topic = google_pubsub_topic.bucket_notification_topic.name
  project = var.project_id
  ack_deadline_seconds = 300

  retry_policy {
    minimum_backoff = "15s"
    maximum_backoff = "600s"
  }
  
  # FIXME: improve this filter
  filter = join(" OR ", "attributes.eventType = ${each.value.events}")

  push_config {
    push_endpoint = each.value.url
    oidc_token {
      service_account_email = each.value.event_token
    }
  }

  expiration_policy {
    ttl = ""
  }
}

# Get the projects google cloud storage account
data "google_project_service_identity" "gcs_service_account" {
  provider = google
  service = "storage-api.googleapis.com"
}

# Create a topic Iam binding for the storage notification topic
resource "google_pubsub_topic_iam_binding" "bucket_notification_topic_iam_binding" {
  count = length(var.notification_targets) > 0 ? 1 : 0
  topic = google_pubsub_topic.bucket_notification_topic.name
  role = "roles/pubsub.publisher"
  members = [
    "serviceAccount:${data.google_project_service_identity.gcs_service_account.email}"
  ]
}