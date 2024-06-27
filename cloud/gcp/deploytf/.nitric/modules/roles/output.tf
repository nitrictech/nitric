output "base_compute_role" {
  value = google_project_iam_custom_role.base_compute_role.id
  description = "The role ID for the Nitric base compute role"
}

output "bucket_read" {
  value = google_project_iam_custom_role.bucket_reader_role.id
  description = "The role ID for the Nitric bucket read role"
}

output "topic_publish" {
  value = google_project_iam_custom_role.topic_publisher_role.id
  description = "The role ID for the Nitric topic publish role"
}

output "bucket_write" {
  value = google_project_iam_custom_role.bucket_writer_role.id
  description = "The role ID for the Nitric bucket write role"
}

output "bucket_delete" {
  value = google_project_iam_custom_role.bucket_deleter_role.id
  description = "The role ID for the Nitric bucket delete role"
}

output "secret_access" {
  value       = google_project_iam_custom_role.secret_access_role.id
  description = "The role ID for the Nitric secrete access role"
}

output "secret_put" {
  value       = google_project_iam_custom_role.secret_put_role.id
  description = "The role ID for the Nitric secrete put role"
}

output "kv_read" {
  value       = google_project_iam_custom_role.kv_reader_role.id
  description = "The role ID for the Nitric kv read role"
}

output "kv_write" {
  value       = google_project_iam_custom_role.kv_writer_role.id
  description = "The role ID for the Nitric kv write role"
}

output "kv_delete" {
  value       = google_project_iam_custom_role.kv_deleter_role.id
  description = "The role ID for the Nitric kv write role"
}

output "queue_enqueue" {
  value       = google_project_iam_custom_role.queue_enqueue_role.id
  description = "The role ID for the Nitric queue enqueue role"
}

output "queue_dequeue" {
  value       = google_project_iam_custom_role.queue_dequeue_role.id
  description = "The role ID for the Nitric queue dequeue role"
}