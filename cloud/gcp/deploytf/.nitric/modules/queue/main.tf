
# Deploy a PubSub topic to serve as the queue
resource "google_pubsub_topic" "queue" {
  name = "queue"
  labels = {
    "x-nitric-${var.stack_id}-name" = var.queue_name
    "x-nitric-${var.stack_id}-type" = "queue"
  }
}

# Create a pull subscription for the topic to emulate a queue
resource "google_pubsub_subscription" "queue_subscription" {
  name = "${var.queue_name}-nitricqueue"
  topic = google_pubsub_topic.queue.name
  expiration_policy {
    # TODO: this is blank in the Pulumi provider - verify this is still correct
    ttl = ""
  }
}