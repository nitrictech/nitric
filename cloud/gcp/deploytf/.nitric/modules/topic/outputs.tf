# TODO: confirm outputs
output "topic_name" {
  description = "The name of the topic."
  value       = google_pubsub_topic.topic.name
}
