module "iam_roles" {
  source = "../roles"
  project_id = var.project_id
}

locals {
  is_bucket = var.resource_type == "Bucket"
  is_secret = var.resource_type == "Secret"
  is_kv = var.resource_type == "KeyValueStore"
  is_queue = var.resource_type == "Queue"
}

# Apply the IAM policy to the resource
resource "google_storage_bucket_iam_member" "bucket_iam_member_read" {
  count  = local.is_bucket && (contains(var.actions, "BucketFileGet") || contains(var.actions, "BucketFileList")) ? 1 : 0
  bucket = var.resource_name
  role   = module.iam_roles.bucket_read
  member = "serviceAccount:${var.service_account_email}"
}

resource "google_storage_bucket_iam_member" "bucket_iam_member_write" {
  count  = local.is_bucket && contains(var.actions, "BucketFilePut") ? 1 : 0
  bucket = var.resource_name
  role   = module.iam_roles.bucket_write
  member = "serviceAccount:${var.service_account_email}"
}

resource "google_storage_bucket_iam_member" "bucket_iam_member_delete" {
  count  = local.is_bucket && contains(var.actions, "BucketFileDelete") ? 1 : 0
  bucket = var.resource_name
  role   = module.iam_roles.bucket_delete
  member = "serviceAccount:${var.service_account_email}"
}

resource "google_secret_manager_secret_iam_member" "secret_iam_member_put" {
  project = var.project_id
  count  = local.is_secret && contains(var.actions, "SecretPut") ? 1 : 0
  secret_id = var.resource_name
  role    = module.iam_roles.secret_put
  member = "serviceAccount:${var.service_account_email}"
}

resource "google_secret_manager_secret_iam_member" "secret_iam_member_access" {
  project = var.project_id
  count  = local.is_secret && contains(var.actions, "SecretAccess") ? 1 : 0
  secret_id = var.resource_name
  role    = module.iam_roles.secret_access
  member = "serviceAccount:${var.service_account_email}"
}

resource "google_project_iam_member" "kv_iam_member_read" {
  project = var.project_id
  count  = local.is_kv && contains(var.actions, "KeyValueStoreRead") ? 1 : 0
  role    = module.iam_roles.kv_read
  member = "serviceAccount:${var.service_account_email}"
}

resource "google_project_iam_member" "kv_iam_member_delete" {
  project = var.project_id
  count  = local.is_kv && contains(var.actions, "KeyValueStoreDelete") ? 1 : 0
  role    = module.iam_roles.kv_read
  member = "serviceAccount:${var.service_account_email}"
}

resource "google_project_iam_member" "kv_iam_member_write" {
  project = var.project_id
  count  = local.is_kv && contains(var.actions, "KeyValueStoreWrite") ? 1 : 0
  role    = module.iam_roles.kv_read
  member = "serviceAccount:${var.service_account_email}"
}

resource "google_project_iam_member" "queue_iam_member_dequeue" {
  project = var.project_id
  count  = local.is_queue && contains(var.actions, "QueueDequeue") ? 1 : 0
  role    = module.iam_roles.queue_dequeue
  member = "serviceAccount:${var.service_account_email}"
}

resource "google_project_iam_member" "queue_iam_member_enqueue" {
  project = var.project_id
  count  = local.is_queue && contains(var.actions, "QueueEnqueue") ? 1 : 0
  role    = module.iam_roles.queue_enqueue
  member = "serviceAccount:${var.service_account_email}"
}