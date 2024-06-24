data "google_project" "project" {
}

locals {
  is_bucket = var.resource_type == "Bucket"
  is_secret = var.resource_type == "Secret"
  is_kv = var.resource_type == "KeyValueStore"
  is_queue = var.resource_type == "Queue"
  is_topic = var.resource_type == "Topic"
}

# Apply the IAM policy to the resource
resource "google_pubsub_iam_member" "topic_iam_member_publish" {
  count  = local.is_topic && contains(var.actions, "TopicPublish") ? 1 : 0
  bucket = var.resource_name
  role   = var.iam_roles.topic_publish
  member = "serviceAccount:${var.service_account_email}"
}

# Apply the IAM policy to the resource
resource "google_storage_bucket_iam_member" "bucket_iam_member_read" {
  count  = local.is_bucket && (contains(var.actions, "BucketFileGet") || contains(var.actions, "BucketFileList")) ? 1 : 0
  bucket = var.resource_name
  role   = var.iam_roles.bucket_read
  member = "serviceAccount:${var.service_account_email}"
}

resource "google_storage_bucket_iam_member" "bucket_iam_member_write" {
  count  = local.is_bucket && contains(var.actions, "BucketFilePut") ? 1 : 0
  bucket = var.resource_name
  role   = var.iam_roles.bucket_write
  member = "serviceAccount:${var.service_account_email}"
}

resource "google_storage_bucket_iam_member" "bucket_iam_member_delete" {
  count  = local.is_bucket && contains(var.actions, "BucketFileDelete") ? 1 : 0
  bucket = var.resource_name
  role   = var.iam_roles.bucket_delete
  member = "serviceAccount:${var.service_account_email}"
}

resource "google_secret_manager_secret_iam_member" "secret_iam_member_put" {
  count  = local.is_secret && contains(var.actions, "SecretPut") ? 1 : 0
  secret_id = var.resource_name
  role    = var.iam_roles.secret_put
  member = "serviceAccount:${var.service_account_email}"
}

resource "google_secret_manager_secret_iam_member" "secret_iam_member_access" {
  count  = local.is_secret && contains(var.actions, "SecretAccess") ? 1 : 0
  secret_id = var.resource_name
  role    = var.iam_roles.secret_access
  member = "serviceAccount:${var.service_account_email}"
}

resource "google_project_iam_member" "kv_iam_member_read" {
  project = data.google_project.project.project_id
  count  = local.is_kv && contains(var.actions, "KeyValueStoreRead") ? 1 : 0
  role    = var.iam_roles.kv_read
  member = "serviceAccount:${var.service_account_email}"
}

resource "google_project_iam_member" "kv_iam_member_delete" {
  project = data.google_project.project.project_id
  count  = local.is_kv && contains(var.actions, "KeyValueStoreDelete") ? 1 : 0
  role    = var.iam_roles.kv_delete
  member = "serviceAccount:${var.service_account_email}"
}

resource "google_project_iam_member" "kv_iam_member_write" {
  project = data.google_project.project.project_id
  count  = local.is_kv && contains(var.actions, "KeyValueStoreWrite") ? 1 : 0
  role    = var.iam_roles.kv_write
  member = "serviceAccount:${var.service_account_email}"
}

resource "google_pubsub_iam_member" "queue_iam_member_dequeue" {
  project = data.google_project.project.project_id
  count  = local.is_queue && contains(var.actions, "QueueDequeue") ? 1 : 0
  role    = var.iam_roles.queue_dequeue
  member = "serviceAccount:${var.service_account_email}"
  topic = var.resource_name
}

resource "google_pubsub_iam_member" "queue_iam_member_enqueue" {
  project = data.google_project.project.project_id
  count  = local.is_queue && contains(var.actions, "QueueEnqueue") ? 1 : 0
  role    = var.iam_roles.queue_enqueue
  member = "serviceAccount:${var.service_account_email}"
  topic = var.resource_name
}