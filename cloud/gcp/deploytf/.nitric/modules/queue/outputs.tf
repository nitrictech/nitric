output "name" {
  description = "The name of the deployed queue."
  value       = google_pubsub_subscription.queue_subscription.name
}